package modules

import (
	"github.com/aidenwallis/thonkbot/bot"
	"github.com/aidenwallis/thonkbot/common"
	"github.com/aidenwallis/thonkbot/mysql"
)

type LogMessage struct {
	bot bot.Bot
}

func (m *LogMessage) Init(bot bot.Bot) {
	m.bot = bot
}

func (m *LogMessage) Run(msg *common.Message) {
	err := mysql.LogMessage(msg.User.Username, msg.Message.Text, msg.ChannelName)
	if err != nil {
		m.bot.Log().WithError(err).Error("Failed to log message")
	}
}
