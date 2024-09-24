#!/bin/bash

export BOTTOKEN=""

export REDISADDR="localhost:6379"

export SUPPORTPASSWORD="1111111111"
export ITPASSWORD="2222222222"
export BILLINGPASSWORD="3333333333"

go run "./../cmd/main.go"
# docker-compose -f docker-compose.yml up -d --remove-orphans