# Loading Screen Standards

## Overview
This document outlines the standards and implementation for loading indicators in the Samskipnad application.

## Problem Solved
Previously, the loading overlay was blocking the entire viewport with a full-screen overlay that prevented users from interacting with the site during HTMX requests. This created a poor user experience where users felt "locked out" of the interface.

## New Loading Indicator Design

### Visual Design
- **Position**: Fixed position in the top-right corner (20px from top and right edges)
- **Style**: Small, non-intrusive indicator with cyberpunk aesthetic
- **Background**: Semi-transparent black with neon green border
- **Size**: Compact (max-width: 200px) with backdrop blur effect

### Behavior
- **Show**: Only during active HTMX requests
- **Hide**: 200ms after request completion (brief delay to show completion)
- **Error handling**: Automatically hides on request errors
- **Accessibility**: Screen reader announcements without overwhelming users

### Technical Implementation

#### CSS Classes
```css
.loading-overlay {
    position: fixed;
    top: 20px;
    right: 20px;
    background: rgba(0, 0, 0, 0.8);
    border: 1px solid var(--neon-green);
    border-radius: 4px;
    padding: 12px 16px;
    display: none;
    z-index: 1000;
    backdrop-filter: blur(4px);
    max-width: 200px;
}
```

#### JavaScript Events
- `htmx:beforeRequest` - Show loading indicator
- `htmx:afterRequest` - Hide loading indicator (with 200ms delay)
- `htmx:responseError` - Hide loading indicator immediately

## Accessibility Features

### Screen Reader Support
- `role="status"` for loading container
- `aria-live="polite"` for non-intrusive announcements
- `aria-busy` state management
- Brief "Content updated" announcement on completion

### Visual Accessibility
- High contrast neon green against dark background
- Sufficient color contrast ratios
- Respects `prefers-reduced-motion` for animations
- Consistent with overall cyberpunk theme

## Implementation Guidelines

### Required HTML Structure
```html
<div class="loading-overlay" id="loading-overlay" role="status" aria-live="polite" aria-label="Loading content">
    <div class="loading-spinner">
        <div class="spinner-ring" aria-hidden="true"></div>
        <div class="loading-text">LOADING</div>
    </div>
</div>
```

### JavaScript Requirements
1. Must handle all three HTMX events: `beforeRequest`, `afterRequest`, `responseError`
2. Must manage `aria-busy` state
3. Must clean up accessibility announcements after timeout
4. Should include error handling for missing DOM elements

### CSS Requirements
1. Must not block user interaction with the main content
2. Must maintain cyberpunk aesthetic (neon green, Space Mono font)
3. Must respect accessibility preferences (`prefers-reduced-motion`)
4. Must be visible but unobtrusive

## Best Practices

### Do's
- Keep loading indicators small and non-blocking
- Provide brief, clear accessibility announcements
- Use consistent visual styling across all pages
- Include error handling for failed requests
- Respect user motion preferences

### Don'ts
- Never block the entire viewport with loading overlays
- Don't show loading indicators for instant operations
- Don't overwhelm screen readers with constant announcements
- Don't use loading indicators without proper accessibility attributes
- Don't ignore error states in loading behavior

## Testing Checklist

### Visual Testing
- [ ] Loading indicator appears in top-right corner during HTMX requests
- [ ] User can still interact with page content while loading
- [ ] Indicator disappears after request completion
- [ ] Styling matches cyberpunk theme
- [ ] Works on both desktop and mobile viewports

### Accessibility Testing
- [ ] Screen reader announces loading state appropriately
- [ ] `aria-busy` state changes correctly
- [ ] Loading indicator has proper ARIA labels
- [ ] Respects `prefers-reduced-motion` setting
- [ ] Keyboard navigation not interrupted by loading states

### Functionality Testing
- [ ] Shows on HTMX requests (`htmx:beforeRequest`)
- [ ] Hides on successful completion (`htmx:afterRequest`)
- [ ] Hides on errors (`htmx:responseError`)
- [ ] No memory leaks from announcement DOM elements
- [ ] Works with all HTMX-enabled navigation and forms

## Browser Support
- Modern browsers with CSS backdrop-filter support
- Graceful degradation for older browsers (no backdrop blur)
- Consistent behavior across Chrome, Firefox, Safari, Edge