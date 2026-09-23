#!/usr/bin/env bash
set -e

# ssgo-deployment automated build and packaging script
echo "=== Phase 1: Building Vue Frontend ==="
pushd ../ssgo-frontend > /dev/null
echo "Installing frontend dependencies..."
yarn install --silent
echo "Building Vue production site..."
yarn build
popd > /dev/null

echo "=== Phase 2: Copying Frontend Assets to Server Public ==="
rm -rf ./server/public
mkdir -p ./server/public
cp -R ../ssgo-frontend/site/* ./server/public/
echo "Frontend assets copied successfully."

echo "=== Phase 3: Building Go Backend Binary ==="
pushd ../ssgo-server > /dev/null
# Default target is Intel Linux (amd64) as most production Linux servers are x86_64.
# Can be overridden by running e.g. ARCH=arm64 ./build_deploy_bundle.sh
TARGET_ARCH="${ARCH:-amd64}"
echo "Building Go backend for Linux (${TARGET_ARCH})..."
GOOS=linux GOARCH=${TARGET_ARCH} go build -o ../deployment/server/ssgo main.go
echo "Backend binary built successfully for Linux (${TARGET_ARCH})."
popd > /dev/null

echo "=== Phase 4: Building Unified Docker Image ==="
pushd ./server > /dev/null
echo "Building Docker image 'ssgo-server:latest'..."
docker build -t ssgo-server:latest .
popd > /dev/null

echo "=== Phase 5: Creating Deployment Bundle ==="
BUNDLE_DIR="./ssgo-deployment"
rm -rf "${BUNDLE_DIR}"
mkdir -p "${BUNDLE_DIR}"
mkdir -p "${BUNDLE_DIR}/database"

echo "Saving Docker image to tar..."
docker save -o "${BUNDLE_DIR}/ssgo-server-image.tar" ssgo-server:latest

echo "Copying deployment files..."
cp ./stack/docker-compose.yml "${BUNDLE_DIR}/"
cp ../database/*.json "${BUNDLE_DIR}/database/"

echo "Writing startup/shutdown scripts..."
cat << 'EOF' > "${BUNDLE_DIR}/start.sh"
#!/usr/bin/env bash
set -e

# Ensure we are in the script's directory for crontab compatibility
cd "$(dirname "$0")"

echo "=========================================="
echo "Starting SSGO Subject Selection System"
echo "=========================================="

if [ -f "ssgo-server-image.tar" ]; then
    echo "Loading Docker image..."
    docker load -i ssgo-server-image.tar
    # Clean up image tar to save disk space on the target machine after loading
    rm ssgo-server-image.tar
fi

echo "Starting container..."
docker compose up -d

echo "------------------------------------------"
echo "Services are UP and RUNNING!"
echo "Port: 6612"
echo "Log: docker compose logs -f server"
echo "=========================================="
EOF

cat << 'EOF' > "${BUNDLE_DIR}/stop.sh"
#!/usr/bin/env bash
# Ensure we are in the script's directory for crontab compatibility
cd "$(dirname "$0")"

echo "Stopping SSGO services..."
docker compose down
echo "Services stopped."
EOF

chmod +x "${BUNDLE_DIR}/start.sh"
chmod +x "${BUNDLE_DIR}/stop.sh"

echo "=== Phase 6: Packaging Bundle into Tarball ==="
TAR_FILE="deploy-bundle.tar.gz"
rm -f "${TAR_FILE}"
tar -czf "${TAR_FILE}" "${BUNDLE_DIR}"

# Clean up staging directories
rm -rf "${BUNDLE_DIR}"
rm -rf ./server/public
rm -f ./server/ssgo

echo "=========================================="
echo "BUILD SUCCESSFUL!"
echo "Deployable tarball created: deployment/${TAR_FILE}"
echo "=========================================="
