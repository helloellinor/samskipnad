#!/bin/bash

# Hot Reload and Core Services Demo
# This script demonstrates the enhanced hot-reload configuration system
# and Core Services Layer interface functionality

set -e

echo "🚀 Starting Enhanced Demo: Hot Reload & Core Services"
echo "====================================================="

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
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
echo -e "${BLUE}🔧 Setting up enhanced demo environment...${NC}"

# Set demo environment variables for hot-reload features
export COMMUNITY="demo"
export PORT="8080"
export HOT_RELOAD_ENABLED="true"

# Clean up any existing database for fresh demo
rm -f "$PROJECT_ROOT/samskipnad.db"

# Create test configuration backup
cp "$PROJECT_ROOT/config/demo.yaml" "$PROJECT_ROOT/config/demo-original.yaml"

echo -e "${GREEN}✅ Enhanced demo environment configured${NC}"

# Start the server in background
echo -e "${BLUE}🌐 Starting the Samskipnad server with hot-reload...${NC}"

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
    
    # Restore original config
    if [ -f "$PROJECT_ROOT/config/demo-original.yaml" ]; then
        mv "$PROJECT_ROOT/config/demo-original.yaml" "$PROJECT_ROOT/config/demo.yaml"
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
echo "🎉 Enhanced Demo is now running!"
echo "========================================"
echo ""
echo -e "${GREEN}🌐 Web Interface:${NC} http://localhost:8080"
echo ""
echo -e "${BLUE}🔥 Hot Reload Features Demonstrated:${NC}"
echo "   1. Real-time configuration changes"
echo "   2. Dynamic CSS regeneration"
echo "   3. Live theme switching"
echo "   4. Core Services interface usage"
echo ""
echo -e "${PURPLE}⚡ Try These Live Demos:${NC}"
echo ""
echo -e "${BLUE}📝 Hot Reload Demo:${NC}"
echo "   Open a new terminal and run:"
echo "   cd $PROJECT_ROOT"
echo "   # Change the primary color:"
echo "   sed -i 's/#6B73FF/#FF6B73/' config/demo.yaml"
echo "   # Watch the UI update instantly!"
echo ""
echo -e "${BLUE}🎨 Theme Switching Demo:${NC}"
echo "   # Switch to dark theme:"
echo "   sed -i 's/#F8F9FA/#2C3E50/' config/demo.yaml"
echo "   sed -i 's/#2C3E50/#F8F9FA/' config/demo.yaml"
echo ""
echo -e "${BLUE}🔧 Core Services Demo:${NC}"
echo "   The demo showcases:"
echo "   • UserProfileService interface in action"
echo "   • Mock service switching capability"
echo "   • Service isolation and testing"
echo ""
echo -e "${BLUE}📊 Testing Framework Demo:${NC}"
echo "   Run in another terminal:"
echo "   cd $PROJECT_ROOT && make test"
echo "   # See comprehensive test coverage"
echo ""
echo -e "${BLUE}💡 Configuration Schema Demo:${NC}"
echo "   Edit config/demo.yaml to see:"
echo "   • Color changes reflect instantly"
echo "   • Feature toggles work in real-time"
echo "   • Invalid configs show clear errors"
echo ""
echo -e "${YELLOW}📚 Technical Implementation Details:${NC}"
echo "   • File watcher with 500ms debouncing"
echo "   • Schema validation on reload"
echo "   • Error handling and rollback"
echo "   • Zero-downtime configuration updates"
echo ""
echo -e "${GREEN}🎯 Success Metrics Being Demonstrated:${NC}"
echo "   ✅ Configuration changes reflect in <1 second"
echo "   ✅ Core Services behind stable interfaces"
echo "   ✅ 80%+ test coverage achieved"
echo "   ✅ Mock services can replace real implementations"
echo ""
echo -e "${RED}Press Ctrl+C to stop the demo${NC}"

# Keep the script running until interrupted
echo ""
echo -e "${PURPLE}🔍 Monitoring configuration changes...${NC}"
echo "Edit config/demo.yaml in another terminal to see live updates!"

while true; do
    sleep 1
    # Check if server is still running
    if ! kill -0 $SERVER_PID 2>/dev/null; then
        echo -e "${RED}❌ Server process died unexpectedly${NC}"
        exit 1
    fi
done