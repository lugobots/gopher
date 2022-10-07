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

Are you familiar with Lugo?
If not, before continuing, please visit [the project website](https://lugobots.dev) and read about the game.

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

