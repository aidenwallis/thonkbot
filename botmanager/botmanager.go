package botmanager

import (
	"strings"
	"sync"

	"github.com/aidenwallis/thonkbot/bot"
	"github.com/aidenwallis/thonkbot/common"
	"github.com/aidenwallis/thonkbot/config"
	"github.com/aidenwallis/thonkbot/mysql"
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/sirupsen/logrus"
)

type BotManager struct {
	newbotFunc     newbotFunc
	conn           *twitch.Client
	mu             sync.Mutex
	log            logrus.FieldLogger
	logger         logrus.FieldLogger
	botsByUsername map[string]bot.Bot
}

type newbotFunc func(common.Channel, logrus.FieldLogger, *twitch.Client, *BotManager) bot.Bot

func New(logger logrus.FieldLogger, newbotFunc newbotFunc) *BotManager {
	m := &BotManager{
		newbotFunc:     newbotFunc,
		logger:         logger,
		log:            logger.WithField("package", "botmanager"),
		botsByUsername: make(map[string]bot.Bot),
	}
	return m
}

func (m *BotManager) Connect(username string, password string, ircAddr string, useTLS bool) {
	client := twitch.NewClient(username, password)
	client.IrcAddress = ircAddr
	client.TLS = useTLS
	client.SetupCmd = "LOGIN thonkbot"
	client.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		msg := &common.Message{
			ChannelName: channel,
			User:        user,
			Message:     message,
			Split:       strings.Split(message.Text, " "),
		}
		m.handleMessage(msg)
	})
	go client.Connect()
	m.conn = client
}

func (m *BotManager) JoinChannels() error {
	channels, err := mysql.FetchChannels()
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, channel := range channels {
		m.joinChannel(channel)
	}
	m.joinChannel(common.Channel{
		ID:   0,
		Name: config.Cfg.TwitchUsername,
	})
	return nil
}

func (m *BotManager) JoinChannel(channel common.Channel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.joinChannel(channel)
}

func (m *BotManager) joinChannel(channel common.Channel) {
	if _, ok := m.botsByUsername[channel.Name]; ok {
		m.log.WithFields(logrus.Fields{
			"channel-id":   channel.ID,
			"channel-name": channel.Name,
		}).Warn("Attempted to remake bot in a channel it's already in!")
		return
	}

	b := m.newbotFunc(channel, m.logger, m.conn, m)
	m.botsByUsername[channel.Name] = b
	m.conn.Join(channel.Name)
}

func (m *BotManager) LeaveChannel(channel string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.botsByUsername, channel)
	m.conn.Depart(channel)
}

func (m *BotManager) handleMessage(msg *common.Message) {
	bot, ok := m.botsByUsername[msg.ChannelName]
	if !ok {
		return
	}
	bot.HandleMessage(msg)
}

func (m *BotManager) ChannelCount() int {
	return len(m.botsByUsername) - 1
}
