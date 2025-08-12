# Samskipnad Platform
## From Community Management Software to Extensible Creator Platform

**Samskipnad** is evolving from a simple community management application into a comprehensive, extensible platform that empowers creators to build and customize their own community experiences. Our mission is to transform from a closed application into an open ecosystem that fosters third-party development and community-driven innovation.

## 🚀 Platform Vision

Samskipnad is becoming a **platform-first** system based on the principles outlined in our [Re-Architecting Roadmap](Re-Architecting-Roadmap.md). The transformation includes:

- **🏗️ Abstraction Layered Architecture (ALA)**: Stable core services with explicit contracts
- **🔌 Plugin Ecosystem**: Safe, isolated extensibility via HashiCorp's go-plugin
- **🎨 Creator Studio**: Tiered customization from YAML configs to full plugins
- **🌐 Multi-Community Platform**: White-label solution supporting diverse communities

### Core Architectural Principles

1. **Stability over Volatility**: Core system provides rock-solid foundation
2. **Explicit over Implicit**: All interactions through well-defined APIs
3. **Composition over Inheritance**: Build features by combining stable components
4. **Safety and Isolation by Default**: Plugin failures cannot crash the core system


## 🎯 Current Implementation Status

**Phase**: Transitioning from MVP to Platform Architecture  
**Architecture**: Legacy monolithic → Core Services Layer (in progress)  
**Plugin System**: Not implemented (Phase 2 target)  
**Creator Studio**: Basic YAML configuration (Tier 1 foundation)

> ⚠️ **Important**: This platform is currently undergoing architectural transformation. The existing codebase serves as a foundation while we implement the new plugin-based architecture outlined in our [Re-Architecting Roadmap](Re-Architecting-Roadmap.md).

## 🏗️ Platform Architecture Overview

### Three-Phase Transformation

#### **Phase 1: Foundation** (Current Focus)
- **Core Services Layer**: Refactor existing logic into stable, versioned interfaces
- **Declarative Customization**: YAML-based theming and configuration
- **Abstraction Layer**: Implement formal decoupling boundaries

#### **Phase 2: Plugin System** (Next Target)  
- **go-plugin Integration**: Safe, isolated plugin architecture
- **Plugin SDK**: Developer tools and documentation
- **Core Service APIs**: gRPC interfaces for plugin communication

#### **Phase 3: Creator Ecosystem** (Future Vision)
- **Creator Studio**: Plugin marketplace and management UI
- **Community Validation**: Plugin review and approval process
- **Ecosystem Growth**: Third-party developer community

### Current Core Services (Legacy → Refactored)

| Service | Legacy Status | Refactoring Status | Target Interface |
|---------|---------------|-------------------|------------------|
| **UserProfileService** | ✅ Implemented | 🔄 In Progress | User management, profiles, authentication |
| **CommunityManagementService** | ✅ Implemented | 🔄 In Progress | Multi-tenant community configuration |
| **ItemManagementService** | ⚠️ Partial | ❌ Pending | Classes, bookings, content management |
| **EventBusService** | ❌ Missing | ❌ Pending | Asynchronous messaging between components |
| **PaymentService** | ⚠️ Partial | ❌ Pending | Stripe integration, subscriptions, billing |

## 🚀 Quick Start (Current MVP)

> **Note**: These instructions are for the current MVP implementation. As we transition to the platform architecture, the setup process will evolve to support plugin-based customization.

### Prerequisites

- Go 1.21 or later
- SQLite (included)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/helloellinor/samskipnad.git
cd samskipnad
```

2. Install dependencies:
```bash
make deps
```

3. Run with your community configuration:
```bash
# Default: Kjernekraft (Scandinavian fitness)
make run

# Or choose a different community:
COMMUNITY=serenity make run    # Traditional yoga studio
COMMUNITY=yourcommunity make run
```

4. Create your own community:
```bash
cp config/kjernekraft.yaml config/mycommunity.yaml
# Edit config/mycommunity.yaml with your branding
COMMUNITY=mycommunity make run
```

The application will be available at `http://localhost:8080`

### Default Login

After first run, you can login with:
- **Email**: admin@kjernekraft.no
- **Password**: admin (change in production!)

### Tier 1 Customization (Current)

