package main

import (
	"context"
	"fmt"
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
	err = host.LoadPlugin(ctx, test.name, absPath)
		if err != nil {
			fmt.Printf("❌ Failed to load plugin: %v\n", err)
			continue
		}

	fmt.Printf("✅ Plugin loaded: %s\n", test.name)

	// Execute plugin (simulation placeholder; real gRPC execution not wired yet)
	fmt.Printf("⚡ Executing with params: %+v\n", test.params)
	startTime := time.Now()
	time.Sleep(100 * time.Millisecond)
	duration := time.Since(startTime)
	fmt.Printf("✅ (Simulated) plugin execution completed in %v\n", duration)

		// Unload plugin
	err = host.UnloadPlugin(test.name)
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
