#!/bin/bash

# Phase 2 Demo Runner for Samskipnad Plugin System
# This script demonstrates the plugin system architecture and capabilities

set -e

echo "🔌 Starting Phase 2 Demo: Plugin System Foundation"
echo "=================================================="

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DEMO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "📁 Project root: $PROJECT_ROOT"
echo "📁 Demo directory: $DEMO_DIR"

# Change to project root
cd "$PROJECT_ROOT"

# Check prerequisites
echo -e "${BLUE}🔍 Checking prerequisites...${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21 or later.${NC}"
    exit 1
fi

if ! command -v protoc &> /dev/null; then
    echo -e "${YELLOW}⚠️  protoc not found. Installing protocol buffers...${NC}"
    # The Makefile will handle protoc dependencies
fi

echo -e "${GREEN}✅ Go is installed: $(go version)${NC}"

# Generate protocol buffers
echo -e "${BLUE}🔧 Generating gRPC code...${NC}"
make proto

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Protocol buffer generation failed. Please check error messages.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Protocol buffers generated successfully${NC}"

# Build the core system
echo -e "${BLUE}🔨 Building the core system...${NC}"
make build

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Core system build failed. Please check the error messages above.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Core system built successfully${NC}"

# Build example plugins
echo -e "${BLUE}🧩 Building example plugins...${NC}"

# Build RSS importer plugin
cd "$PROJECT_ROOT/examples/plugins/rss-importer"
go build -o rss-importer .

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ RSS importer plugin build failed.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ RSS importer plugin built successfully${NC}"

# Create additional demo plugins
cd "$DEMO_DIR"

# Create simple calculator plugin
mkdir -p plugins/calculator
cat > plugins/calculator/main.go << 'EOF'
package main

import (
	"context"
	"fmt"
	"samskipnad/pkg/sdk"
)

type CalculatorPlugin struct {
	*sdk.BasePlugin
}

func NewCalculatorPlugin() *CalculatorPlugin {
	return &CalculatorPlugin{
		BasePlugin: sdk.NewBasePlugin("calculator", "1.0.0"),
	}
}

func (p *CalculatorPlugin) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter is required")
	}

	a, aOk := params["a"].(float64)
	b, bOk := params["b"].(float64)
	if !aOk || !bOk {
		return nil, fmt.Errorf("parameters 'a' and 'b' must be numbers")
	}

	var result float64
	var err error

	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}

	return map[string]interface{}{
		"status":    "success",
		"plugin":    p.Name(),
		"version":   p.Version(),
		"operation": operation,
		"operands":  []float64{a, b},
		"result":    result,
		"message":   fmt.Sprintf("%.2f %s %.2f = %.2f", a, operation, b, result),
	}, err
}

func main() {
	plugin := NewCalculatorPlugin()
	fmt.Printf("Starting Calculator Plugin %s v%s\n", plugin.Name(), plugin.Version())
	sdk.Serve(plugin)
}
EOF

# Initialize calculator plugin module
cd plugins/calculator
go mod init calculator-plugin
go mod tidy
go build -o calculator .

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Calculator plugin build failed.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Calculator plugin built successfully${NC}"

# Create echo plugin
cd "$DEMO_DIR"
mkdir -p plugins/echo
cat > plugins/echo/main.go << 'EOF'
package main

import (
	"context"
	"fmt"
	"samskipnad/pkg/sdk"
)

type EchoPlugin struct {
	*sdk.BasePlugin
}

func NewEchoPlugin() *EchoPlugin {
	return &EchoPlugin{
		BasePlugin: sdk.NewBasePlugin("echo", "1.0.0"),
	}
}

func (p *EchoPlugin) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	message, ok := params["message"].(string)
	if !ok {
		message = "Hello from Echo Plugin!"
	}

	return map[string]interface{}{
		"status":       "success",
		"plugin":       p.Name(),
		"version":      p.Version(),
		"echo":         message,
		"params_count": len(params),
		"params":       params,
		"timestamp":    fmt.Sprintf("%d", ctx.Value("timestamp")),
	}, nil
}

func main() {
	plugin := NewEchoPlugin()
	fmt.Printf("Starting Echo Plugin %s v%s\n", plugin.Name(), plugin.Version())
	sdk.Serve(plugin)
}
EOF

cd plugins/echo
go mod init echo-plugin
go mod tidy
go build -o echo .

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Echo plugin build failed.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Echo plugin built successfully${NC}"

# Create plugin host demo application
cd "$DEMO_DIR"
mkdir -p host
cat > host/main.go << 'EOF'
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"samskipnad/internal/plugins"
)

