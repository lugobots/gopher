package main

import (
	"log"

	"github.com/lugobots/lugo4go/v3"

	"github.com/lugobots/the-dummies-go/v3/bot"
)

func main() {
	connectionStarter, defaultFieldMapper, err := lugo4go.NewDefaultStarter()
	if err != nil {
		log.Fatalf("failed to load the bot configuration: %s", err)
	}

	//
	// Optional: define your own field mapper
	// defaultFieldMapper, err = mapper.NewMapper(NUM_COLS, NUM_ROWS, connectionStarter.Config.TeamSide)
	// if err != nil {
	// 	log.Fatalf("failed to create a field mapper: %s", err)
	// }

	// create your bot as you wish
	// in this example, the bot requires the field mapper, the connection config, and a logger.
	myBot := bot.NewBot(
		defaultFieldMapper,
		connectionStarter.Config,
		connectionStarter.Logger,
	)

	// then lets play
	if err := connectionStarter.Run(myBot); err != nil {
		log.Fatalf("bot stopped: %s", err)
	}
}
