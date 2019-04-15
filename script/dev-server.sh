#!/usr/bin/env bash
redis-server &
export REDIS_ADDRESS="localhost:6379"
export REDIS_PASSWORD=""
export PORT="8001"
gin -p 8000 run main.go