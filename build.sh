#!/bin/bash
set -e

echo "=== Building Vue Frontend ==="
pushd ./ssgo-frontend > /dev/null
yarn install --silent
yarn build
popd > /dev/null

echo "=== Copying Frontend Assets to Server Public ==="
rm -rf ./deployment/server/public
mkdir -p ./deployment/server/public
cp -R ./ssgo-frontend/site/* ./deployment/server/public/

echo "=== Building Go Backend Binary (Linux AMD64) ==="
pushd ./ssgo-server > /dev/null
GOOS=linux GOARCH=amd64 go build -o ../deployment/server/ssgo main.go
popd > /dev/null

echo "=== Building Docker Image ==="
pushd ./deployment/server > /dev/null
docker buildx build --platform linux/amd64 -t ssgo-server:latest .
popd > /dev/null

echo "=== Saving Docker Image to app.tar ==="
docker save -o app.tar ssgo-server:latest

echo "🎉 Build complete! app.tar created successfully."
