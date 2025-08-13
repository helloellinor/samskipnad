#!/bin/bash

# Multi-Community Demo Runner for Samskipnad Platform
# This script demonstrates multiple community configurations

set -e

echo "🚀 Starting Multi-Community Demo: Yoga Studio & Hackerspace"
echo "=========================================================="

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DEMO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "📁 Project root: $PROJECT_ROOT"
echo "📁 Demo directory: $DEMO_DIR"

# Parse command line arguments
COMMUNITY_TYPE=""
if [ "$1" = "yoga" ] || [ "$1" = "yoga-studio" ]; then
    COMMUNITY_TYPE="yoga-studio"
elif [ "$1" = "hack" ] || [ "$1" = "hackerspace" ]; then
    COMMUNITY_TYPE="hackerspace"
elif [ "$1" = "help" ] || [ "$1" = "--help" ]; then
    echo ""
    echo "Usage: $0 [community-type]"
    echo ""
    echo "Community types:"
    echo "  yoga, yoga-studio    - Start Zen Flow Yoga Studio"
    echo "  hack, hackerspace    - Start Oslo Hackerspace"
    echo "  [no argument]        - Interactive mode to choose"
    echo ""
    exit 0
fi

# Interactive community selection if not specified
if [ -z "$COMMUNITY_TYPE" ]; then
    echo ""
    echo "🎯 Choose a community to demo:"
    echo "  1) 🧘 Zen Flow Yoga Studio"
    echo "  2) 💻 Oslo Hackerspace"
    echo ""
    read -p "Enter your choice (1 or 2): " choice
    
    case $choice in
        1)
            COMMUNITY_TYPE="yoga-studio"
            ;;
        2)
            COMMUNITY_TYPE="hackerspace"
            ;;
        *)
            echo -e "${RED}❌ Invalid choice. Please run again and select 1 or 2.${NC}"
            exit 1
            ;;
    esac
fi

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
echo -e "${BLUE}🔧 Setting up multi-community demo environment...${NC}"

# Copy community configurations to the config directory (if needed)
cp "$PROJECT_ROOT/config/yoga-studio.yaml" "$PROJECT_ROOT/config/" 2>/dev/null || echo "Yoga studio config already in place"
cp "$PROJECT_ROOT/config/hackerspace.yaml" "$PROJECT_ROOT/config/" 2>/dev/null || echo "Hackerspace config already in place"

# Set demo environment variables
export COMMUNITY="$COMMUNITY_TYPE"
export PORT="8080"
export HOT_RELOAD_ENABLED="true"
export LOG_LEVEL="info"

# Clean up any existing database for fresh demo
rm -f "$PROJECT_ROOT/samskipnad.db"

echo -e "${GREEN}✅ Multi-community demo environment configured for: $COMMUNITY_TYPE${NC}"

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

# Display demo information based on community type
echo ""
if [ "$COMMUNITY_TYPE" = "yoga-studio" ]; then
    echo "🎉 Zen Flow Yoga Studio Demo is now running!"
    echo "=============================================="
    echo ""
    echo -e "${PURPLE}🧘 Zen Flow Yoga Studio${NC}"
    echo -e "   ${GREEN}🌐 URL:${NC} http://localhost:8080"
    echo -e "   ${BLUE}👤 Student:${NC} student@zenflow.example.com / namaste123"
    echo -e "   ${BLUE}👩‍🏫 Teacher:${NC} teacher@zenflow.example.com / yoga123"
    echo ""
    echo -e "${BLUE}🎯 Things to Try:${NC}"
    echo "   1. Explore the calming, wellness-focused design"
    echo "   2. Browse different yoga class types and pricing"
    echo "   3. Try the prenatal and advanced class packages"
    echo "   4. Check out workshop and private session options"
    echo "   5. Notice the elegant typography and peaceful colors"
elif [ "$COMMUNITY_TYPE" = "hackerspace" ]; then
    echo "🎉 Oslo Hackerspace Demo is now running!"
    echo "========================================"
    echo ""
    echo -e "${CYAN}💻 Oslo Hackerspace${NC}"
    echo -e "   ${GREEN}🌐 URL:${NC} http://localhost:8080"
    echo -e "   ${BLUE}👤 Member:${NC} maker@hackerspace.example.com / build123"
    echo -e "   ${BLUE}🧑‍🎓 Mentor:${NC} mentor@hackerspace.example.com / hack123"
    echo ""
    echo -e "${BLUE}🎯 Things to Try:${NC}"
    echo "   1. Experience the terminal-inspired dark design"
    echo "   2. Explore tech-focused pricing (3D printing, laser cutting)"
    echo "   3. Check out workshops and mentorship options"
    echo "   4. Try the 24/7 access membership model"
    echo "   5. Notice the monospace fonts and green-on-black theme"
fi

echo ""
echo -e "${YELLOW}🔧 Admin Access:${NC}"
echo -e "   ${BLUE}👑 Admin:${NC} admin@example.com / admin123"
echo ""
echo -e "${BLUE}🔧 Configuration Hot-Reload Test:${NC}"
echo -e "   • Edit config/${COMMUNITY_TYPE}.yaml"
echo -e "   • Changes will be applied immediately without restarting"
echo -e "   • Try changing colors, pricing, or content and refresh the page"
echo ""
echo -e "${BLUE}📂 Configuration File:${NC}"
echo -e "   • ${GREEN}Active Config:${NC} config/${COMMUNITY_TYPE}.yaml"
echo ""
echo -e "${BLUE}🔄 To try the other community:${NC}"
if [ "$COMMUNITY_TYPE" = "yoga-studio" ]; then
    echo -e "   • Stop this demo (Ctrl+C) and run: ${GREEN}./run-demo.sh hack${NC}"
else
    echo -e "   • Stop this demo (Ctrl+C) and run: ${GREEN}./run-demo.sh yoga${NC}"
fi
echo ""
echo -e "${BLUE}🏗️ Architecture Highlights:${NC}"
echo "   • Multi-tenant design with shared infrastructure"
echo "   • Community-specific branding and features"
echo "   • Hot-reload configuration system"
echo "   • Dynamic CSS generation"
echo "   • YAML-driven customization"
echo "   • Flexible pricing and feature models"
echo ""
echo -e "${YELLOW}📚 For more details, see: demos/multi-community/README.md${NC}"
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