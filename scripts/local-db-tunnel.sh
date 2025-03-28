#!/bin/bash

# Local Database Tunnel Script
# Creates a secure connection to the local development database

# Configuration
LOCAL_PORT=5434  # Using 5434 to avoid conflicts
DB_NAME=${POSTGRES_DB:-qkoapi}
DB_USER=${POSTGRES_USER:-postgres}

# Display connection information
echo "========================================"
echo " QKO API Local Database Connection"
echo "========================================"
echo ""
echo "When connected, use these settings in your database tool:"
echo " Host: localhost"
echo " Port: $LOCAL_PORT"
echo " User: $DB_USER"
echo " Database: $DB_NAME"
echo ""
echo "Press Ctrl+C to close when finished."
echo "========================================"

# Check if socat is installed
if ! command -v socat &> /dev/null; then
    echo "Error: socat is not installed. Please install it with:"
    echo "  brew install socat  # macOS"
    echo "  apt-get install socat  # Ubuntu/Debian"
    exit 1
fi

# Create the tunnel using socat
echo "Creating tunnel to local database container..."
socat TCP-LISTEN:$LOCAL_PORT,reuseaddr,fork TCP:localhost:5432 &
SOCAT_PID=$!

# Cleanup on exit
trap "kill $SOCAT_PID 2>/dev/null" EXIT

# Keep the script running until user presses Ctrl+C
echo "Tunnel active. Press Ctrl+C to close."
while true; do
    sleep 1
done