Each community is defined by a YAML configuration file in the `config/` directory. This represents the foundation of our Tier 1 declarative customization system. See [COMMUNITY_CONFIG.md](COMMUNITY_CONFIG.md) for the complete customization guide.

**Example community switch**:
```bash
# Scandinavian fitness (edgy, modern)
COMMUNITY=kjernekraft ./server

# Traditional yoga (peaceful, warm)  
COMMUNITY=serenity ./server
```

## 📁 Project Structure

### Current Structure (Transitioning)
```
samskipnad/
├── cmd/server/               # Application entry point
├── config/                   # Community configuration files (Tier 1 foundation)
│   ├── kjernekraft.yaml     # Default Scandinavian fitness community
│   └── serenity.yaml        # Traditional yoga studio example
├── internal/                 # Private application code (being refactored to ALA)
│   ├── auth/                # Authentication & authorization → UserProfileService
│   ├── config/              # Configuration management → CommunityManagementService  
│   ├── database/            # Database connection & migrations
│   ├── handlers/            # HTTP handlers (Application Logic Layer)
│   ├── middleware/          # HTTP middleware
│   ├── models/              # Data models → Core Services Layer
│   └── payments/            # Payment processing → PaymentService
├── web/                     # Frontend assets (Presentation Layer)
│   ├── static/              # CSS, JS, images
│   └── templates/           # HTML templates
├── docs/                    # Documentation
├── tools/                   # Development tools
├── Re-Architecting-Roadmap.md  # Source of truth for platform transformation
├── ROADMAP.md               # Implementation roadmap
├── PROGRESS_TRACKER.md      # Development progress tracking
├── COMMUNITY_CONFIG.md      # Tier 1 customization guide
├── Makefile                 # Build commands
└── go.mod                   # Go modules
```

### Target Architecture (Post-Refactoring)
```
samskipnad/
├── cmd/server/              # Host application entry point
├── internal/
│   ├── core/                # Core Services Layer (stable interfaces)
│   │   ├── interfaces/      # Service interface definitions
│   │   ├── services/        # Default service implementations  
│   │   └── events/          # EventBusService implementation
│   ├── application/         # Application Logic Layer
│   ├── presentation/        # Presentation Layer
│   └── plugins/             # Plugin Host & Management
├── pkg/                     # Public APIs for plugins
│   ├── sdk/                 # Plugin SDK
│   └── interfaces/          # Public service interfaces
├── plugins/                 # Plugin directory
│   ├── community-templates/ # Community template plugins
│   ├── payment-providers/   # Payment integration plugins
│   └── analytics/           # Analytics plugins
├── config/                  # Tier 1 YAML configurations
├── web/                     # Presentation layer assets
└── creator-studio/          # Creator Studio implementation (Phase 3)
```

## 🔧 Key Technologies

### Current Stack
- **Backend**: Go with Gorilla Mux
- **Database**: SQLite with migrations  
- **Frontend**: HTMX + Bootstrap 5
- **Payments**: Stripe API
- **Authentication**: Session-based with bcrypt
- **Configuration**: YAML-based community configs

### Future Platform Stack
- **Plugin System**: HashiCorp go-plugin with gRPC
- **Core Services**: Versioned interfaces with ALA
- **Creator Tools**: Plugin SDK and CLI toolchain
- **Marketplace**: Plugin discovery and management UI
- **Security**: Process isolation, mTLS, validation pipeline

## 🛣️ Implementation Roadmap

Detailed implementation plans are tracked in:
- **[Re-Architecting Roadmap](Re-Architecting-Roadmap.md)**: Complete architectural transformation plan (source of truth)
- **[ROADMAP.md](ROADMAP.md)**: Three-phase implementation timeline  
- **[PROGRESS_TRACKER.md](PROGRESS_TRACKER.md)**: Current development status and metrics

### Phase Overview

| Phase | Focus | Duration | Key Deliverables |
|-------|-------|----------|------------------|
| **Phase 1** | Core Services Layer | 3-4 months | ALA refactoring, YAML hot-reload, stable interfaces |
| **Phase 2** | Plugin Architecture | 2-3 months | go-plugin integration, SDK, proof-of-concept plugin |
| **Phase 3** | Creator Ecosystem | 3-4 months | Creator Studio UI, plugin marketplace, validation pipeline |

## 🎨 Customization Tiers

