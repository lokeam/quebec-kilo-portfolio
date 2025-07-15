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

  # Check if port 4321 is already in use (Astro's default port)
  info "Checking if port 4321 is in use..."
  if lsof -i :4321 -t &>/dev/null; then
    error "Port 4321 is already in use by:"
    lsof -i :4321
  else
    success "Port 4321 is available"
  fi

  # Test connectivity to 127.0.0.1
  info "Testing connection to 127.0.0.1:4321..."
  if curl -s -m 1 http://127.0.0.1:4321 &>/dev/null; then
    success "Connection to 127.0.0.1:4321 succeeded"
  else
    info "Connection to 127.0.0.1:4321 failed (this is expected if Astro is not running)"
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
  echo -e "${GREEN}Astro Marketing Site Control Script${NC} - A single tool to manage your Astro development server"
  echo
  echo "Usage: ./vite-control.sh [command]"
  echo
  echo "Commands:"
  echo "  start       Clean restart Astro with optimal settings"
  echo "  stop        Stop all running Astro/Vite instances"
  echo "  diagnose    Run diagnostics without changing anything"
  echo "  clean       Clean Astro/Vite cache and temporary files"
  echo "  restore     Restore original configuration"
  echo "  network     Test network and DNS configuration"
  echo "  help        Show this help message"
  echo
  echo "Examples:"
  echo "  ./vite-control.sh start    # Most common usage - clean start Astro"
  echo "  ./vite-control.sh diagnose # Check what might be wrong"
  echo
  echo -e "${YELLOW}IMPORTANT NOTES:${NC}"
  echo "• Always use ${GREEN}http://127.0.0.1:4321${NC} instead of localhost:4321 to avoid DNS resolution issues"
  echo "• If connections are failing, try running ${GREEN}./vite-control.sh network${NC} to test your network config"
  echo "• If Astro still won't start, try running ${GREEN}./vite-control.sh clean${NC} before starting again"
  echo "• This script disables caching for immediate development feedback"
  echo
}

# Kill all Astro/Vite-related processes
kill_astro() {
  section "Stopping Astro/Vite processes"

  # Kill astro processes
  ASTRO_PIDS=$(pgrep -f astro)
  if [ -n "$ASTRO_PIDS" ]; then
    info "Found Astro processes running:"
    ps -f $ASTRO_PIDS
    info "Killing Astro processes..."
    pkill -9 -f astro
    success "Astro processes killed"
  else
    success "No Astro processes running"
  fi

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

  # Check for processes on port 4321 (Astro default)
  PORT_PIDS=$(lsof -i :4321 -t 2>/dev/null)
  if [ -n "$PORT_PIDS" ]; then
    info "Found processes using port 4321:"
    ps -f $PORT_PIDS
    info "Killing processes on port 4321..."
    lsof -i :4321 -t 2>/dev/null | xargs kill -9 2>/dev/null
    success "Port 4321 freed"
  else
    success "Port 4321 is available"
  fi

  # Wait for processes to terminate
  sleep 1
}

# Clean Astro/Vite cache and temporary files
clean_astro() {
  section "Cleaning Astro/Vite"

  # Remove Astro cache
  if [ -d ".astro" ]; then
    info "Removing Astro cache..."
    rm -rf .astro
    success "Astro cache removed"
  else
    success "No Astro cache to remove"
  fi

  # Remove Vite cache
  if [ -d "node_modules/.vite" ]; then
    info "Removing Vite cache..."
    rm -rf node_modules/.vite
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

  success "Temporary files cleaned"
}

# Run diagnostics
diagnose_astro() {
  section "ASTRO DIAGNOSTICS REPORT"

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
    echo "Try accessing your app using http://127.0.0.1:4321 directly instead of localhost:4321"
  fi

  # Check for running Astro processes
  section "Processes"
  ASTRO_PIDS=$(pgrep -f astro)
  if [ -n "$ASTRO_PIDS" ]; then
    error "Found Astro processes running:"
    ps -f $ASTRO_PIDS
  else
    success "No Astro processes running"
  fi

  # Network diagnostics
  section "Port Status"
  if lsof -i :4321 | grep -q LISTEN; then
    success "Port 4321 is in use by:"
    lsof -i :4321 | grep LISTEN

    # Attempt connection to verify service is responding
    if curl -s -m 2 http://127.0.0.1:4321 > /dev/null; then
      success "Astro server is responding at http://127.0.0.1:4321"
    else
      error "Astro server is running but not responding to requests"
      info "This could indicate a configuration issue or a stalled server"
      info "Try running ./vite-control.sh stop and then ./vite-control.sh clean"
    fi
  else
    error "Port 4321 is not in use - no server is listening"
  fi

  # Check file system
  section "Configuration"
  info "Checking Astro config files..."
  find . -maxdepth 1 -name "astro.config.*" | while read -r config; do
    success "Found config: $config"
  done

  # Check cache directories
  info "Checking cache directories..."
  if [ -d ".astro" ]; then
    info "Found .astro cache directory"
    du -sh .astro 2>/dev/null || info "Unable to check size"
  fi

  if [ -d "node_modules/.vite" ]; then
    info "Found node_modules/.vite cache directory"
    du -sh node_modules/.vite 2>/dev/null || info "Unable to check size"
  fi

  # Check main source files
  info "Checking main source files..."
  if [ -d "src" ]; then
    success "Found src directory"
  else
    error "src directory not found"
  fi

  if [ -f "src/pages/index.astro" ]; then
    success "Found index.astro entry point"
  else
    error "index.astro entry point not found"
  fi

  # Check dependencies
  info "Checking dependencies..."
  if grep -q '"astro"' package.json; then
    success "Astro dependency found in package.json"
  else
    error "Astro dependency may be missing"
  fi

  info "Astro version: $(grep '"astro"' package.json | head -1 | awk -F: '{print $2}' | tr -d '", ')"

  section "Recommendation"
  echo "Based on diagnostics, try these steps:"
  echo "1. Always use ${GREEN}http://127.0.0.1:4321${NC} instead of localhost:4321"
  echo "2. Run ${GREEN}./vite-control.sh stop${NC} to kill any stuck processes"
  echo "3. Run ${GREEN}./vite-control.sh clean${NC} to clear caches"
  echo "4. Run ${GREEN}./vite-control.sh start${NC} to start fresh"
}

