package sdk

import (
	"context"
	"log"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	pb "samskipnad/pkg/proto/v1"
)

// BasePlugin provides a basic implementation that other plugins can embed
type BasePlugin struct {
	name     string
	version  string
	services PluginServices
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

// NewBasePlugin creates a new base plugin
func NewBasePlugin(name, version string) *BasePlugin {
	return &BasePlugin{
		name:    name,
		version: version,
	}
}

// Initialize implements the SamskipnadPlugin interface
func (p *BasePlugin) Initialize(ctx context.Context, services PluginServices) error {
	p.services = services
	log.Printf("Plugin %s v%s initialized", p.name, p.version)
	return nil
}

// Name returns the plugin's name
func (p *BasePlugin) Name() string {
	return p.name
}

// Version returns the plugin's version
func (p *BasePlugin) Version() string {
	return p.version
}

// Execute is a default implementation that should be overridden by actual plugins
func (p *BasePlugin) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"status":  "success",
		"message": "Base plugin executed successfully",
		"plugin":  p.name,
		"version": p.version,
	}, nil
}

// Shutdown implements the SamskipnadPlugin interface
func (p *BasePlugin) Shutdown(ctx context.Context) error {
	log.Printf("Plugin %s v%s shutting down", p.name, p.version)
	return nil
}

// GetServices returns the core services available to the plugin
func (p *BasePlugin) GetServices() PluginServices {
	return p.services
}

// Plugin main function that plugins should call to start serving
func Serve(pluginImpl interface{}) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"samskipnad_plugin": &SamskipnadPluginGRPC{Impl: pluginImpl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

// Configuration from host (shared with internal/plugins)
var (
	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "SAMSKIPNAD_PLUGIN",
		MagicCookieValue: "samskipnad_v1",
	}
)

// SamskipnadPlugin is the interface that all plugins must implement
type SamskipnadPlugin interface {
	Initialize(ctx context.Context, services PluginServices) error
	Name() string
	Version() string
	Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error)
	Shutdown(ctx context.Context) error
}

// SamskipnadPluginGRPC is the gRPC implementation wrapper
type SamskipnadPluginGRPC struct {
	plugin.Plugin
	Impl interface{}
}

func (p *SamskipnadPluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	// TODO: Register the actual plugin service
	return nil
}

func (p *SamskipnadPluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	// TODO: Return the actual plugin client
	return p.Impl, nil
}