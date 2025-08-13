package handlers

import (
	"html/template"
	"net/http"

	"samskipnad/internal/services"

	"github.com/gorilla/mux"
)

// HandlerContainer holds all domain-specific handlers
// This implements the Application Logic Layer pattern from the Re-Architecting Roadmap
type HandlerContainer struct {
	Auth    *AuthHandlers
	Classes *ClassHandlers
	User    *UserHandlers
	// TODO: Add PaymentHandlers, AdminHandlers as they're implemented
}

// NewHandlerContainer creates a new handler container with all domain handlers
func NewHandlerContainer(serviceContainer *services.ServiceContainer, templates *template.Template) *HandlerContainer {
	return &HandlerContainer{
		Auth: NewAuthHandlers(
			serviceContainer.UserProfile,
			templates,
		),
		Classes: NewClassHandlers(
			serviceContainer.ItemManagement,
			serviceContainer.CommunityManagement,
			serviceContainer.UserProfile,
			templates,
		),
		User: NewUserHandlers(
			serviceContainer.UserProfile,
			serviceContainer.ItemManagement,
			serviceContainer.CommunityManagement,
			templates,
		),
	}
}

// SetupRoutes configures all HTTP routes using the domain-specific handlers
func (h *HandlerContainer) SetupRoutes(router *mux.Router) {
	// Home page
	router.HandleFunc("/", h.Home).Methods("GET")

	// Authentication routes
	router.HandleFunc("/login", h.Auth.Login).Methods("GET", "POST")
	router.HandleFunc("/register", h.Auth.Register).Methods("GET", "POST")
	router.HandleFunc("/logout", h.Auth.Logout).Methods("POST")

	// User routes
	router.HandleFunc("/dashboard", h.User.Dashboard).Methods("GET")
	router.HandleFunc("/profile", h.User.Profile).Methods("GET", "POST")
	router.HandleFunc("/my-bookings", h.User.MyBookings).Methods("GET")

	// Class routes
	router.HandleFunc("/classes", h.Classes.Classes).Methods("GET")
	router.HandleFunc("/book-class", h.Classes.BookClass).Methods("POST")
	router.HandleFunc("/cancel-booking", h.Classes.CancelBooking).Methods("POST")
	router.HandleFunc("/search-classes", h.Classes.SearchClasses).Methods("GET")

	// Admin routes
	router.HandleFunc("/admin/classes", h.Classes.AdminClasses).Methods("GET", "POST")
	router.HandleFunc("/admin/classes/edit", h.Classes.EditClass).Methods("GET", "POST")
	router.HandleFunc("/admin/classes/delete", h.Classes.DeleteClass).Methods("POST")

	// TODO: Add payment routes when PaymentHandlers are implemented
	// TODO: Add other admin routes when AdminHandlers are implemented
}

// Home handles the home page - temporary implementation
func (h *HandlerContainer) Home(w http.ResponseWriter, r *http.Request) {
	// Simple home page showing the new architecture is active
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Samskipnad Platform</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { color: #2c3e50; border-bottom: 3px solid #3498db; padding-bottom: 20px; margin-bottom: 30px; }
        .status { background: #e8f5e8; border: 1px solid #4caf50; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .architecture { background: #f0f8ff; border: 1px solid #2196f3; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .links { margin-top: 30px; }
        .links a { display: inline-block; background: #3498db; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin: 5px; }
        .links a:hover { background: #2980b9; }
        .architecture ul { margin: 10px 0; padding-left: 20px; }
        .architecture li { margin: 5px 0; }
        .code { font-family: monospace; background: #f8f9fa; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎯 Samskipnad Platform</h1>
            <p>Creator-driven extensible community management platform</p>
        </div>
        
        <div class="status">
            <h3>✅ Core Services Layer Active</h3>
            <p>The re-architecting Phase 1 infrastructure is now operational!</p>
        </div>

        <div class="architecture">
            <h3>🏗️ Current Architecture Status</h3>
            <ul>
                <li><strong>✅ UserProfileService:</strong> Authentication & profile management</li>
                <li><strong>✅ CommunityManagementService:</strong> Multi-tenant configuration & features</li>
                <li><strong>🔄 ItemManagementService:</strong> Classes & content management (in progress)</li>
                <li><strong>🔄 PaymentService:</strong> Billing & subscriptions (planned)</li>
                <li><strong>✅ EventBusService:</strong> Messaging & notifications</li>
            </ul>
        </div>

        <div class="architecture">
            <h3>📋 Handler Refactoring Progress</h3>
            <ul>
                <li><strong>✅ AuthHandlers:</strong> Login, register, logout</li>
                <li><strong>✅ ClassHandlers:</strong> Class management & booking</li>
                <li><strong>✅ UserHandlers:</strong> Dashboard & profile</li>
                <li><strong>🔄 PaymentHandlers:</strong> Payment processing (planned)</li>
                <li><strong>🔄 AdminHandlers:</strong> Admin operations (planned)</li>
            </ul>
        </div>

        <div class="links">
            <h3>🚀 Available Routes</h3>
            <a href="/login">Login</a>
            <a href="/register">Register</a>
            <a href="/dashboard">Dashboard</a>
            <a href="/classes">Classes</a>
            <a href="/profile">Profile</a>
            <a href="/my-bookings">My Bookings</a>
        </div>

        <div class="architecture">
            <h3>📖 Documentation</h3>
            <p>For implementation details, see:</p>
            <ul>
                <li><span class="code">Re-Architecting-Roadmap.md</span> - Strategic vision</li>
                <li><span class="code">IMPLEMENTATION_SUMMARY.md</span> - Current progress</li>
                <li><span class="code">CODEBASE_ASSESSMENT.md</span> - Migration strategy</li>
            </ul>
        </div>
    </div>
</body>
</html>
    `))
}