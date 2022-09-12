#!/bin/bash

max_retry=10
counter=0

FILE=/app/my_bot_binary

# let's wait just a bit to ensure the old binary was already deleted
sleep 0.3

if [ ! -f "$FILE" ]; then
  sleep 0.5
fi

until [[ -f "$FILE" ]]
do
   sleep 0.5
   [[ counter -eq $max_retry ]] && echo "Failed!" && exit 1
   echo "Waiting bot be compiled again. Try #$counter"
   ((counter++))
done

/app/my_bot_binary
