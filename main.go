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

	// OPTIONAL
	// define your own field mapper! The default number of col/rows are defined by lugo4go.DefaultFieldMapCols and lugo4go.DefaultFieldMapRows
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

	// Here you define the initial position of your bot. It's important to use the field mapper instead of points because
	// the field mapper won't be affected when your bot is playing on the away side
	initialPosition := bot.DefaultInitialPositions[connectionStarter.Config.Number]
	region, err := defaultFieldMapper.GetRegion(initialPosition.Col, initialPosition.Row)
	if err != nil {
		log.Fatalf("failed to define initial position using field mapper: %s", err)
	}
	connectionStarter.Config.InitialPosition = region.Center()

	// then lets play
	if err := connectionStarter.Run(myBot); err != nil {
		log.Fatalf("bot stopped: %s", err)
	}
}
