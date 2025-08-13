package plugins

import (
	"context"
	"fmt"
	"log"
	"os"

	"samskipnad/pkg/sdk"
)

// PluginHost manages the lifecycle of plugins
type PluginHost struct {
	plugins map[string]*Plugin
	config  *PluginConfig
}

// Plugin represents a loaded plugin instance
type Plugin struct {
	Name     string
	Path     string
	Services PluginServices
}

// PluginServices contains the gRPC clients for core services exposed to plugins
// Use shared PluginServices from SDK
type PluginServices = sdk.PluginServices

// PluginConfig contains configuration for the plugin system
type PluginConfig struct {
	PluginDir     string
	MaxPlugins    int
	PluginTimeout int // seconds
}

// NewPluginHost creates a new plugin host
func NewPluginHost(config *PluginConfig) *PluginHost {
	if config == nil {
		config = &PluginConfig{
			PluginDir:     "./plugins",
			MaxPlugins:    10,
			PluginTimeout: 30,
		}
	}

	return &PluginHost{
		plugins: make(map[string]*Plugin),
		config:  config,
	}
}

// LoadPlugin loads a plugin from the specified path
func (h *PluginHost) LoadPlugin(ctx context.Context, name, path string) error {
	if _, exists := h.plugins[name]; exists {
		return fmt.Errorf("plugin %s is already loaded", name)
	}

	if len(h.plugins) >= h.config.MaxPlugins {
		return fmt.Errorf("maximum number of plugins (%d) reached", h.config.MaxPlugins)
	}
	// Verify the plugin binary exists
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("plugin not found at path %s: %v", path, err)
	}

	// For this demo, we won't start the process yet; we just prepare services
	services := h.createPluginServices(ctx)

	// Store the plugin
	h.plugins[name] = &Plugin{
		Name:     name,
		Path:     path,
		Services: services,
	}

	log.Printf("Plugin %s registered from %s", name, path)
	return nil
}

// UnloadPlugin unloads a plugin
func (h *PluginHost) UnloadPlugin(name string) error {
	_, exists := h.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	// Remove from our map
	delete(h.plugins, name)

	log.Printf("Plugin %s unloaded successfully", name)
	return nil
}

// GetLoadedPlugins returns a list of currently loaded plugin names
func (h *PluginHost) GetLoadedPlugins() []string {
	names := make([]string, 0, len(h.plugins))
	for name := range h.plugins {
		names = append(names, name)
	}
	return names
}

// Shutdown gracefully shuts down the plugin host
func (h *PluginHost) Shutdown() {
	log.Println("Shutting down plugin host...")
	for name := range h.plugins {
		if err := h.UnloadPlugin(name); err != nil {
			log.Printf("Error unloading plugin %s: %v", name, err)
		}
	}
}

// createPluginServices creates gRPC clients for core services that plugins can use
func (h *PluginHost) createPluginServices(ctx context.Context) PluginServices {
	// TODO: In a real implementation, these would connect to actual gRPC servers
	// For now, we'll return empty clients as this is the foundation
	return PluginServices{
		// UserProfile: pb.NewUserProfileServiceClient(conn),
		// CommunityManagement: pb.NewCommunityManagementServiceClient(conn),
	}
}

// Use the SDK's SamskipnadPlugin interface and GRPC wrapper; no local types needed here
