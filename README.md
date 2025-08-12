# Samskipnad - Configurable Community Platform

A lightweight, modern platform for managing fitness communities and wellness studios. Built with Go and HTMX for fast, interactive experiences without JavaScript complexity.

## ✨ Unique Features

- **🎨 Fully Configurable Branding**: Each community gets its own colors, fonts, content, and identity
- **📝 YAML-based Configuration**: Human-readable config files for easy customization
- **🌍 Multi-Community Support**: One platform, infinite possibilities
- **🎯 White-label Ready**: "Built with Samskipnad" attribution system

## 🏋️ Example Communities

### Kjernekraft (Default)
A Scandinavian fitness community for busy parents with an edgy, sarcastic tone. Features Nordic color palette and modern typography.

### Serenity Yoga Studio
A traditional yoga studio with warm, earthy colors and peaceful, mindful messaging.

**Configure your own**: See [COMMUNITY_CONFIG.md](COMMUNITY_CONFIG.md) for the complete guide.

## Features

- **🎨 Community Configuration**: YAML-based branding and content customization
- **🏢 Multi-tenant Support**: Communities can create their own instances
- **📅 Class Management**: Schedule and manage classes and events
- **🎫 Flexible Booking**: Support for both memberships and class tickets
- **👥 Role-based Permissions**: Admin, instructor, and member roles
- **💳 Payment Processing**: Stripe integration with configurable pricing
- **⚡ HTMX-powered UI**: Fast, interactive frontend without JavaScript
- **📱 Mobile Responsive**: Works seamlessly on all devices

## Quick Start

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

### Community Configuration

Each community is defined by a YAML configuration file in the `config/` directory. See [COMMUNITY_CONFIG.md](COMMUNITY_CONFIG.md) for the complete customization guide.

**Example community switch**:
```bash
# Scandinavian fitness (edgy, modern)
COMMUNITY=kjernekraft ./server

# Traditional yoga (peaceful, warm)  
COMMUNITY=serenity ./server
```

### Development

For development with hot reload:
```bash
make dev
```

## Default Login

After first run, you can login with:
- **Email**: admin@kjernekraft.no
- **Password**: admin (change in production!)

## Project Structure

```
samskipnad/
├── cmd/server/          # Application entry point
├── config/              # Community configuration files
│   ├── kjernekraft.yaml # Default Scandinavian fitness community
│   └── serenity.yaml    # Traditional yoga studio example
├── internal/            # Private application code
│   ├── auth/           # Authentication & authorization
│   ├── config/         # Configuration management
│   ├── database/       # Database connection & migrations
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # HTTP middleware
│   └── models/         # Data models
├── web/                # Frontend assets
│   ├── static/         # CSS, JS, images
│   └── templates/      # HTML templates
├── COMMUNITY_CONFIG.md  # Configuration guide
├── Makefile            # Build commands
└── go.mod              # Go modules
```

## Key Technologies

- **Backend**: Go with Gorilla Mux
- **Database**: SQLite with migrations
- **Frontend**: HTMX + Bootstrap 5
- **Payments**: Stripe API
- **Authentication**: Session-based with bcrypt

## API Endpoints

### Public Routes
- `GET /` - Homepage
- `GET /login` - Login page
- `POST /login` - Login submission
- `GET /register` - Registration page
- `POST /register` - Registration submission

### Protected Routes
- `GET /dashboard` - User dashboard
- `GET /classes` - Browse classes
- `POST /classes/{id}/book` - Book a class
- `GET /memberships` - View membership options
- `GET /profile` - User profile

### Admin Routes
- `GET /admin` - Admin dashboard
- `GET /admin/classes` - Manage classes
- `POST /admin/classes` - Create new class
- `GET /admin/users` - Manage users
- `GET /admin/roles` - Manage roles

## Configuration

### Environment Variables
- `PORT` - Server port (default: 8080)
- `COMMUNITY` - Community configuration to load (default: kjernekraft)
- `DATABASE_PATH` - SQLite database file path (default: ./samskipnad.db)
- `STRIPE_SECRET_KEY` - Stripe secret key for payments
- `STRIPE_PUBLISHABLE_KEY` - Stripe publishable key

### Community Configuration
Each community is configured via YAML files in the `config/` directory. See [COMMUNITY_CONFIG.md](COMMUNITY_CONFIG.md) for details on:

- **Branding**: Colors, fonts, logos, and visual identity
- **Content**: Welcome messages, feature descriptions, community voice
- **Pricing**: Membership costs, class prices, currency settings
- **Features**: Enable/disable platform features per community
- **Localization**: Language, timezone, and regional settings

## Development Commands

```bash
make build      # Build the application
make run        # Build and run
make dev        # Development mode with hot reload
make test       # Run tests
make clean      # Clean build artifacts
make fmt        # Format code
make lint       # Lint code (requires golangci-lint)
make db-reset   # Reset database
```

## Multi-Community Platform

Samskipnad is designed as a white-label platform supporting multiple communities:

- **Isolated Communities**: Each community has separate users, classes, and payments
- **Custom Branding**: Unique colors, fonts, content, and identity per community  
- **Flexible Pricing**: Different membership and class pricing per community
- **Regional Support**: Multiple currencies, languages, and timezones
- **Feature Toggles**: Enable/disable features per community needs

Perfect for:
- **Yoga Studios**: Traditional, hot yoga, prenatal, etc.
- **Fitness Centers**: CrossFit, pilates, barre, etc.  
- **Dance Studios**: Ballet, hip-hop, salsa, etc.
- **Wellness Centers**: Meditation, mindfulness, therapy, etc.
- **Corporate Programs**: Employee wellness initiatives

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if needed
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Inspiration

This project draws inspiration from:
- [Yogo.no](https://yogo.no) - Yoga studio management
- [Bruce Studios](https://www.brucestudios.com/nb) - Fitness platform
- [Bitraf](https://bitraf.no) - Community management

Built for the yoga community at [Kjernekraft Oslo](https://www.kjernekraftoslo.no).