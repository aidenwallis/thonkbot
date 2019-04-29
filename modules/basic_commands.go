package modules

import (
	"strings"
	"time"

	"github.com/aidenwallis/thonkbot/bot"
	"github.com/aidenwallis/thonkbot/common"
	"github.com/aidenwallis/thonkbot/mysql"
)

type BasicCommands struct {
	commands map[string]commandFunc
	bot      bot.Bot
}

func (m *BasicCommands) Init(bot bot.Bot) {
	m.bot = bot
	m.commands = make(map[string]commandFunc)

	m.registerCommand([]string{"rq"}, m.randomQuote)
	m.registerCommand([]string{"scan"}, m.scanCommand)
	m.registerCommand([]string{"globalscan"}, m.globalscanCommand)
	// m.registerCommand([]string{"firstmsg"}, m.firstmsg)
	m.registerCommand([]string{"linecount"}, m.linecount)
}

func (m *BasicCommands) Run(msg *common.Message) {
	if len(msg.Split) == 0 {
		return
	}

	prefix := msg.Split[0]
	if prefix[0] != '!' || len(prefix) < 2 {
		return
	}

	command := strings.ToLower(prefix[1:])
	commandFunc, ok := m.commands[command]
	if !ok {
		return
	}

	commandFunc(msg)
}

func (m *BasicCommands) registerCommand(aliases []string, cb commandFunc) {
	for _, alias := range aliases {
		m.commands[alias] = cb
	}
}

func (m *BasicCommands) randomQuote(msg *common.Message) {
	target := msg.User.Username
	if len(msg.Split) >= 2 && len(msg.Split[1]) > 1 {
		target = strings.TrimPrefix(strings.ToLower(msg.Split[1]), "@")
	}
	quote, err := mysql.FetchRandomQuote(target, msg.ChannelName)
	if err != nil {
		m.bot.Log().WithError(err).Error("Failed to get random quote")
		m.bot.Sayf("@%s, Failed to get quote for user %s", msg.User.Username, target)
		return
	}

	if quote == nil {
		m.bot.Sayf("@%s, No quotes exist for user %s!", msg.User.Username, target)
		return
	}

	formattedDate := quote.CreatedAt.UTC().Format(time.ANSIC)
	m.bot.Sayf("[%s] %s: %s", formattedDate, quote.Username, quote.Message)
}

func (m *BasicCommands) linecount(msg *common.Message) {
	target := msg.User.Username
	if len(msg.Split) >= 2 && len(msg.Split[1]) > 1 {
		target = strings.TrimPrefix(strings.ToLower(msg.Split[1]), "@")
	}
	count, err := mysql.FetchLineCount(target, msg.ChannelName)
	if err != nil {
		m.bot.Log().WithError(err).Error("Failed to fetch line count")
		m.bot.Sayf("@%s, Failed to get line count for user %s!", msg.User.Username, target)
	}

	m.bot.Sayf("@%s, User %s has sent %d lines in chat so far.", msg.User.Username, target, count)
}

func (m *BasicCommands) scanCommand(msg *common.Message) {
	if len(msg.Split) < 3 {
		return
	}
	target := msg.Split[1]
	if len(target) == 0 {
		return
	}
	if target[0] == '@' {
		target = target[1:]
	}

	query := strings.Join(msg.Split[2:], " ")
	if len(query) == 0 {
		return
	}

	count, err := mysql.ScanMessages(msg.ChannelName, msg.User.Username, query)
	if err != nil {
		m.bot.Log().WithField("username", target).WithError(err).Error("Failed to scan the messages")
		m.bot.Sayf("@%s, Failed to scan user %s's messages.", msg.User.Username, target)
		return
	}

	plural := "s"
	if count == 1 {
		plural = ""
	}
	m.bot.Sayf(`@%s, User %s has said "%s" %d time%s in #%s.`, msg.User.Username, target, query, count, plural, msg.ChannelName)
}

func (m *BasicCommands) globalscanCommand(msg *common.Message) {
	if len(msg.Split) < 2 {
		return
	}

	query := strings.Join(msg.Split[1:], " ")
	if len(query) == 0 {
		return
	}

	count, err := mysql.GlobalScanMessages(msg.ChannelName, query)
	if err != nil {
		m.bot.Log().WithError(err).Error("Failed to global scan the messages")
		m.bot.Sayf("@%s, Failed to scan the global message database.", msg.User.Username)
		return
	}

	plural := "s"
	if count == 1 {
		plural = ""
	}
	m.bot.Sayf(`@%s, "%s" has been said %d time%s in #%s.`, msg.User.Username, query, count, plural, msg.ChannelName)
}
