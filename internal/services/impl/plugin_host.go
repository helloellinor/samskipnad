package impl

import (
	"context"
	"fmt"
	"log"

	"samskipnad/internal/plugins"
	"samskipnad/internal/services"
)

// PluginHostServiceImpl implements the PluginHostService interface
type PluginHostServiceImpl struct {
	host *plugins.PluginHost
}

// NewPluginHostService creates a new plugin host service
func NewPluginHostService() services.PluginHostService {
	config := &plugins.PluginConfig{
		PluginDir:     "./plugins",
		MaxPlugins:    10,
		PluginTimeout: 30,
	}
	
	return &PluginHostServiceImpl{
		host: plugins.NewPluginHost(config),
	}
}

// Initialize initializes the plugin host
func (p *PluginHostServiceImpl) Initialize(ctx context.Context) error {
	log.Println("Initializing plugin host service")
	// TODO: Auto-discover and load plugins from plugin directory
	return nil
}

// LoadPlugin loads a plugin from the specified path
func (p *PluginHostServiceImpl) LoadPlugin(ctx context.Context, name, path string) error {
	return p.host.LoadPlugin(ctx, name, path)
}

// UnloadPlugin unloads a plugin
func (p *PluginHostServiceImpl) UnloadPlugin(ctx context.Context, name string) error {
	return p.host.UnloadPlugin(name)
}

// GetLoadedPlugins returns a list of currently loaded plugin names
func (p *PluginHostServiceImpl) GetLoadedPlugins(ctx context.Context) []string {
	return p.host.GetLoadedPlugins()
}

// ExecutePlugin executes a plugin with the given parameters
func (p *PluginHostServiceImpl) ExecutePlugin(ctx context.Context, name string, params map[string]interface{}) (map[string]interface{}, error) {
	loadedPlugins := p.host.GetLoadedPlugins()
	for _, pluginName := range loadedPlugins {
		if pluginName == name {
			// TODO: Implement actual plugin execution
			// For now, return a placeholder response
			return map[string]interface{}{
				"status":  "success",
				"message": fmt.Sprintf("Plugin %s executed successfully", name),
				"params":  params,
			}, nil
		}
	}
	
	return nil, fmt.Errorf("plugin %s is not loaded", name)
}

// Shutdown gracefully shuts down the plugin host
func (p *PluginHostServiceImpl) Shutdown(ctx context.Context) error {
	log.Println("Shutting down plugin host service")
	p.host.Shutdown()
	return nil
}