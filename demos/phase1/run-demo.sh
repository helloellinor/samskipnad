#!/bin/bash

# Phase 1 Demo Runner for Samskipnad Platform
# This script demonstrates the core platform functionality

set -e

echo "🚀 Starting Phase 1 Demo: Core Samskipnad Platform"
echo "=================================================="

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DEMO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "📁 Project root: $PROJECT_ROOT"
echo "📁 Demo directory: $DEMO_DIR"

# Change to project root
cd "$PROJECT_ROOT"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21 or later.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go is installed: $(go version)${NC}"

# Build the application
echo -e "${BLUE}🔨 Building the Samskipnad platform...${NC}"
make build

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Build failed. Please check the error messages above.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Build successful!${NC}"

# Set up demo environment
echo -e "${BLUE}🔧 Setting up demo environment...${NC}"

# Copy demo configuration
cp "$DEMO_DIR/demo-config.yaml" "$PROJECT_ROOT/community-config.yaml" 2>/dev/null || true

# Set demo environment variables
export COMMUNITY="demo"
export PORT="8080"

# Clean up any existing database for fresh demo
rm -f "$PROJECT_ROOT/samskipnad.db"

echo -e "${GREEN}✅ Demo environment configured${NC}"

# Start the server in background
echo -e "${BLUE}🌐 Starting the Samskipnad server...${NC}"

# Run the server and capture the PID
"$PROJECT_ROOT/bin/samskipnad" &
SERVER_PID=$!

# Function to cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}🧹 Cleaning up...${NC}"
    if kill -0 $SERVER_PID 2>/dev/null; then
        kill $SERVER_PID
        wait $SERVER_PID 2>/dev/null || true
    fi
    echo -e "${GREEN}✅ Demo cleanup complete${NC}"
}

# Set trap to cleanup on script exit
trap cleanup EXIT INT TERM

# Wait for server to start
echo -e "${BLUE}⏳ Waiting for server to start...${NC}"
sleep 3

# Check if server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo -e "${RED}❌ Server failed to start. Check the logs above.${NC}"
    exit 1
fi

# Test if server is responding
for i in {1..10}; do
    if curl -s http://localhost:8080 > /dev/null; then
        break
    fi
    if [ $i -eq 10 ]; then
        echo -e "${RED}❌ Server is not responding after 10 attempts${NC}"
        exit 1
    fi
    sleep 1
done

echo -e "${GREEN}✅ Server is running successfully!${NC}"

# Display demo information
echo ""
echo "🎉 Phase 1 Demo is now running!"
echo "======================================"
echo ""
echo -e "${GREEN}🌐 Web Interface:${NC} http://localhost:8080"
echo ""
echo -e "${BLUE}📝 Demo Credentials:${NC}"
echo "   Regular User:"
echo "   • Email: demo@example.com"
echo "   • Password: demo123"
echo ""
echo "   Admin User:"
echo "   • Email: admin@example.com"
echo "   • Password: admin123"
echo ""
echo -e "${BLUE}🎯 Things to Try:${NC}"
echo "   1. Visit the home page and explore the community"
echo "   2. Register a new account or login with demo credentials"
echo "   3. Browse the class calendar and book a class"
echo "   4. Try purchasing a membership or klippekort"
echo "   5. Login as admin and explore the admin interface"
echo "   6. Check out the dynamic styling and responsive design"
echo ""
echo -e "${BLUE}🔧 Technical Features:${NC}"
echo "   • Multi-tenant community configuration"
echo "   • Role-based authentication and authorization"
echo "   • Real-time HTMX interactions"
echo "   • Secure payment processing simulation"
echo "   • Mobile-responsive design"
echo "   • Administrative management interface"
echo ""
echo -e "${YELLOW}📚 For more details, see: demos/phase1/README.md${NC}"
echo ""
echo -e "${RED}Press Ctrl+C to stop the demo${NC}"

# Keep the script running until interrupted
while true; do
    sleep 1
    # Check if server is still running
    if ! kill -0 $SERVER_PID 2>/dev/null; then
        echo -e "${RED}❌ Server process died unexpectedly${NC}"
        exit 1
    fi
done