# Modify astro.config.mjs for development (disable caching)
patch_astro_config() {
  section "Patching Astro configuration for development"

  # Create a backup of the original config
  if [ ! -f "astro.config.mjs.bak" ]; then
    info "Creating backup of original astro.config.mjs..."
    cp astro.config.mjs astro.config.mjs.bak
    success "Backup created at astro.config.mjs.bak"
  else
    info "Backup already exists, using it for reference"
  fi

  # Update the configuration to disable caching for development
  info "Updating astro.config.mjs with development settings..."
  cat > astro.config.mjs << EOL
// @ts-check
import { defineConfig } from 'astro/config';
import tailwind from '@astrojs/tailwind';

// https://astro.build/config
export default defineConfig({
  integrations: [
    tailwind(),
  ],
  // Disable build caching
  build: {
    inlineStylesheets: 'auto',
  },
  // Disable Vite caching and force reload
  vite: {
    // Force reload on file changes
    server: {
      hmr: {
        overlay: true,
      },
    },
    // Disable caching for immediate development feedback
    optimizeDeps: {
      force: true,
    },
  },
});
EOL
  success "astro.config.mjs updated with development settings (caching disabled)"
}

# Start Astro with the standard dev script but with optimizations
start_astro() {
  section "Starting Astro with reliable configuration"

  # First stop any running instances and clean up
  kill_astro
  clean_astro

  # Patch the Astro config
  patch_astro_config

  # Start Astro with optimized settings
  info "Starting Astro server..."
  info "IMPORTANT: Please use http://127.0.0.1:4321 to access your app"
  info "Press Ctrl+C to stop the server when finished"
  echo "----------------------------------------------"

  # Use the standard npm run dev command but with increased memory and specific settings
  NODE_OPTIONS="--max-old-space-size=4096" npm run dev &

  # Add a PID capture to track our process
  ASTRO_PID=$!

  # Wait for server to be ready with retries
  info "Waiting for Astro server to be ready..."
  MAX_RETRIES=15
  RETRY_COUNT=0
  CONNECTED=false

  while [ $RETRY_COUNT -lt $MAX_RETRIES ] && [ "$CONNECTED" = false ]; do
    sleep 1
    if curl -s http://127.0.0.1:4321 > /dev/null 2>&1 || curl -s http://localhost:4321 > /dev/null 2>&1; then
      success "Astro server is ready at http://localhost:4321 or http://127.0.0.1:4321"
      CONNECTED=true
    else
      RETRY_COUNT=$((RETRY_COUNT + 1))
      info "Waiting for server to start... ($RETRY_COUNT/$MAX_RETRIES)"
    fi
  done

  if [ "$CONNECTED" = false ]; then
    error "Astro server did not start properly in the allocated time"
    info "You can try accessing http://localhost:4321 or http://127.0.0.1:4321 manually"
  else
    # Open the URL in the default browser on macOS
    if [[ "$OSTYPE" == "darwin"* ]]; then
      open http://localhost:4321
    fi
  fi

  # Bring the Astro process back to foreground
  wait $ASTRO_PID
}

# Restore original config if needed
restore_config() {
  section "Restoring original configuration"

  if [ -f "astro.config.mjs.bak" ]; then
    info "Restoring astro.config.mjs from backup..."
    cp astro.config.mjs.bak astro.config.mjs
    success "Original configuration restored"
  else
    error "No backup file found at astro.config.mjs.bak"
  fi
}

# Main command handler
case "$1" in
  start)
    start_astro
    ;;
  stop)
    kill_astro
    success "All Astro/Vite processes stopped"
    ;;
  diagnose)
    diagnose_astro
    ;;
  clean)
    kill_astro
    clean_astro
    success "Astro/Vite cleaned successfully"
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