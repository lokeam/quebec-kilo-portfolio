#!/bin/bash

# Database SSH Tunnel Script
# Creates a secure tunnel to the production database for GUI tool access
#
# Usage:
#   ./scripts/db-tunnel.sh
#
# Before using, set these environment variables or edit the defaults below:
#   PROD_SSH_USER - Your SSH username for the production server
#   PROD_SSH_HOST - Your production server hostname or IP

# Configuration
REMOTE_USER="${PROD_SSH_USER:-your-default-username}"
REMOTE_HOST="${PROD_SSH_HOST:-your-default-host}"
LOCAL_PORT=5433
REMOTE_PORT=5432
DB_NAME=${POSTGRES_DB:-qkoapi}
DB_USER=${POSTGRES_USER:-postgres}

# Display connection information
echo "========================================"
echo " QKO API Database SSH Tunnel"
echo "========================================"
echo "Creating secure tunnel to qko production database..."
echo ""
echo "When connected, use these settings in your database tool:"
echo " Host: localhost"
echo " Port: $LOCAL_PORT"
echo " User: $DB_USER"
echo " Database: $DB_NAME"
echo ""
echo "Press Ctrl+C to close the tunnel when finished."
echo "========================================"

# Test SSH connection
echo "Testing connection..."
if ! ssh -q -o "BatchMode=yes" $REMOTE_USER@$REMOTE_HOST "exit" 2>/dev/null; then
  echo "Error: Cannot connect to $REMOTE_HOST. Please check your SSH credentials."
  exit 1
fi

# Create SSH tunnel
ssh -L $LOCAL_PORT:localhost:$REMOTE_PORT $REMOTE_USER@$REMOTE_HOST

# Check for errors
if [ $? -ne 0 ]; then
  echo "Error: SSH connection failed or was terminated unexpectedly."
  exit 1
fi

# This will only execute after the SSH connection is closed
echo ""
echo "SSH tunnel closed. Database connection terminated."
