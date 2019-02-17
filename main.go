package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aidenwallis/thonkbot/bot/singlebot"
	"github.com/aidenwallis/thonkbot/botmanager"
	"github.com/aidenwallis/thonkbot/config"
	"github.com/aidenwallis/thonkbot/mysql"
	"github.com/aidenwallis/thonkbot/web"
	"github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig()

	logger := logrus.StandardLogger()
	logger.SetLevel(logrus.DebugLevel)

	err := mysql.Connect(config.Cfg.MysqlDSN)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to mysql")
		os.Exit(1)
	}

	manager := botmanager.New(logger, singlebot.New)
	manager.Connect(config.Cfg.TwitchUsername, config.Cfg.TwitchPassword, config.Cfg.Twitchrelay, config.Cfg.UseTLS)
	err = manager.JoinChannels()
	if err != nil {
		logger.WithError(err).Error("Failed to fetch channels")
		os.Exit(1)
	}

	w := web.New(config.Cfg.Webhost, manager, logger)
	go w.Start()

	logrus.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	mysql.Close()
}
