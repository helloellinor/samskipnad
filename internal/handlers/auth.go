package handlers

import (
	"html/template"
	"net/http"

	"samskipnad/internal/services"
)

// AuthHandlers handles authentication-related HTTP requests
// This represents the Application Logic Layer for authentication operations
type AuthHandlers struct {
	userProfileService services.UserProfileService
	templates          *template.Template
}

// NewAuthHandlers creates a new AuthHandlers instance
func NewAuthHandlers(userProfileService services.UserProfileService, templates *template.Template) *AuthHandlers {
	return &AuthHandlers{
		userProfileService: userProfileService,
		templates:          templates,
	}
}

// Login handles user login
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Migrate login logic from monolithic handlers.go
	// This will use h.userProfileService.Authenticate() instead of direct database access
	
	if r.Method == "GET" {
		// Show login form
		err := h.templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == "POST" {
		// Process login
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := h.userProfileService.Authenticate(r.Context(), email, password)
		if err != nil {
			// Render login form with error
			data := struct {
				Error string
			}{
				Error: "Invalid email or password",
			}
			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
			return
		}

		// Create session
		sessionID, err := h.userProfileService.CreateSession(r.Context(), user.ID)
		if err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
		})

		// Redirect to dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Register handles user registration
func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Migrate registration logic from monolithic handlers.go
	// This will use h.userProfileService.Register() instead of direct database access
	
	if r.Method == "GET" {
		// Show registration form
		err := h.templates.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == "POST" {
		// TODO: Implement registration using userProfileService
		http.Error(w, "Registration not yet implemented in refactored version", http.StatusNotImplemented)
	}
}

// Logout handles user logout
func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session from cookie
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Revoke session
		h.userProfileService.RevokeSession(r.Context(), cookie.Value)
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}