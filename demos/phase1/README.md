# Phase 1 Demo: Core Samskipnad Platform

This demo showcases the core Samskipnad platform with its web interface, authentication system, community management, and payment processing capabilities.

## What This Demo Shows

- **Multi-tenant Community Management**: Dynamic configuration per community
- **User Authentication & Authorization**: Secure login/registration with role-based access
- **Class Booking System**: Interactive calendar with real-time availability
- **Payment Processing**: Membership purchases and Klippekort (punch card) system
- **Admin Interface**: Complete administrative controls for community managers
- **Responsive Web UI**: Modern interface using HTMX for dynamic interactions

## Quick Start

1. **Build and run the demo:**
   ```bash
   ./run-demo.sh
   ```

2. **Open your browser to:** http://localhost:8080

3. **Try the demo features:**
   - Register a new account or login with demo credentials
   - Browse the class calendar and book classes
   - Purchase memberships or klippekort
   - Explore the admin interface (admin credentials provided)

## Demo Features

### User Experience
- **Home Page**: Welcome with community branding
- **Registration/Login**: Secure authentication flow
- **Dashboard**: Personal activity overview
- **Class Calendar**: Interactive booking interface
- **Profile Management**: Update personal information
- **Payment Flow**: Secure membership and klippekort purchases

### Admin Experience
- **Admin Dashboard**: Platform management overview
- **User Management**: View and manage community members
- **Class Management**: Create, edit, and schedule classes
- **Payment Tracking**: Monitor transactions and memberships
- **Role Management**: Assign administrative permissions

### Technical Features
- **Dynamic CSS**: Community-specific theming and branding
- **Real-time Updates**: HTMX-powered dynamic content
- **Secure Sessions**: Proper authentication and authorization
- **Database Integration**: SQLite with automatic schema management
- **Mobile Responsive**: Works on all device sizes

## Demo Data

The demo includes sample data for:
- **Demo Community**: "Yoga Studio Oslo" configuration
- **Sample Classes**: Various yoga classes with different instructors
- **Demo Users**: Pre-created accounts for testing
- **Mock Payments**: Simulated payment flow without real processing

### Demo Credentials

**Regular User:**
- Email: `demo@example.com`
- Password: `demo123`

**Admin User:**
- Email: `admin@example.com`
- Password: `admin123`

## Architecture Highlights

This demo showcases the core architectural principles:

### Multi-tenant Design
- Community-specific configuration and branding
- Isolated data per community while sharing infrastructure
- Dynamic CSS generation based on community settings

### Service-Oriented Architecture
- Clean separation between auth, payment, and core services
- Dependency injection for service management
- Easy testing and mocking of service components

### Web-First Approach
- Server-side rendering with progressive enhancement
- HTMX for smooth user interactions without full page reloads
- Minimal JavaScript, maximum accessibility

### Security & Privacy
- Secure session management
- Role-based access control
- SQL injection prevention
- XSS protection

## Technical Stack

- **Backend**: Go with Gorilla Mux for routing
- **Database**: SQLite for simplicity (production would use PostgreSQL)
- **Frontend**: HTML templates with HTMX for interactivity
- **Styling**: CSS with dynamic community theming
- **Authentication**: Session-based with secure cookies
- **Payments**: Simulated payment processing (Stripe integration ready)

## Files in This Demo

- `run-demo.sh` - Main demo runner script
- `demo-config.yaml` - Community configuration for demo
- `setup-demo-data.sql` - Sample data for the demo
- `README.md` - This documentation

## Exploring the Code

The demo runs the actual platform code from:
- `../../cmd/server/` - Main application server
- `../../internal/` - Core business logic and services
- `../../web/` - HTML templates and static assets

To understand the implementation:
1. Start with `../../cmd/server/main.go` for the application entry point
2. Explore `../../internal/handlers/` for web request handling
3. Review `../../internal/services/` for business logic
4. Check `../../web/templates/` for the user interface

## What's Next

This Phase 1 demo provides the foundation for:
- **Phase 2**: Plugin system that extends this core functionality
- **Phase 3**: Creator Studio for community customization
- **Future Phases**: Advanced features and integrations

The plugin system in Phase 2 will allow extending this core platform without modifying the base code, demonstrating true architectural flexibility.