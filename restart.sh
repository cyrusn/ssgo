#!/bin/bash
# Ensure we are in the script's directory for crontab compatibility
cd "$(dirname "$0")"

./load.sh

docker compose down

./start.sh

docker container prune -f
docker image prune -a -f
