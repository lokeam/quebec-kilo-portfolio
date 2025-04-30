#!/bin/bash

# Colors for better output readability
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print section headers
section() {
  echo -e "\n${BLUE}==== $1 ====${NC}"
}

# Function to print success messages
success() {
  echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error messages
error() {
  echo -e "${RED}✗ $1${NC}"
}

# Function to print info messages
info() {
  echo -e "${YELLOW}$1${NC}"
}

# Test network and DNS configuration
test_network() {
  section "Testing Network Configuration"

  # Check if localhost resolves to 127.0.0.1
  info "Checking 'localhost' resolution..."
  HOST_IP=$(ping -c 1 localhost | grep PING | awk -F'(' '{print $2}' | awk -F')' '{print $1}')
  if [ "$HOST_IP" = "127.0.0.1" ]; then
    success "localhost resolves to 127.0.0.1"
  else
    error "localhost resolves to $HOST_IP - this may cause connection issues"
    info "Consider adding 'localhost 127.0.0.1' to your /etc/hosts file"
  fi

  # Check if port 3000 is already in use
  info "Checking if port 3000 is in use..."
  if lsof -i :3000 -t &>/dev/null; then
    error "Port 3000 is already in use by:"
    lsof -i :3000
  else
    success "Port 3000 is available"
  fi

  # Test connectivity to 127.0.0.1
  info "Testing connection to 127.0.0.1:3000..."
  if curl -s -m 1 http://127.0.0.1:3000 &>/dev/null; then
    success "Connection to 127.0.0.1:3000 succeeded"
  else
    info "Connection to 127.0.0.1:3000 failed (this is expected if Vite is not running)"
  fi

  # Check hosts file
  info "Checking /etc/hosts file..."
  if grep -q "^127.0.0.1.*localhost" /etc/hosts; then
    success "localhost is properly configured in /etc/hosts"
  else
    error "localhost entry may be missing or incorrect in /etc/hosts"
  fi

  # Check network interfaces
  info "Available network interfaces:"
  ifconfig | grep -E 'inet (127|192|10|172)' | awk '{print $2}'
}

# Show help message
show_help() {
  echo -e "${GREEN}Vite Control Script${NC} - A single tool to manage your Vite development server"
  echo
  echo "Usage: ./vite-control.sh [command]"
  echo
  echo "Commands:"
  echo "  start       Clean restart Vite with optimal settings"
  echo "  stop        Stop all running Vite instances"
  echo "  diagnose    Run diagnostics without changing anything"
  echo "  clean       Clean Vite cache and temporary files"
  echo "  restore     Restore original configuration"
  echo "  network     Test network and DNS configuration"
  echo "  help        Show this help message"
  echo
  echo "Examples:"
  echo "  ./vite-control.sh start    # Most common usage - clean start Vite"
  echo "  ./vite-control.sh diagnose # Check what might be wrong"
  echo
  echo -e "${YELLOW}IMPORTANT NOTES:${NC}"
  echo "• Always use ${GREEN}http://127.0.0.1:3000${NC} instead of localhost:3000 to avoid DNS resolution issues"
  echo "• If connections are failing, try running ${GREEN}./vite-control.sh network${NC} to test your network config"
  echo "• If Vite still won't start, try running ${GREEN}./vite-control.sh clean${NC} before starting again"
  echo
}

# Kill all Vite-related processes
kill_vite() {
  section "Stopping Vite processes"

  # Kill vite processes
  VITE_PIDS=$(pgrep -f vite)
  if [ -n "$VITE_PIDS" ]; then
    info "Found Vite processes running:"
    ps -f $VITE_PIDS
    info "Killing Vite processes..."
    pkill -9 -f vite
    success "Vite processes killed"
  else
    success "No Vite processes running"
  fi

  # Kill esbuild processes (Vite dependency)
  ESBUILD_PIDS=$(pgrep -f esbuild)
  if [ -n "$ESBUILD_PIDS" ]; then
    info "Found esbuild processes running:"
    ps -f $ESBUILD_PIDS
    info "Killing esbuild processes..."
    pkill -9 -f esbuild
    success "esbuild processes killed"
  else
    success "No esbuild processes running"
  fi

  # Check for processes on port 3000
  PORT_PIDS=$(lsof -i :3000 -t 2>/dev/null)
  if [ -n "$PORT_PIDS" ]; then
    info "Found processes using port 3000:"
    ps -f $PORT_PIDS
    info "Killing processes on port 3000..."
    lsof -i :3000 -t 2>/dev/null | xargs kill -9 2>/dev/null
    success "Port 3000 freed"
  else
    success "Port 3000 is available"
  fi

  # Wait for processes to terminate
  sleep 1
}

# Clean Vite cache and temporary files
clean_vite() {
  section "Cleaning Vite"

  # Remove Vite cache
  if [ -d "node_modules/.vite-cache" ]; then
    info "Removing Vite cache..."
    rm -rf node_modules/.vite-cache
    success "Vite cache removed"
  else
    success "No Vite cache to remove"
  fi

  # Clean up dist directory if it exists
  if [ -d "dist" ]; then
    info "Removing dist directory..."
    rm -rf dist
    success "dist directory removed"
  fi

  # Clean up vite cache directories
  if [ -d ".vite" ]; then
    info "Removing .vite directory..."
    rm -rf .vite
    success ".vite directory removed"
  fi

  if [ -d "node_modules/.vite" ]; then
    info "Removing node_modules/.vite directory..."
    rm -rf node_modules/.vite
    success "node_modules/.vite directory removed"
  fi

  success "Temporary files cleaned"
}

# Run diagnostics
diagnose_vite() {
  section "VITE DIAGNOSTICS REPORT"

  # Check Node.js version
  section "Environment"
  NODE_VERSION=$(node -v)
  info "Node.js version: $NODE_VERSION"
  NPM_VERSION=$(npm -v)
  info "npm version: $NPM_VERSION"

  # Check for DNS and network issues
  section "Network Diagnostics"
  LOCALHOST_IP=$(ping -c 1 localhost | grep PING | awk -F'(' '{print $2}' | awk -F')' '{print $1}')
  info "localhost resolves to: $LOCALHOST_IP"

  if [ "$LOCALHOST_IP" != "127.0.0.1" ]; then
    error "DNS Resolution Issue: localhost should resolve to 127.0.0.1"
    info "This is likely causing your connection problems"
    echo
    echo "Try accessing your app using http://127.0.0.1:3000 directly instead of localhost:3000"
  fi

  # Check for running Vite processes
  section "Processes"
  VITE_PIDS=$(pgrep -f vite)
  if [ -n "$VITE_PIDS" ]; then
    error "Found Vite processes running:"
    ps -f $VITE_PIDS
  else
    success "No Vite processes running"
  fi

  # Network diagnostics
  section "Port Status"
  if lsof -i :3000 | grep -q LISTEN; then
    success "Port 3000 is in use by:"
    lsof -i :3000 | grep LISTEN

    # Attempt connection to verify service is responding
    if curl -s -m 2 http://127.0.0.1:3000 > /dev/null; then
      success "Vite server is responding at http://127.0.0.1:3000"
    else
      error "Vite server is running but not responding to requests"
      info "This could indicate a configuration issue or a stalled server"
      info "Try running ./vite-control.sh stop and then ./vite-control.sh clean"
    fi
  else
    error "Port 3000 is not in use - no server is listening"
  fi

  # Check file system
  section "Configuration"
  info "Checking Vite config files..."
  find . -maxdepth 1 -name "vite*.config.*" | while read -r config; do
    success "Found config: $config"
  done

  # Check cache directories
  info "Checking cache directories..."
  if [ -d "node_modules/.vite" ]; then
    info "Found node_modules/.vite cache directory"
    du -sh node_modules/.vite 2>/dev/null || info "Unable to check size"
  fi

  if [ -d "node_modules/.vite-cache" ]; then
    info "Found node_modules/.vite-cache directory"
    du -sh node_modules/.vite-cache 2>/dev/null || info "Unable to check size"
  fi

  # Check main source files
  info "Checking main source files..."
  if [ -f "src/main.tsx" ]; then
    success "Found main.tsx entry point"
  else
    error "main.tsx entry point not found"
  fi

  if [ -f "index.html" ]; then
    success "Found index.html"
  else
    error "index.html not found"
  fi

  # Check dependencies
  info "Checking dependencies..."
  if grep -q '"react"' package.json && grep -q '"vite"' package.json; then
    success "Core dependencies found in package.json"
  else
    error "Some core dependencies may be missing"
  fi

  info "Vite version: $(grep '"vite"' package.json | head -1 | awk -F: '{print $2}' | tr -d '", ')"

  section "Recommendation"
  echo "Based on diagnostics, try these steps:"
  echo "1. Always use ${GREEN}http://127.0.0.1:3000${NC} instead of localhost:3000"
  echo "2. Run ${GREEN}./vite-control.sh stop${NC} to kill any stuck processes"
  echo "3. Run ${GREEN}./vite-control.sh clean${NC} to clear caches"
  echo "4. Run ${GREEN}./vite-control.sh start${NC} to start fresh"
}

# Modify vite.base.config.ts for better reliability
patch_vite_config() {
  section "Patching Vite configuration for reliability"

  # Create a backup of the original config
  if [ ! -f "vite.base.config.ts.bak" ]; then
    info "Creating backup of original vite.base.config.ts..."
    cp vite.base.config.ts vite.base.config.ts.bak
    success "Backup created at vite.base.config.ts.bak"
  else
    info "Backup already exists, using it for reference"
  fi

  # Update the configuration to be more reliable
  info "Updating vite.base.config.ts with more reliable settings..."
  cat > vite.base.config.ts << EOL
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: [
      { find: '@', replacement: resolve(__dirname, './src') },
      { find: '@test', replacement: resolve(__dirname, './test') }
    ]
  },
  server: {
    host: '127.0.0.1', // Using 127.0.0.1 instead of 0.0.0.0 for better reliability
    port: 3000,
    strictPort: true, // Added this to fail if port is in use
    hmr: {
      protocol: 'ws',
      host: '127.0.0.1',
      port: 3000,
      overlay: false,
      clientPort: 3000, // Ensure consistent port for HMR
      timeout: 30000 // Increase timeout to 30 seconds
    },
    watch: {
      usePolling: false, // Set to true if you're having file watching issues
      interval: 100
    },
    proxy: {
      '/api': {
        target: 'http://localhost:80',
        changeOrigin: true,
        secure: false,
        headers: {
          'Host': 'api.localhost'
        },
        configure: (proxy) => {
          proxy.on('error', (err) => {
            console.log('proxy error', err);
          });
        }
      }
    }
  },
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      'react-router-dom',
      '@tanstack/react-query',
      '@auth0/auth0-react'
    ],
    force: true // Force dependency optimization
  },
  build: {
    sourcemap: true, // Add sourcemaps for better debugging
    chunkSizeWarningLimit: 1000 // Increase the limit for larger chunks
  },
  cacheDir: 'node_modules/.vite-cache'
});
EOL
  success "vite.base.config.ts updated with more reliable settings"
}

