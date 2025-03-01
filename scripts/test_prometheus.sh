#!/bin/bash

# This script starts Prometheus for monitoring the PocketBase application

# Default variables
PROMETHEUS_VERSION="2.44.0"
PROMETHEUS_PORT=9090
# Use environment variable METRICS_PORT if set, otherwise default to 9091
METRICS_PORT=${METRICS_PORT:-9091}
DASHBOARD_PORT=8090
PROMETHEUS_DIR="./prometheus"

# Check if Prometheus is already downloaded
if [ ! -f "${PROMETHEUS_DIR}/prometheus" ]; then
    echo "Downloading Prometheus..."
    mkdir -p "${PROMETHEUS_DIR}"
    
    # Determine platform
    PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
    if [ "$PLATFORM" == "darwin" ]; then
        PLATFORM="darwin"
    elif [ "$PLATFORM" == "linux" ]; then
        PLATFORM="linux"
    else
        echo "Unsupported platform. Please download Prometheus manually."
        exit 1
    fi
    
    # Determine architecture
    ARCH=$(uname -m)
    if [ "$ARCH" == "x86_64" ]; then
        ARCH="amd64"
    elif [ "$ARCH" == "arm64" ]; then
        ARCH="arm64"
    else
        echo "Unsupported architecture. Please download Prometheus manually."
        exit 1
    fi
    
    # Download Prometheus
    TEMP_DIR=$(mktemp -d)
    PROMETHEUS_URL="https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.${PLATFORM}-${ARCH}.tar.gz"
    
    echo "Downloading from: ${PROMETHEUS_URL}"
    curl -L "${PROMETHEUS_URL}" -o "${TEMP_DIR}/prometheus.tar.gz"
    
    echo "Extracting Prometheus..."
    tar -xzf "${TEMP_DIR}/prometheus.tar.gz" -C "${TEMP_DIR}"
    
    # Find the extracted directory (should be the only dir in temp)
    EXTRACTED_DIR=$(find "${TEMP_DIR}" -maxdepth 1 -type d -name 'prometheus*' | head -1)
    
    if [ -z "${EXTRACTED_DIR}" ]; then
        echo "Failed to extract Prometheus. Please check the downloaded file."
        exit 1
    fi
    
    # Copy necessary files to our prometheus directory
    cp "${EXTRACTED_DIR}/prometheus" "${PROMETHEUS_DIR}/"
    cp "${EXTRACTED_DIR}/promtool" "${PROMETHEUS_DIR}/" 2>/dev/null || true
    cp -r "${EXTRACTED_DIR}/console_libraries" "${PROMETHEUS_DIR}/" 2>/dev/null || true
    cp -r "${EXTRACTED_DIR}/consoles" "${PROMETHEUS_DIR}/" 2>/dev/null || true
    
    # Clean up
    rm -rf "${TEMP_DIR}"
    
    # Make sure it's executable
    chmod +x "${PROMETHEUS_DIR}/prometheus"
    
    echo "Prometheus successfully installed to ${PROMETHEUS_DIR}"
fi

# Create Prometheus configuration file
cat > "${PROMETHEUS_DIR}/prometheus.yml" << EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'pocketbase'
    static_configs:
      - targets: ['localhost:${METRICS_PORT}']
EOF

echo "Starting Prometheus on port ${PROMETHEUS_PORT}..."
echo "Configured to scrape PocketBase metrics from port ${METRICS_PORT}..."
cd "${PROMETHEUS_DIR}" && ./prometheus --config.file=prometheus.yml --web.listen-address=":${PROMETHEUS_PORT}" &
PROMETHEUS_PID=$!

if [ $? -ne 0 ]; then
    echo "Failed to start Prometheus. Please check if the binary exists and is executable."
    exit 1
fi

echo "Prometheus started with PID: $PROMETHEUS_PID"
echo "Prometheus UI available at: http://localhost:${PROMETHEUS_PORT}"
echo "PocketBase metrics available at: http://localhost:${METRICS_PORT}/metrics"
echo ""
echo "Press Ctrl+C to stop Prometheus"

# Wait for user to press Ctrl+C
trap "kill $PROMETHEUS_PID 2>/dev/null && echo 'Prometheus stopped'" INT
wait $PROMETHEUS_PID 