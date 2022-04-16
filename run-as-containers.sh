#!/usr/bin/env bash

if [ -z "$1" ]
  then
    echo "Please, pass the first argument (home or away) to set the team side"
    exit 1
fi

docker build -t my-bots .  || { echo "building has failed"; exit 1; }
for i in `seq 1 11`
do
  docker run --net=host my-bots -team=$1 -number=$i &
  sleep 0.1
done

echo ""
