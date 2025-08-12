# Phase 2 Demo: Plugin System Foundation

This demo showcases the revolutionary plugin system architecture that enables safe, isolated extensibility of the Samskipnad platform.

## What This Demo Shows

- **Plugin Host Architecture**: Process isolation with gRPC communication
- **Plugin SDK**: Developer-friendly interface for creating plugins
- **Real Plugin Examples**: Working RSS importer and sample plugins
- **Service Integration**: How plugins access core platform services
- **Plugin Lifecycle**: Loading, execution, and cleanup processes
- **Security & Isolation**: How plugin failures don't affect the core system

## Quick Start

1. **Run the plugin system demo:**
   ```bash
   ./run-demo.sh
   ```

2. **Watch the output** to see plugins being loaded and executed

3. **Explore the plugin code** to understand development patterns

## Demo Components

### 1. Plugin Host Application
A command-line application that demonstrates:
- Loading plugins from the filesystem
- Establishing gRPC communication channels
- Executing plugin operations safely
- Managing plugin lifecycles
- Handling plugin failures gracefully

### 2. Example Plugins

#### RSS Importer Plugin
- **Purpose**: Demonstrates content import capabilities
- **Features**: URL parsing, RSS feed simulation, service integration
- **Location**: `plugins/rss-importer/`

#### Simple Calculator Plugin  
- **Purpose**: Shows basic computation plugins
- **Features**: Mathematical operations, parameter validation
- **Location**: `plugins/calculator/`

#### Echo Plugin
- **Purpose**: Minimal example for learning
- **Features**: Parameter echo, basic plugin structure
- **Location**: `plugins/echo/`

### 3. Plugin Development SDK
The demo includes the complete SDK that makes plugin development simple:
- **BasePlugin**: Reusable base implementation
- **Service Clients**: Access to core platform services
- **gRPC Integration**: Automatic protocol handling
- **Error Management**: Robust error handling patterns

## Architecture Highlights

### Process Isolation
Each plugin runs as a separate OS process, providing:
- **Crash Protection**: Plugin failures don't affect the core system
- **Resource Isolation**: Plugins can't consume unlimited resources
- **Security Boundaries**: Plugins operate in sandboxed environments
- **Version Independence**: Plugins can use different library versions

### gRPC Communication
Type-safe communication between host and plugins:
- **Protocol Buffers**: Well-defined service interfaces
- **Streaming Support**: For large data transfers
- **Error Propagation**: Proper error handling across processes
- **Performance**: Efficient binary protocol

### Service-Oriented Architecture
Plugins access core services through well-defined APIs:
- **UserProfileService**: User and tenant management
- **Future Services**: Content, payment, notification services
- **Versioned APIs**: Backward compatibility guarantee
- **Authentication**: Secure service access

## Demo Scenarios

### Scenario 1: RSS Content Import
1. Plugin host loads the RSS importer plugin
2. Plugin receives RSS URL and tenant information
3. Plugin fetches and parses RSS feed (simulated)
4. Plugin creates content items via ItemManagementService
5. Plugin returns import summary

### Scenario 2: Plugin Failure Handling
1. Host loads a plugin that deliberately fails
2. Plugin crashes during execution
3. Host detects failure and continues operating
4. Core system remains completely stable
5. Other plugins continue working normally

### Scenario 3: Service Integration
1. Plugin requests user profile information
2. Host provides gRPC client for UserProfileService  
3. Plugin makes authenticated calls to core services
4. Plugin processes data and returns results
5. All communication is type-safe and logged

## Plugin Development Example

Creating a new plugin is straightforward:

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
    sdk.Serve(NewMyPlugin())
}
```

## Technical Stack

- **Core**: Go with HashiCorp go-plugin library
- **Communication**: gRPC with Protocol Buffers
- **Isolation**: OS processes with controlled interfaces
- **Service Discovery**: Plugin manifest and registration
- **Build System**: Integrated with project Makefile

## Files in This Demo

- `run-demo.sh` - Main demo runner script
- `host/` - Plugin host application demonstrating plugin management
- `plugins/` - Example plugins showing different patterns
- `README.md` - This documentation

## Plugin Development Workflow

1. **Create Plugin Structure**:
   ```bash
   mkdir my-plugin
   cd my-plugin
   go mod init my-plugin
   ```

2. **Implement Plugin Interface**:
   - Embed `sdk.BasePlugin`
   - Implement `Execute()` method
   - Add any custom functionality

3. **Build Plugin**:
   ```bash
   go build -o my-plugin
   ```

4. **Test with Host**:
   ```bash
   ./plugin-host --plugin ./my-plugin
   ```

## Security Features

### Process Boundaries
- Plugins cannot access host memory directly
- Filesystem access is controlled by OS permissions  
- Network access can be restricted via containers

### API Boundaries  
- All communication through defined gRPC interfaces
- Input validation at service boundaries
- Authentication required for service access

### Resource Management
- Plugin timeout controls prevent hanging
- Memory and CPU limits can be enforced
- Graceful shutdown procedures

## What's Next: Phase 3

This plugin foundation enables:

1. **Creator Studio**: Visual plugin development environment
2. **Plugin Marketplace**: Community-driven plugin ecosystem  
3. **Advanced Services**: More core services exposed to plugins
4. **Auto-discovery**: Dynamic plugin loading and configuration
5. **Plugin Composition**: Plugins that use other plugins

## Performance Considerations

The demo shows how the plugin system maintains performance:
- **Lazy Loading**: Plugins loaded only when needed
- **Connection Pooling**: Efficient gRPC connection management
- **Parallel Execution**: Multiple plugins can run simultaneously
- **Resource Cleanup**: Proper cleanup prevents resource leaks

## Troubleshooting

Common issues and solutions:
- **Plugin Won't Load**: Check file permissions and Go build errors
- **gRPC Errors**: Verify protocol buffer generation
- **Service Access**: Ensure proper authentication setup
- **Process Cleanup**: Use the provided cleanup scripts

The plugin system is designed to be robust and developer-friendly while maintaining security and performance.