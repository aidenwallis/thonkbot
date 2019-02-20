package modules

import (
	"strings"

	"github.com/aidenwallis/thonkbot/bot"
	"github.com/aidenwallis/thonkbot/common"
	"github.com/aidenwallis/thonkbot/mysql"
)

type MainBotCommands struct {
	bot      bot.Bot
	commands map[string]commandFunc
}

func (m *MainBotCommands) Init(bot bot.Bot) {
	m.bot = bot
	m.commands = make(map[string]commandFunc)

	m.registerCommand([]string{"join"}, m.joinCommand)
	m.registerCommand([]string{"channels"}, m.channelsCount)
}

func (m *MainBotCommands) Run(msg *common.Message) {
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

func (m *MainBotCommands) registerCommand(aliases []string, cb commandFunc) {
	for _, alias := range aliases {
		m.commands[alias] = cb
	}
}

func (m *MainBotCommands) joinCommand(msg *common.Message) {
	toJoin := msg.User.Username
	if len(msg.Split) > 1 {
		toJoin = msg.Split[1]
	}

	if toJoin != msg.User.Username {
		isAdmin, err := mysql.IsAdmin(msg.User.Username)
		if err != nil {
			m.bot.Log().WithError(err).Error("Error while fetching if they're a bot admin")
			m.bot.Sayf("@%s, Failed to verify as to whether you are allowed to join me to other people's channels! (Bot admins only)", msg.User.Username)
			return
		}
		if !isAdmin {
			toJoin = msg.User.Username
		}
	}

	isJoined, err := mysql.CheckJoined(toJoin)
	if err != nil {
		m.bot.Log().WithError(err).Error("Failed to check if user has already joined Thonkbot")
		m.bot.Sayf("@%s, Failed to join channel, please try again later.", msg.User.Username)
		return
	}

	if isJoined {
		m.bot.Sayf("@%s, I have already joined #%s! cmonBruh", msg.User.Username, toJoin)
		return
	}

	channel, err := mysql.AddChannel(toJoin, msg.User.Username)
	if err != nil {
		m.bot.Log().WithError(err).Error("Failed to create channel in database")
		m.bot.Sayf("@%s, Failed to join channel. Please try again later.", msg.User.Username)
		return
	}

	m.bot.JoinChannel(channel)

	m.bot.Sayf("@%s, I am now joining channel #%s", msg.User.Username, toJoin)
	m.bot.Sayf(`I am required to inform you that I store your chat messages so that I can run queries against them to make me work. If at any time you want the bot to leave your channel and delete all recorded messages, please type "!thonkleave" in the channels\' chatroom. By using this command you are giving me permission to log messages!`)
	m.bot.Sayf("Thanks for trying this out! Feel free to share feedback with me on Twitter: https://twitter.com/WallisDev")
}

func (m *MainBotCommands) channelsCount(msg *common.Message) {
	count := m.bot.ChannelCount()
	m.bot.Sayf("@%s, I am currently joined to %d channels! PogChamp", msg.User.Username, count)
}
