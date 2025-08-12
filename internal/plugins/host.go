package plugins

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	pb "samskipnad/pkg/proto/v1"
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
	Client   *plugin.Client
	Process  *exec.Cmd
	Services PluginServices
}

// PluginServices contains the gRPC clients for core services exposed to plugins
type PluginServices struct {
	UserProfile pb.UserProfileServiceClient
	// TODO: Add other services as we generate their protobuf definitions
	// CommunityManagement pb.CommunityManagementServiceClient
	// ItemManagement     pb.ItemManagementServiceClient
	// Payment            pb.PaymentServiceClient
	// EventBus           pb.EventBusServiceClient
}

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

	// Create plugin client configuration
	pluginConfig := &plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         PluginMap,
		Cmd:             exec.Command(path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Output: log.Writer(),
			Level:  hclog.Info,
		}),
	}

	// Start the plugin
	client := plugin.NewClient(pluginConfig)
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to create RPC client: %v", err)
	}

	// Get the plugin interface
	raw, err := rpcClient.Dispense("samskipnad_plugin")
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense plugin: %v", err)
	}

	// Type assert to our plugin interface
	pluginInterface, ok := raw.(SamskipnadPlugin)
	if !ok {
		client.Kill()
		return fmt.Errorf("plugin does not implement SamskipnadPlugin interface")
	}

	// Initialize the plugin with core services
	services := h.createPluginServices(ctx)
	if err := pluginInterface.Initialize(ctx, services); err != nil {
		client.Kill()
		return fmt.Errorf("failed to initialize plugin: %v", err)
	}

	// Store the plugin
	h.plugins[name] = &Plugin{
		Name:     name,
		Path:     path,
		Client:   client,
		Services: services,
	}

	log.Printf("Plugin %s loaded successfully from %s", name, path)
	return nil
}

// UnloadPlugin unloads a plugin
func (h *PluginHost) UnloadPlugin(name string) error {
	plugin, exists := h.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	// Kill the plugin process
	plugin.Client.Kill()

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

// Plugin interface and handshake configuration
var (
	// HandshakeConfig is the configuration for the plugin handshake
	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "SAMSKIPNAD_PLUGIN",
		MagicCookieValue: "samskipnad_v1",
	}

	// PluginMap is the map of plugins we can dispense
	PluginMap = map[string]plugin.Plugin{
		"samskipnad_plugin": &SamskipnadPluginGRPC{},
	}
)

// SamskipnadPlugin is the interface that all plugins must implement
type SamskipnadPlugin interface {
	// Initialize is called when the plugin is loaded
	Initialize(ctx context.Context, services PluginServices) error

	// Name returns the plugin's name
	Name() string

	// Version returns the plugin's version
	Version() string

	// Execute performs the plugin's main functionality
	Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error)

	// Shutdown is called when the plugin is being unloaded
	Shutdown(ctx context.Context) error
}

// SamskipnadPluginGRPC is the gRPC implementation of the plugin interface
type SamskipnadPluginGRPC struct {
	plugin.Plugin
	Impl SamskipnadPlugin
}

func (p *SamskipnadPluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	// TODO: Register the plugin service with the gRPC server
	// This would require defining a plugin-specific protobuf service
	return nil
}

func (p *SamskipnadPluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	// TODO: Return a gRPC client for the plugin service
	// This would require implementing the client-side of the plugin protocol
	return nil, nil
}