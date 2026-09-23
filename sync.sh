#!/bin/bash
set -e

LOCATION='ssgo-docker'
DEST='root@calp'

echo "=== Creating remote directories ==="
ssh $DEST "mkdir -p ~/$LOCATION"

echo "=== Syncing deployment scripts, .env & compose file ==="
rsync -rvv \
    .env \
    load.sh \
    restart.sh \
    start.sh \
    stop.sh \
    deployment/stack/docker-compose.yml \
    $DEST:~/$LOCATION/

echo "=== Syncing massive app.tar (Docker Image) ==="
rsync -rvv app.tar $DEST:~/$LOCATION/
