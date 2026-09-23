#!/bin/bash
set -e

LOCATION='ssgo-docker'
DEST='root@calp'

echo "🔄 Triggering remote app load and restart..."
ssh ${DEST} "bash -c 'cd ~/${LOCATION} && \
  chmod +x start.sh load.sh restart.sh stop.sh && \
  ./restart.sh'"

echo "🎉 Deployment successful!"