### Tier 1: Declarative YAML Configuration (Current)
- **Target Users**: Community managers, designers, power users
- **Capabilities**: Theming, feature toggles, content customization
- **Technology**: YAML files with hot-reload
- **Example**: Change community branding without code

### Tier 2: Plugin Development (Phase 2)
- **Target Users**: Developers, integrators
- **Capabilities**: Custom business logic, third-party integrations, new features
- **Technology**: Go plugins with gRPC APIs
- **Example**: RSS feed importer, custom payment providers

### Tier 3: Creator Studio (Phase 3)
- **Target Users**: Non-technical administrators
- **Capabilities**: Plugin marketplace, one-click installs, configuration UIs
- **Technology**: Web-based management interface
- **Example**: Install Slack integration via marketplace

## ⚡ Development Commands

```bash
make build      # Build the application
make run        # Build and run with default community
make dev        # Development mode with hot reload
make test       # Run tests (Phase 1 priority: expand coverage)
make clean      # Clean build artifacts
make fmt        # Format code
make lint       # Lint code (requires golangci-lint)
make db-reset   # Reset database
make deps       # Download dependencies
make setup      # Setup development environment
```

### Community Development
```bash
# Run with specific community
COMMUNITY=serenity make run

# Create new community configuration
cp config/kjernekraft.yaml config/mycommunity.yaml
COMMUNITY=mycommunity make run
```

### Future Plugin Development (Phase 2)
```bash
# Plugin development commands (coming soon)
make plugin-scaffold name=myplugin    # Generate plugin template
make plugin-build plugin=myplugin     # Build specific plugin  
make plugin-test plugin=myplugin      # Test plugin
make plugin-install plugin=myplugin   # Install to local instance
```

## 🌍 Multi-Community Platform

Samskipnad is designed as a white-label platform supporting diverse communities:

### Current Communities
- **🏋️ Fitness Centers**: CrossFit, pilates, barre, strength training
- **🧘 Yoga Studios**: Traditional, hot yoga, prenatal, meditation
- **💻 Hackerspaces**: Community events, hackathons, workshops, fix parties  
- **🎭 Creative Spaces**: Art studios, maker spaces, craft workshops
- **📚 Learning Communities**: Language cafes, study groups, book clubs
- **🤝 Unions and Societies**: Meetups, professional groups, interest communities

### Platform Features
- **🏢 Isolated Communities**: Separate users, classes, and payments per community
- **🎨 Custom Branding**: Unique colors, fonts, content, and identity
- **💰 Flexible Pricing**: Different membership and class pricing structures
- **🌐 Regional Support**: Multiple currencies, languages, and timezones
- **🔧 Feature Toggles**: Enable/disable features per community needs
- **🔌 Extensibility**: Custom plugins for specialized community needs (Phase 2+)

## 🤝 Contributing to the Platform

### Current Contributors
- **Feature Development**: Implement Phase 1 core services refactoring
- **Testing**: Build comprehensive test coverage (critical need)
- **Documentation**: Improve guides and API documentation
- **Community Configs**: Create new community templates

### Future Contributors (Phase 2+)
- **Plugin Development**: Build plugins for the ecosystem
- **Creator Studio**: Contribute to the management UI
- **Marketplace**: Help with plugin validation and review
- **SDK Development**: Improve developer tools and experience

### Getting Started
1. Fork the repository
2. Check [PROGRESS_TRACKER.md](PROGRESS_TRACKER.md) for current priorities
3. Pick up an unassigned task from [ROADMAP.md](ROADMAP.md)
4. Follow the architectural principles in [Re-Architecting-Roadmap.md](Re-Architecting-Roadmap.md)
5. Add tests for new functionality
6. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details.

## 💡 Inspiration & Acknowledgments

This platform transformation draws inspiration from:
- **Community Platforms**: [Yogo.no](https://yogo.no), [Bruce Studios](https://www.brucestudios.com/nb)
- **Extensible Architectures**: HashiCorp's plugin ecosystem, WordPress plugin system
- **Academic Research**: Community Support Platform architectures and ALA patterns
- **Real Communities**: Built for [Kjernekraft Oslo](https://www.kjernekraftoslo.no) and growing

---

**🎯 Next Steps**: See [ROADMAP.md](ROADMAP.md) for Phase 1 implementation priorities and [PROGRESS_TRACKER.md](PROGRESS_TRACKER.md) for current development status.
