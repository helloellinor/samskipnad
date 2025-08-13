package handlers_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"samskipnad/internal/handlers"
	"samskipnad/internal/services"
	"samskipnad/internal/services/mocks"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestHandlerContainerWithMockServices demonstrates that the refactored handlers
// work with the Core Services Layer and can use mock implementations for testing
func TestHandlerContainerWithMockServices(t *testing.T) {
	// Create mock services - this proves the abstraction layer is working
	mockServices := &services.ServiceContainer{
		UserProfile:         &mocks.MockUserProfileService{},
		CommunityManagement: &mocks.MockCommunityManagementService{},
		ItemManagement:      &mocks.MockItemManagementService{},
		Payment:             &mocks.MockPaymentService{},
		EventBus:            &mocks.MockEventBusService{},
		PluginHost:          nil, // Phase 2
	}

	// Create minimal templates for testing
	templates := template.New("")
	template.Must(templates.New("login.html").Parse(`<h1>Login</h1>`))
	template.Must(templates.New("register.html").Parse(`<h1>Register</h1>`))
	template.Must(templates.New("dashboard.html").Parse(`<h1>{{.Title}}</h1>`))
	template.Must(templates.New("classes.html").Parse(`<h1>{{.Title}}</h1>`))
	template.Must(templates.New("profile.html").Parse(`<h1>{{.Title}}</h1>`))
	template.Must(templates.New("bookings.html").Parse(`<h1>{{.Title}}</h1>`))

	// Create handler container using mock services
	handlerContainer := handlers.NewHandlerContainer(mockServices, templates)

	// Create router and setup routes
	router := mux.NewRouter()
	handlerContainer.SetupRoutes(router)

	t.Run("HomePageLoads", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Core Services Layer Active")
		assert.Contains(t, rr.Body.String(), "AuthHandlers")
		assert.Contains(t, rr.Body.String(), "ClassHandlers")
		assert.Contains(t, rr.Body.String(), "UserHandlers")
	})

	t.Run("LoginPageLoads", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/login", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "<h1>Login</h1>")
	})

	t.Run("ClassesPageLoads", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/classes", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "<h1>Classes</h1>")
	})

	t.Run("RoutesAreConfigured", func(t *testing.T) {
		// Test that all expected routes exist by trying to match them
		routes := []string{
			"/",
			"/login",
			"/register",
			"/logout",
			"/dashboard",
			"/profile",
			"/my-bookings",
			"/classes",
			"/book-class",
			"/cancel-booking",
			"/search-classes",
			"/admin/classes",
			"/admin/classes/edit",
			"/admin/classes/delete",
		}

		for _, route := range routes {
			req, err := http.NewRequest("GET", route, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// We expect either 200 (working), 405 (method not allowed, but route exists),
			// 501 (not implemented, but route exists), or 500 (template error but route exists)
			assert.True(t, 
				rr.Code == http.StatusOK || 
				rr.Code == http.StatusMethodNotAllowed || 
				rr.Code == http.StatusNotImplemented ||
				rr.Code == http.StatusInternalServerError,
				"Route %s should exist (got %d)", route, rr.Code)
		}
	})
}