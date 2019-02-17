package bot

import (
	"github.com/aidenwallis/thonkbot/common"
	"github.com/sirupsen/logrus"
)

type Bot interface {
	HandleMessage(*common.Message)
	JoinChannel(common.Channel)
	LeaveChannel(string)
	Channel() *common.Channel
	Close()
	Log() logrus.FieldLogger
	Sayf(string, ...interface{})
	ChannelCount() int
}
