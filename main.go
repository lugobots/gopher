package main

import (
	"log"

	clientGo "github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/pkg/util"

	"github.com/lugobots/the-dummies-go/v2/bot"
)

func main() {
	// DefaultInitBundle is a shortcut for stuff that usually we define in init functions
	playerConfig, logger, err := util.DefaultInitBundle()
	if err != nil {
		log.Fatalf("could not init default config or logger: %s", err)
	}

	dummy := bot.NewBot(logger, playerConfig.TeamSide, playerConfig.Number)

	playerConfig.InitialPosition = dummy.MyInitialPosition()

	log.Printf("%s-%d: %v", playerConfig.TeamSide, playerConfig.Number, dummy.MyInitialPosition())

	player, err := clientGo.NewClient(playerConfig)
	if err != nil {
		log.Fatalf("could not init the gRPC client (server at %s): %s", playerConfig.GRPCAddress, err)
	}
	logger.Info("connected to the game server")

	if err := player.PlayAsBot(dummy, logger.Named("bot")); err != nil {
		logger.With("error", err).Warnf("got interruption signal")
	}

	logger.Infof("process finished")
}
