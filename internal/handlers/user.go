package handlers

import (
	"html/template"
	"net/http"

	"samskipnad/internal/services"
)

// UserHandlers handles user dashboard and profile-related HTTP requests
// This represents the Application Logic Layer for user operations
type UserHandlers struct {
	userProfileService      services.UserProfileService
	itemManagementService   services.ItemManagementService
	communityService        services.CommunityManagementService
	templates               *template.Template
}

// NewUserHandlers creates a new UserHandlers instance
func NewUserHandlers(
	userProfileService services.UserProfileService,
	itemManagementService services.ItemManagementService,
	communityService services.CommunityManagementService,
	templates *template.Template,
) *UserHandlers {
	return &UserHandlers{
		userProfileService:    userProfileService,
		itemManagementService: itemManagementService,
		communityService:      communityService,
		templates:             templates,
	}
}

// Dashboard handles the user dashboard
func (h *UserHandlers) Dashboard(w http.ResponseWriter, r *http.Request) {
	// TODO: Migrate dashboard logic from monolithic handlers.go
	// This will use services to get user data instead of direct database access
	
	// For now, show a placeholder dashboard
	data := struct {
		Title   string
		Message string
	}{
		Title:   "Dashboard",
		Message: "Dashboard will be implemented using Core Services Layer",
	}

	err := h.templates.ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// Profile handles user profile viewing and editing
func (h *UserHandlers) Profile(w http.ResponseWriter, r *http.Request) {
	// TODO: Migrate profile logic from monolithic handlers.go
	// This will use h.userProfileService.GetProfile() and UpdateProfile() instead of direct database access
	
	if r.Method == "GET" {
		// Show profile form
		// TODO: Get current user from session and load profile data
		data := struct {
			Title   string
			Message string
		}{
			Title:   "Profile",
			Message: "Profile management will use UserProfileService",
		}

		err := h.templates.ExecuteTemplate(w, "profile.html", data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		// TODO: Update profile using userProfileService
		http.Error(w, "Profile updates not yet implemented in refactored version", http.StatusNotImplemented)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// MyBookings handles listing user's bookings
func (h *UserHandlers) MyBookings(w http.ResponseWriter, r *http.Request) {
	// TODO: Migrate bookings listing logic from monolithic handlers.go
	// This will use h.itemManagementService.ListBookings() instead of direct database access
	
	// TODO: Get current user from session
	// TODO: Load bookings using itemManagementService
	
	data := struct {
		Title   string
		Message string
	}{
		Title:   "My Bookings",
		Message: "Bookings listing will use ItemManagementService",
	}

	err := h.templates.ExecuteTemplate(w, "bookings.html", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}