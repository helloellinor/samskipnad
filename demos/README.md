# Samskipnad Platform Demos

This directory contains demonstrations of the Samskipnad platform functionality organized by implementation phases.

## Demo Structure

### Phase 1: Core Platform
- **Location**: `phase1/`
- **Description**: Demonstrates the core Samskipnad platform with web interface, authentication, community management, and payment processing
- **Technologies**: Go, SQLite, HTML/CSS, HTMX

### Phase 2: Plugin System Foundation
- **Location**: `phase2/`
- **Description**: Demonstrates the plugin system architecture with gRPC APIs, process isolation, and example plugins
- **Technologies**: Go, gRPC, Protocol Buffers, HashiCorp go-plugin

### Multi-Community Demo: Yoga Studio & Hackerspace
- **Location**: `multi-community/`
- **Description**: Showcases multi-tenant capabilities with two distinct community types: yoga studio and hackerspace
- **Technologies**: YAML-driven configuration, hot-reload, dynamic theming, multi-tenant architecture

## Quick Start

Each phase demo includes its own README with specific instructions. To run a demo:

```bash
# Navigate to the desired demo
cd demos/phase1         # Core platform demo
cd demos/phase2         # Plugin system demo  
cd demos/multi-community # Multi-tenant community demo

# Follow the instructions in that demo's README.md
cat README.md
```

## Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (for Phase 2)
- Make (for build automation)

For detailed setup instructions, see the main project README.

## Demo Philosophy

These demos are designed to:
1. **Showcase real functionality** - Each demo runs actual working code
2. **Demonstrate architectural principles** - Show how the system is designed and why
3. **Enable developer understanding** - Provide clear examples for extending the platform
4. **Validate implementation** - Prove that the architecture works as intended

## Phase Progression

The demos show the evolution of the Samskipnad platform:

1. **Phase 1** establishes the core platform with essential features
2. **Phase 2** introduces the plugin system foundation for extensibility  
3. **Multi-Community Demo** demonstrates real-world multi-tenant usage with distinct community types
4. **Phase 3** (future) will demonstrate the Creator Studio and advanced plugin ecosystem

Each phase builds upon the previous one while maintaining stability and backward compatibility.