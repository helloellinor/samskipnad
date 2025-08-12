# Samskipnad - Yoga Community Platform

A lightweight, modern platform for managing yoga studios and fitness communities. Built with Go and HTMX for fast, interactive experiences without JavaScript complexity.

## Features

- **Multi-tenant Support**: Communities can create their own instances
- **Class Management**: Schedule and manage yoga classes and events
- **Flexible Booking**: Support for both memberships and class tickets
- **Role-based Permissions**: Admin, instructor, and member roles
- **Payment Processing**: Stripe integration for secure payments
- **HTMX-powered UI**: Fast, interactive frontend without JavaScript
- **Mobile Responsive**: Works seamlessly on all devices

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

3. Build and run:
```bash
make run
```

The application will be available at `http://localhost:8080`

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
├── internal/            # Private application code
│   ├── auth/           # Authentication & authorization
│   ├── database/       # Database connection & migrations
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # HTTP middleware
│   └── models/         # Data models
├── web/                # Frontend assets
│   ├── static/         # CSS, JS, images
│   └── templates/      # HTML templates
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

Environment variables:
- `PORT` - Server port (default: 8080)
- `DATABASE_PATH` - SQLite database file path (default: ./samskipnad.db)
- `STRIPE_SECRET_KEY` - Stripe secret key for payments
- `STRIPE_PUBLISHABLE_KEY` - Stripe publishable key

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

## Multi-tenant Usage

The platform supports multiple communities (tenants). Each tenant has:
- Isolated user base
- Separate class schedules
- Independent payment processing
- Custom branding (planned)

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