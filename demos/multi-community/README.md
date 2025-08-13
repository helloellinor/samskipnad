# Multi-Community Demo: Yoga Studio & Hackerspace

This demo showcases the Samskipnad platform's multi-tenant capabilities by demonstrating two distinct community types:

1. **Zen Flow Yoga Studio** - A modern mindful movement community
2. **Oslo Hackerspace** - A collaborative workspace for makers and creators

## What This Demo Shows

- **Multi-tenant Architecture**: Same platform, completely different experiences
- **Community-Specific Branding**: Unique visual identity and theming per community
- **Tailored Features**: Different pricing models, class types, and community features
- **Flexible Configuration**: YAML-driven customization without code changes
- **Hot-reload Configuration**: Real-time updates to community settings
- **Cross-Community User Management**: Users can be members of multiple communities

## Quick Start

1. **Build and run the demo:**
   ```bash
   ./run-demo.sh
   ```

2. **Access the communities:**
   - **Yoga Studio**: http://localhost:8080/?community=yoga-studio
   - **Hackerspace**: http://localhost:8080/?community=hackerspace

3. **Switch between communities:**
   - Use the community selector in the navigation
   - Or manually change the URL parameter

## Demo Features

### Zen Flow Yoga Studio
- **Wellness-focused design**: Calming colors, elegant typography
- **Diverse class offerings**: All levels, hot yoga, prenatal, workshops
- **Flexible pricing**: Drop-ins, class packages, unlimited memberships
- **Holistic approach**: Meditation, breathwork, philosophy workshops
- **Expert instruction**: Certified teachers with specialized training

### Oslo Hackerspace
- **Tech-inspired design**: Terminal-style colors, monospace fonts
- **Maker-focused features**: Workshop access, tool usage, project collaboration
- **Usage-based pricing**: 3D printing credits, laser cutter time, workshop passes
- **Learning opportunities**: Technical workshops, mentorship, hackathons
- **24/7 access**: Round-the-clock workspace for members

## Demo Data

### Yoga Studio Users
**Student:**
- Email: `student@zenflow.example.com`
- Password: `namaste123`

**Instructor:**
- Email: `teacher@zenflow.example.com`
- Password: `yoga123`

### Hackerspace Users
**Member:**
- Email: `maker@hackerspace.example.com`
- Password: `build123`

**Mentor:**
- Email: `mentor@hackerspace.example.com`
- Password: `hack123`

**Admin (both communities):**
- Email: `admin@example.com`
- Password: `admin123`

## Architecture Highlights

### Community Isolation
- Separate configuration files define each community's identity
- Isolated data per community while sharing infrastructure
- Dynamic CSS generation based on community settings
- Custom content and pricing per community

### Configuration-Driven Design
- YAML files control all aspects of community appearance and functionality
- Hot-reload allows real-time configuration changes
- No code changes needed to add new communities
- Flexible feature flags enable/disable functionality per community

### Cross-Community Features
- Single user account can access multiple communities
- Shared authentication and session management
- Community-specific roles and permissions
- Unified admin interface for managing multiple communities

## Technical Implementation

### Configuration Files
- `yoga-studio.yaml` - Yoga studio community configuration
- `hackerspace.yaml` - Hackerspace community configuration
- Both inherit from the same configuration schema

### Dynamic Theming
- CSS variables populated from community configuration
- Real-time theme switching without page reload
- Responsive design adapts to community branding

### Multi-Tenant Data
- Community-scoped data isolation
- Shared user accounts across communities
- Community-specific pricing and features

## Demo Scenarios

### Scenario 1: Yoga Student Journey
1. Visit the Zen Flow Yoga Studio
2. Browse available classes and teachers
3. Register for a beginner-friendly Hatha class
4. Purchase a 10-class package
5. Book additional classes using purchased credits

### Scenario 2: Hackerspace Member Journey
1. Visit the Oslo Hackerspace
2. Explore available workshops and equipment
3. Sign up for a 3D printing workshop
4. Purchase laser cutter time credits
5. Book 1-on-1 mentorship session

### Scenario 3: Admin Management
1. Login as admin user
2. Switch between community admin interfaces
3. Manage users, classes, and payments for each community
4. View community-specific analytics and reports

### Scenario 4: Configuration Changes
1. Edit community YAML files while demo is running
2. Observe real-time updates to branding and features
3. Add new class types or pricing options
4. Test hot-reload functionality

## What's Next

This demo represents the foundation for the plugin ecosystem. Future phases will add:

- **Plugin System**: Community-specific custom functionality
- **Creator Studio**: GUI for managing community configurations
- **Plugin Marketplace**: Discover and install community enhancements
- **Advanced Theming**: Visual theme builder and customization tools

## Files in This Demo

```
multi-community/
├── README.md                    # This file
├── run-demo.sh                  # Demo runner script
├── yoga-studio.yaml             # Yoga studio configuration
├── hackerspace.yaml             # Hackerspace configuration
└── demo-data/
    ├── users.sql                # Demo user accounts
    ├── yoga-classes.sql         # Sample yoga classes
    └── hackerspace-events.sql   # Sample workshops and events
```

## Exploring the Code

The multi-community demo showcases key architectural decisions:

1. **Configuration Schema**: See `internal/config/config.go` for the Community struct
2. **Hot-reload System**: Check `internal/config/hotreload.go` for real-time updates
3. **Dynamic CSS**: View `web/static/css/styles.css` for theming integration
4. **Community Switching**: Explore the middleware in `internal/middleware/`

This demonstration proves that Samskipnad can support radically different community types using the same core platform, setting the stage for the plugin ecosystem that will allow unlimited customization and extensibility.