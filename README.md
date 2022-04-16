# Lugo - The Dummies Go

[![GoDoc](https://godoc.org/github.com/lugobots/the-dummies-go?status.svg)](https://godoc.org/github.com/lugobots/the-dummies-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lugobots/the-dummies-go)](https://goreportcard.com/report/github.com/lugobots/the-dummies-go)

The Dummies Go is a [Go](http://golang.org/) implementation of a player (bot) for [Lugo](https://lugobots.dev) game.
This bot was made using the [Go Client Player](https://github.com/lugobots/client-player-go).

As this name suggest, _The Dummies_ are not that smart, but they may play well enough to help you to test your bot.

### Requirements

0. Docker >= 18.03 (https://docs.docker.com/install/)
0. Docker Compose >= 1.21 (https://docs.docker.com/compose/install/)
0. Go Lang >= 1.17 (https://golang.org/doc/install)

### Usage 

There are several ways to run _The Dummies_, the 3 easiest ones are described below.
 
#### Option A - Running them in containers (no Git Clone needed)

1. Start the game server
```shell
# starts the server 
docker run -p 8080:8080 -p 5000:5000 lugobots/server:v1.0.0-beta play

```
2. Start the home and away teams.
```shell
# start the home team
docker run --net=host lugobots/the-dummies-go:v1.0.0-beta -team=[home away] -number=[1 ... 11]

```

:dart: Each team will require 11 players, so that command needs to be run 11 times. Or you can use this command:

```shell
./start-team-container.sh lugobots/the-dummies-go:v1.0.0-beta away

# (optional) if do not have your bot yet, you can run The Dummies as the opponent team, so you can watch the match  
./start-team-container.sh lugobots/the-dummies-go:v1.0.0-beta home
```

#### Option B - Running in containers after cloning the repo (requires Docker Compose)

If you clonned the repo and you are changing the code to build your bot, you can still use the local files to run the original _The Dummies_.


```
AWAY_TEAM=lugobots/the-dummies-go:v1.0.0-beta \
SERVER_VERSION=v1.0.0-beta \
docker-compose up
```

That command will run _The _Dummies_ as the **away** team (defined by the env variable `AWAY_TEAM`).
Now you may start **your** bot to play against _The Dummies_.

#### Option C - Running the processes directly on your machine (recommended for developing environment because the startup is a faster)

If you are working in your bot, and you want to play against _The Dummies_ several times to test your bot, so I recommend
you having a copy of _The Dummies_ in you machine because the bots will start up faster than running them as containers. 

1. Clone the repository to your machine
2. Start the game server
   ```
   docker run -p 8080:8080 -p 5000:5000 lugobots/server:v1.0.0-beta play
   ```

and then, you may execute the script `./play.sh [home|away]` in that directory when you want to start the team.

