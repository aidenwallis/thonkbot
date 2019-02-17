package common

import twitch "github.com/gempir/go-twitch-irc"

type Message struct {
	ChannelName string
	User        twitch.User
	Message     twitch.Message
	Split       []string
}
