#!/bin/bash
# Ensure we are in the script's directory for crontab compatibility
cd "$(dirname "$0")"

docker compose up -d
