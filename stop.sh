#!/bin/bash
# Ensure we are in the script's directory for crontab compatibility
cd "$(dirname "$0")"

echo "Stopping SSGO services..."
docker compose down
echo "Services stopped."