# Start Vite with the standard dev script but with optimizations
start_vite() {
  section "Starting Vite with reliable configuration"

  # First stop any running instances and clean up
  kill_vite
  clean_vite

  # Patch the Vite config
  patch_vite_config

  # Start Vite with optimized settings
  info "Starting Vite server..."
  info "IMPORTANT: Please use http://127.0.0.1:3000 to access your app"
  info "Press Ctrl+C to stop the server when finished"
  echo "----------------------------------------------"

  # Use the standard npm run dev command but with increased memory and specific settings
  NODE_OPTIONS="--max-old-space-size=4096" npm run dev:clean &

  # Add a PID capture to track our process
  VITE_PID=$!

  # Wait for server to be ready with retries
  info "Waiting for Vite server to be ready..."
  MAX_RETRIES=15
  RETRY_COUNT=0
  CONNECTED=false

  while [ $RETRY_COUNT -lt $MAX_RETRIES ] && [ "$CONNECTED" = false ]; do
    sleep 1
    if curl -s http://127.0.0.1:3000 > /dev/null 2>&1; then
      success "Vite server is ready at http://127.0.0.1:3000"
      CONNECTED=true
    else
      RETRY_COUNT=$((RETRY_COUNT + 1))
      info "Waiting for server to start... ($RETRY_COUNT/$MAX_RETRIES)"
    fi
  done

  if [ "$CONNECTED" = false ]; then
    error "Vite server did not start properly in the allocated time"
    info "You can try accessing http://127.0.0.1:3000 manually"
  else
    # Open the URL in the default browser on macOS
    if [[ "$OSTYPE" == "darwin"* ]]; then
      open http://127.0.0.1:3000
    fi
  fi

  # Bring the Vite process back to foreground
  wait $VITE_PID
}

# Restore original config if needed
restore_config() {
  section "Restoring original configuration"

  if [ -f "vite.base.config.ts.bak" ]; then
    info "Restoring vite.base.config.ts from backup..."
    cp vite.base.config.ts.bak vite.base.config.ts
    success "Original configuration restored"
  else
    error "No backup file found at vite.base.config.ts.bak"
  fi
}

# Main command handler
case "$1" in
  start)
    start_vite
    ;;
  stop)
    kill_vite
    success "All Vite processes stopped"
    ;;
  diagnose)
    diagnose_vite
    ;;
  clean)
    kill_vite
    clean_vite
    success "Vite cleaned successfully"
    ;;
  restore)
    restore_config
    ;;
  network)
    test_network
    ;;
  help|--help|-h)
    show_help
    ;;
  *)
    if [ -z "$1" ]; then
      error "No command specified"
      show_help
    else
      error "Unknown command: $1"
      show_help
    fi
    exit 1
    ;;
esac