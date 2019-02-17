package singlebot

import (
	"fmt"

	"github.com/aidenwallis/thonkbot/bot"
	"github.com/aidenwallis/thonkbot/botmanager"
	"github.com/aidenwallis/thonkbot/common"
	"github.com/aidenwallis/thonkbot/modules"
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/sirupsen/logrus"
)

type SingleBot struct {
	conn     *twitch.Client
	channel  common.Channel
	log      logrus.FieldLogger
	modules  []modules.Module
	manager  *botmanager.BotManager
	messages chan *common.Message
}

func New(channel common.Channel, logger logrus.FieldLogger, conn *twitch.Client, manager *botmanager.BotManager) bot.Bot {
	b := &SingleBot{
		conn:    conn,
		channel: channel,
		manager: manager,
		log:     logger.WithField("channel-name", channel.Name),
		modules: []modules.Module{
			&modules.BasicCommands{},
			&modules.LogMessage{},
			&modules.BroadcasterControl{},
		},
	}
	if channel.ID == 0 {
		b.modules = []modules.Module{
			&modules.MainBotCommands{},
		}
	}
	for _, module := range b.modules {
		module.Init(b)
	}
	return b
}

func (b *SingleBot) HandleMessage(msg *common.Message) {
	if b.messages == nil {
		b.messages = make(chan *common.Message, 15)
		go b.handleMessages()
	}
	select {
	case b.messages <- msg:
	default:
		b.log.Warn("Message buffer is full, am dropping a message! %+v", msg)
	}
}

func (b *SingleBot) ChannelCount() int {
	return b.manager.ChannelCount()
}

func (b *SingleBot) Channel() *common.Channel {
	return &b.channel
}

func (b *SingleBot) Log() logrus.FieldLogger {
	return b.log
}

func (b *SingleBot) Close() {
	if b.messages != nil {
		close(b.messages)
	}
}

func (b *SingleBot) JoinChannel(channel common.Channel) {
	b.manager.JoinChannel(channel)
}

func (b *SingleBot) LeaveChannel(channel string) {
	b.manager.LeaveChannel(channel)
}

func (b *SingleBot) Say(line string) {
	b.conn.Say(b.channel.Name, line)
}

func (b *SingleBot) Sayf(msg string, a ...interface{}) {
	b.Say(fmt.Sprintf(msg, a...))
}

func (b *SingleBot) handleMessages() {
	for msg := range b.messages {
		go b.handleMessage(msg)
	}
}

func (b *SingleBot) handleMessage(msg *common.Message) {
	for _, module := range b.modules {
		go module.Run(msg)
	}
}
