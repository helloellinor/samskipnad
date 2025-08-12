# Samskipnad Plugin Development Guide

## Phase 2: Plugin System Overview

The Samskipnad platform now supports a plugin architecture that allows developers to extend the platform's functionality while maintaining system stability and security. This guide will help you get started building plugins.

## Architecture

### Core Principles
- **Process Isolation**: Plugins run as separate processes, ensuring system stability
- **gRPC Communication**: Plugins communicate with the core system via gRPC APIs
- **Standard Interface**: All plugins implement the `SamskipnadPlugin` interface
- **Service Access**: Plugins can access core services through well-defined APIs

### Plugin Lifecycle
1. **Discovery**: Host discovers plugin binaries
2. **Loading**: Host starts plugin process
3. **Handshake**: Plugin and host establish communication
4. **Initialization**: Plugin receives core service clients
5. **Execution**: Plugin performs its functionality
6. **Shutdown**: Plugin gracefully terminates

## Getting Started

### 1. Plugin Structure

Every plugin must implement the `SamskipnadPlugin` interface:

```go
type SamskipnadPlugin interface {
    Initialize(ctx context.Context, services PluginServices) error
    Name() string
    Version() string
    Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error)
    Shutdown(ctx context.Context) error
}
```

### 2. Basic Plugin Template

```go
package main

import (
    "context"
    "samskipnad/pkg/sdk"
)

type MyPlugin struct {
    *sdk.BasePlugin
}

func NewMyPlugin() *MyPlugin {
    return &MyPlugin{
        BasePlugin: sdk.NewBasePlugin("my-plugin", "1.0.0"),
    }
}

func (p *MyPlugin) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    // Your plugin logic here
    return map[string]interface{}{
        "status": "success",
        "message": "Plugin executed successfully",
    }, nil
}

func main() {
    plugin := NewMyPlugin()
    sdk.Serve(plugin)
}
```

### 3. Building Your Plugin

```bash
# Build your plugin as an executable
go build -o my-plugin main.go

# Make it executable
chmod +x my-plugin
```

## Available Services

Currently available core services for plugins:

### UserProfileService
Access user authentication, profiles, and permissions:

```go
services := p.GetServices()
if services.UserProfile != nil {
    user, err := services.UserProfile.GetProfile(ctx, &pb.GetProfileRequest{
        UserId: userID,
    })
    if err != nil {
        return nil, err
    }
    // Use user data...
}
```

## Example: RSS Feed Importer Plugin

The RSS importer plugin demonstrates how to build a content processing plugin:

```go
func (p *RSSImporterPlugin) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    rssURL := params["rss_url"].(string)
    
    // 1. Fetch RSS feed
    // 2. Parse content
    // 3. Create items via ItemManagementService
    // 4. Return results
    
    return map[string]interface{}{
        "status": "success",
        "items_processed": 5,
    }, nil
}
```

## Plugin Configuration

Plugins can accept configuration through the `Execute` parameters:

```go
// Host calls plugin with configuration
result, err := pluginHost.ExecutePlugin(ctx, "my-plugin", map[string]interface{}{
    "config_option": "value",
    "api_key": "secret",
})
```

## Development Tools

### Makefile Targets

```bash
# Generate gRPC code (when adding new services)
make proto

# Build the main application
make build

# Run tests
make test

# Development mode with hot reload
make dev
```

### Plugin Development Workflow

1. **Create Plugin Directory**: `mkdir -p examples/plugins/my-plugin`
2. **Implement Plugin**: Create `main.go` with your plugin logic
3. **Build Plugin**: `go build -o my-plugin main.go`
4. **Test Integration**: Load plugin into host and test functionality

## Best Practices

### Error Handling
- Always handle errors gracefully
- Return meaningful error messages
- Don't panic - return errors instead

### Performance
- Keep plugin execution fast
- Use context for cancellation
- Avoid blocking operations in critical paths

### Security
- Validate all input parameters
- Don't store sensitive data in plugin state
- Use provided service APIs rather than direct database access

## Future Roadmap

### Coming Soon (Phase 2 Completion)
- [ ] Complete gRPC APIs for all core services
- [ ] Plugin auto-discovery from directory
- [ ] Plugin configuration files
- [ ] Hot-reload plugin capability

### Phase 3: Creator Studio
- [ ] Plugin marketplace UI
- [ ] One-click plugin installation
- [ ] Plugin validation pipeline
- [ ] Visual plugin configuration

## Troubleshooting

### Common Issues

**Plugin won't load**: Check that the plugin binary is executable and implements the correct interface.

**gRPC errors**: Ensure the plugin uses the same protobuf definitions as the host.

**Service unavailable**: Verify that the required core services are initialized and available.

### Debugging

Enable debug logging in the plugin host:
```go
// In host configuration
Logger: hclog.New(&hclog.LoggerOptions{
    Level: hclog.Debug,
})
```

## Support

For plugin development support:
- Check the examples in `examples/plugins/`
- Review the SDK documentation in `pkg/sdk/`
- See the core service interfaces in `internal/services/interfaces.go`

---

**Note**: This plugin system is currently in Phase 2 development. APIs may change as we complete the implementation. Pin to specific versions for production use.