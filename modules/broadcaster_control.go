package modules

import (
	"strings"

	"github.com/aidenwallis/thonkbot/bot"
	"github.com/aidenwallis/thonkbot/common"
	"github.com/aidenwallis/thonkbot/mysql"
)

type BroadcasterControl struct {
	bot      bot.Bot
	commands map[string]commandFunc
}

func (m *BroadcasterControl) Init(bot bot.Bot) {
	m.bot = bot
	m.commands = make(map[string]commandFunc)

	m.registerCommand([]string{"thonkleave"}, m.leaveChannel)
}

func (m *BroadcasterControl) Run(msg *common.Message) {
	_, isBroadcaster := msg.Message.Tags["broadcaster"]
	if !isBroadcaster && msg.Message.ChannelID == msg.User.UserID {
		isBroadcaster = true
	}

	if !isBroadcaster {
		return
	}

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

func (m *BroadcasterControl) registerCommand(aliases []string, cb commandFunc) {
	for _, alias := range aliases {
		m.commands[alias] = cb
	}
}

func (m *BroadcasterControl) leaveChannel(msg *common.Message) {
	err := mysql.DeleteChannel(m.bot.Channel().ID, msg.ChannelName)
	if err != nil {
		m.bot.Log().WithError(err).Error("Failed to leave channel")
		m.bot.Sayf("@%s, Failed to leave channel #%s. Please try again later.", msg.User.Username, msg.ChannelName)
		return
	}
	m.bot.Sayf("@%s, Now leaving channel #%s and deleting all logs associated!", msg.User.Username, msg.ChannelName)
	m.bot.Close()
	m.bot.LeaveChannel(msg.ChannelName)
}
