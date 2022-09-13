# Lugo - The Dummies Go

[![GoDoc](https://godoc.org/github.com/lugobots/the-dummies-go?status.svg)](https://godoc.org/github.com/lugobots/the-dummies-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lugobots/the-dummies-go)](https://goreportcard.com/report/github.com/lugobots/the-dummies-go)

The Dummies Go is a [Go](http://golang.org/) implementation of a player (bot) for [Lugo](https://lugobots.dev) game.
This bot was made using the [Go Client Player](https://github.com/lugobots/lugo4go).

As this name suggest, _The Dummies_ are not that smart, but they may play well enough to help you to test your bot.

## Requirements

0. Docker >= 18.03 (https://docs.docker.com/install/)
0. Docker Compose >= 1.21 (https://docs.docker.com/compose/install/)

## Before starting

Are you familiar with Lugo? If not, before continuing, please visit [the project website](https://lugobots.dev) and read
about the game.

## How to use this source code

1. **Checkout the code** or download the most recent tag release
2. **Test it out**: Before any change, make The Dummies Go to play to ensure you are not working on a broken code.

   _Note_: this step can take a little longer at the first time.
   ```sh 
   docker-compose up
   ```
   and open [http://localhost:8080/](http://localhost:8080/) to watch the game.
4. **Now, make your changes**: (see :question:[How to change the bot](#how-to-edit-the-bot))
5. Play again to see your changes results:

   ```sh 
   docker-compose up
   ```
6. **Are you ready to compete? Build your Docker image:**

    ```sh 
   docker build -t my-super-bot .
   ```
7. :checkered_flag: Before pushing your changes

   ```sh 
   MY_BOT=my-super-bot docker-compose --file docker-compose-test.yml up
   ```

## How to edit the bot

The only files that you may need to edit are the ones inside [./bot](./bot) directory. Ignore all the other ones.

### Helper functions file [bot/helpers.go](./bot/helpers.go)

You may need to change this file if you want to change the player disposition in the field. E.g. if you want to change
the 4-4-2 disposition.

### Settings file [bot/settings.go](./bot/settings.go)

This file defines where the players will act based on the team state (attacking, defending, etc.) and also define the
field grid.

### Main file [bot/bot.go](./bot/bot.go)

:eyes: This is the most important file!

There will be 5 important methods that you must edit to change the bot behaviour.

```go

type Bot interface {
   // OnDisputing is called when no one has the ball possession
   OnDisputing(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
   
   // OnDefending is called when an opponent player has the ball possession
   OnDefending(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
   
   // OnHolding is called when this bot has the ball possession
   OnHolding(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
   
   // OnSupporting is called when a teammate player has the ball possession
   OnSupporting(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
   
   // AsGoalkeeper is only called when this bot is the goalkeeper (number 1). This method is called on every turn,
   // and the player state is passed at the last parameter.
   AsGoalkeeper(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot, state PlayerState) error
}
```
