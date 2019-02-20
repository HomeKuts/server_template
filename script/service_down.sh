#!/bin/bash

SERVICE="cmd"

ps aux | grep -i $SERVICE | grep -v grep | awk {'print $2'} | xargs kill -2

