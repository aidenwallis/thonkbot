package modules

import (
	"github.com/aidenwallis/thonkbot/bot"
	"github.com/aidenwallis/thonkbot/common"
)

type Module interface {
	Init(bot.Bot)
	Run(*common.Message)
}

type commandFunc func(*common.Message)
