#!/usr/bin/env bash
set -e

echo "=========================================="
echo "Starting SSGO Local Server Setup"
echo "=========================================="

# 1. Ensure database folder exists
if [ ! -d "./database" ]; then
    echo "Creating empty database folder..."
    mkdir -p ./database
fi

# Clean up any potential duplicate database folders in subdirectories
rm -rf ./ssgo-server/database ./deployment/server/database

# 2. Build the Vue frontend
echo "Step 1: Building Vue frontend..."
pushd ./ssgo-frontend > /dev/null
yarn install --silent
NODE_ENV=development yarn build
popd > /dev/null

# 3. Copy frontend assets to local Go public folder
echo "Step 2: Copying frontend assets to ssgo-server/public/..."
rm -rf ./ssgo-server/public
mkdir -p ./ssgo-server/public
cp -R ./ssgo-frontend/site/* ./ssgo-server/public/

# 4. Building Go Server (Verify build)
echo "Step 3: Building Go backend..."
pushd ./ssgo-server > /dev/null
go build -o ssgo main.go
popd > /dev/null

# 5. Run the server
echo "Step 4: Launching local server on port 3000..."
echo "------------------------------------------"
echo "Setup complete! Open your browser at:"
echo "👉 http://localhost:3000"
echo ""
echo "Onboarding Instructions (If first time):"
echo "1. Go to: http://localhost:3000/#/config"
echo "2. Login with Username: 'root' and Password found in your .env file"
echo "3. Create a cohort (e.g. '2026-27_Mock'), activate it, and paste sample data!"
echo "------------------------------------------"
echo "Press Ctrl+C to stop the server."
echo "=========================================="

pushd ./ssgo-server > /dev/null
./ssgo serve --port :3000
popd > /dev/null