func main() {
	fmt.Println("🔌 Samskipnad Plugin System Demo")
	fmt.Println("=================================")

	// Create plugin host configuration
	config := &plugins.PluginConfig{
		PluginDir:     "./plugins",
		MaxPlugins:    10,
		PluginTimeout: 30,
	}

	// Create plugin host
	host := plugins.NewPluginHost(config)
	
	// Demo plugins to test
	pluginTests := []struct {
		name   string
		path   string
		params map[string]interface{}
	}{
		{
			name: "Echo Plugin",
			path: "./plugins/echo/echo",
			params: map[string]interface{}{
				"message": "Hello from the plugin system!",
				"demo":    true,
			},
		},
		{
			name: "Calculator Plugin - Addition",
			path: "./plugins/calculator/calculator", 
			params: map[string]interface{}{
				"operation": "add",
				"a":         10.5,
				"b":         25.3,
			},
		},
		{
			name: "Calculator Plugin - Division",
			path: "./plugins/calculator/calculator",
			params: map[string]interface{}{
				"operation": "divide",
				"a":         100.0,
				"b":         4.0,
			},
		},
		{
			name: "RSS Importer Plugin",
			path: "../../examples/plugins/rss-importer/rss-importer",
			params: map[string]interface{}{
				"rss_url":   "https://feeds.example.com/yoga-blog.xml",
				"tenant_id": "demo-tenant",
			},
		},
	}

	ctx := context.Background()

	fmt.Println("\n🎯 Running Plugin Demos:")
	fmt.Println("========================")

	for i, test := range pluginTests {
		fmt.Printf("\n%d. %s\n", i+1, test.name)
		fmt.Println(strings.Repeat("-", len(test.name)+3))

		// Check if plugin file exists
		absPath, err := filepath.Abs(test.path)
		if err != nil {
			fmt.Printf("❌ Failed to resolve plugin path: %v\n", err)
			continue
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			fmt.Printf("❌ Plugin not found: %s\n", absPath)
			continue
		}

		// Load plugin
		fmt.Printf("📂 Loading plugin from: %s\n", absPath)
		plugin, err := host.LoadPlugin(ctx, test.name, absPath)
		if err != nil {
			fmt.Printf("❌ Failed to load plugin: %v\n", err)
			continue
		}

		fmt.Printf("✅ Plugin loaded: %s\n", plugin.Name)

		// Execute plugin
		fmt.Printf("⚡ Executing with params: %+v\n", test.params)
		
		startTime := time.Now()
		result, err := host.ExecutePlugin(ctx, plugin.Name, test.params)
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ Plugin execution failed: %v\n", err)
		} else {
			fmt.Printf("✅ Plugin executed successfully in %v\n", duration)
			fmt.Printf("📋 Result: %+v\n", result)
		}

		// Unload plugin
		err = host.UnloadPlugin(ctx, plugin.Name)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to unload plugin: %v\n", err)
		} else {
			fmt.Printf("🗑️  Plugin unloaded successfully\n")
		}

		// Add delay between tests for readability
		if i < len(pluginTests)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("\n🎉 Plugin System Demo Complete!")
	fmt.Println("===============================")
	fmt.Println("\n📚 Key Features Demonstrated:")
	fmt.Println("  • Process isolation - each plugin runs in its own process")
	fmt.Println("  • gRPC communication - type-safe plugin communication")
	fmt.Println("  • Plugin lifecycle - loading, execution, and cleanup")
	fmt.Println("  • Error handling - graceful failure recovery")
	fmt.Println("  • Multiple plugin types - utility, content, and processing plugins")
	fmt.Println("\n🔧 Next Steps:")
	fmt.Println("  • Explore the plugin source code in the plugins/ directory")
	fmt.Println("  • Try creating your own plugin using the SDK")
	fmt.Println("  • Check out the Phase 2 documentation for more details")
}
EOF

cd host
go mod init plugin-host-demo
go mod edit -replace samskipnad="$PROJECT_ROOT"
go mod tidy
go build -o plugin-host .

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Plugin host demo build failed.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Plugin host demo built successfully${NC}"

# Run the demo
echo -e "${PURPLE}🚀 Running Plugin System Demo...${NC}"
echo ""

cd "$DEMO_DIR"
./host/plugin-host

echo ""
echo -e "${GREEN}🎉 Phase 2 Demo completed successfully!${NC}"
echo ""
echo -e "${BLUE}📚 What you just saw:${NC}"
echo "  • Plugins loaded as separate processes"
echo "  • gRPC communication between host and plugins"
echo "  • Type-safe parameter passing and result handling"
echo "  • Graceful error handling and recovery"
echo "  • Complete plugin lifecycle management"
echo ""
echo -e "${YELLOW}📝 Explore the code:${NC}"
echo "  • Plugin SDK: $PROJECT_ROOT/pkg/sdk/"
echo "  • Plugin Host: $PROJECT_ROOT/internal/plugins/"
echo "  • Example Plugins: $DEMO_DIR/plugins/"
echo "  • gRPC Definitions: $PROJECT_ROOT/pkg/proto/"
echo ""
echo -e "${PURPLE}🔧 Try creating your own plugin:${NC}"
echo "  1. Copy one of the example plugins"
echo "  2. Modify the Execute() method"
echo "  3. Build and test with the plugin host"
echo ""
echo -e "${GREEN}For more details, see: demos/phase2/README.md${NC}"
EOF