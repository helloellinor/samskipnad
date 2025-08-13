package main

import (
	"context"
	"fmt"
	"log"

	"samskipnad/pkg/sdk"
)

// RSSImporterPlugin is a simple plugin that imports RSS feeds
type RSSImporterPlugin struct {
	*sdk.BasePlugin
}

// NewRSSImporterPlugin creates a new RSS importer plugin
func NewRSSImporterPlugin() *RSSImporterPlugin {
	return &RSSImporterPlugin{
		BasePlugin: sdk.NewBasePlugin("rss-importer", "1.0.0"),
	}
}

// Execute implements the plugin's main functionality
func (p *RSSImporterPlugin) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	log.Printf("RSS Importer Plugin executing with params: %+v", params)

	// Extract RSS URL from parameters
	rssURL, ok := params["rss_url"].(string)
	if !ok {
		return nil, fmt.Errorf("rss_url parameter is required")
	}

	// Simulate RSS feed processing
	result := map[string]interface{}{
		"status":      "success",
		"plugin":      p.Name(),
		"version":     p.Version(),
		"rss_url":     rssURL,
		"items_found": 5, // Simulated
		"message":     fmt.Sprintf("Successfully processed RSS feed from %s", rssURL),
	}

	// In a real implementation, this would:
	// 1. Fetch the RSS feed from the URL
	// 2. Parse the RSS XML
	// 3. Use the ItemManagementService to create content items
	// 4. Handle any errors gracefully

	// Example of how we would use the core services:
	// services := p.GetServices()
	// if services.ItemManagement != nil {
	//     // Create items using the ItemManagementService
	//     for _, item := range rssItems {
	//         _, err := services.ItemManagement.CreateItem(ctx, &pb.CreateItemRequest{
	//             TenantId: tenantID,
	//             ItemType: "rss_article",
	//             Data: convertRSSItemToStruct(item),
	//         })
	//         if err != nil {
	//             return nil, fmt.Errorf("failed to create item: %v", err)
	//         }
	//     }
	// }

	return result, nil
}

func main() {
	// Create and serve the plugin
	plugin := NewRSSImporterPlugin()
	// Avoid writing to stdout before handshake; use stderr (log)
	log.Printf("Starting RSS Importer Plugin %s v%s", plugin.Name(), plugin.Version())

	// Start serving the plugin
	sdk.Serve(plugin)
}
