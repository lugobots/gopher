#!/bin/bash

max_retry=10
counter=0
until /app/my_bot_binary
do
   sleep 3
   [[ counter -eq $max_retry ]] && echo "Failed!" && exit 1
   echo "Trying again. Try #$counter"
   ((counter++))
done
