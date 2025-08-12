# Community Configuration Guide

Samskipnad is a flexible platform that can be configured for different types of fitness and wellness communities. Each community has its own YAML configuration file that controls branding, content, pricing, and features.

## Quick Start

1. **Set your community**: Use the `COMMUNITY` environment variable
   ```bash
   COMMUNITY=kjernekraft ./server    # Norwegian fitness community (default)
   COMMUNITY=serenity ./server       # Traditional yoga studio
   ```

2. **Create your config**: Copy an existing config and customize it
   ```bash
   cp config/kjernekraft.yaml config/yourcommunity.yaml
   # Edit config/yourcommunity.yaml
   COMMUNITY=yourcommunity ./server
   ```

## Configuration Structure

### Basic Information
```yaml
name: "Your Community Name"
description: "Brief description"
tagline: "Your motto or tagline"
```

### Branding & Colors
```yaml
colors:
  primary: "#2e3440"      # Main brand color
  secondary: "#5e81ac"    # Secondary brand color  
  accent: "#d08770"       # Accent/highlight color
  success: "#a3be8c"      # Success state color
  warning: "#ebcb8b"      # Warning state color
  danger: "#bf616a"       # Error/danger color
  background: "#eceff4"   # Page background
  surface: "#f4f6fa"      # Card/surface background
  text: "#2e3440"         # Primary text color
  muted: "#4c566a"        # Muted/secondary text
```

### Typography
```yaml
fonts:
  primary: "Inter"        # Main font (Google Fonts)
  secondary: "JetBrains Mono"  # Accent/code font
  size_base: "16px"       # Base font size
```

### Features
```yaml
features:
  classes: true           # Class booking system
  memberships: true       # Membership management
  community: true         # Community features
  payments: true          # Stripe payments
  calendar: true          # Calendar view
```

### Content
```yaml
content:
  home:
    title: "Welcome to Your Community"
    subtitle: "Optional subtitle"
    description: "Main description text"
    
  features:
    - title: "Feature 1"
      description: "Description of feature 1"
    - title: "Feature 2" 
      description: "Description of feature 2"
    # Add up to 3 feature boxes
```

### Pricing
```yaml
pricing:
  currency: "NOK"         # Currency code
  monthly: 299            # Monthly membership price
  yearly: 2990            # Annual membership price  
  drop_in: 89             # Single class price
```

### Localization
```yaml
locale:
  language: "en"          # Language code
  country: "NO"           # Country code
  timezone: "Europe/Oslo" # Timezone
```

### Attribution
```yaml
attribution:
  show: true              # Show "Built with Samskipnad"
  text: "Built with Samskipnad"
  link: "https://github.com/helloellinor/samskipnad"
```

## Example Configurations

### Kjernekraft (Default)
A Scandinavian fitness community for busy parents with an edgy, slightly sarcastic tone. Uses Nordic color palette (blues, grays) and modern typography.

### Serenity Yoga Studio  
A traditional yoga studio with warm, earthy colors (browns, golds) and peaceful, mindful messaging.

## Creating Your Configuration

1. **Choose a name**: Pick a short, URL-friendly name for your community config file
2. **Copy a template**: Start with the config closest to your needs
3. **Customize branding**: Update colors, fonts, and visual identity
4. **Write your content**: Craft messaging that fits your community's voice
5. **Set pricing**: Configure membership and class pricing in your currency
6. **Test it**: Run with `COMMUNITY=yourname ./server`

## Design Philosophy

The configuration system is designed to let you create completely different experiences:

- **Corporate gym**: Modern, sleek, performance-focused
- **Yoga studio**: Peaceful, warm, mindfulness-oriented  
- **CrossFit box**: Bold, energetic, achievement-driven
- **Dance studio**: Creative, expressive, artistic
- **Senior center**: Accessible, friendly, community-focused

Each configuration creates a unique brand experience while using the same underlying platform.

## Advanced Customization

For deeper customization beyond the YAML config:

1. **Custom CSS**: The dynamic CSS system generates styles from your config
2. **Template overrides**: Create community-specific template variants
3. **Feature flags**: Enable/disable features per community
4. **Payment processors**: Configure different payment providers per region

## Contributing

When adding new configuration options:

1. Update the `Community` struct in `internal/config/config.go`
2. Add template usage in relevant HTML files
3. Update the CSS generation in `DynamicCSS` handler
4. Document the new option in this guide
5. Add examples to existing configurations

This system lets you build a "white-label" version of Samskipnad for any wellness community while maintaining the core functionality and user experience.