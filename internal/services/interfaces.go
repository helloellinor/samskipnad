package services

import (
	"context"
	"time"

	"samskipnad/internal/config"
	"samskipnad/internal/models"
)

// Core Services Layer Interfaces
// These interfaces define the stable contracts for the Core Services Layer
// as outlined in the Re-Architecting Roadmap. They provide the abstraction
// boundary between the volatile customization layer and the stable core.

// UserProfileService handles all user-related operations including authentication,
// profile management, and role-based access control
type UserProfileService interface {
	// Authentication
	Authenticate(ctx context.Context, email, password string) (*models.User, error)
	Register(ctx context.Context, user *models.User) error
	
	// Profile Management
	GetProfile(ctx context.Context, userID int) (*models.User, error)
	UpdateProfile(ctx context.Context, userID int, updates *models.User) error
	DeleteProfile(ctx context.Context, userID int) error
	
	// Role and Permission Management
	AssignRole(ctx context.Context, userID int, role string) error
	GetUserRoles(ctx context.Context, userID int) ([]string, error)
	HasPermission(ctx context.Context, userID int, permission string) (bool, error)
	
	// Session Management
	CreateSession(ctx context.Context, userID int) (string, error)
	ValidateSession(ctx context.Context, sessionID string) (*models.User, error)
	RevokeSession(ctx context.Context, sessionID string) error
	
	// Password Management
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, email string) error
}

// CommunityManagementService handles multi-tenant community configuration,
// group structures, memberships, and permissions
type CommunityManagementService interface {
	// Community Configuration
	GetCommunity(ctx context.Context, tenantID int) (*models.Tenant, error)
	LoadConfiguration(ctx context.Context, communitySlug string) (*config.Community, error)
	UpdateConfiguration(ctx context.Context, tenantID int, config *config.Community) error
	
	// Multi-Tenant Management
	CreateTenant(ctx context.Context, tenant *models.Tenant) error
	GetTenantBySlug(ctx context.Context, slug string) (*models.Tenant, error)
	ListTenants(ctx context.Context) ([]*models.Tenant, error)
	
	// Member Management
	AddMember(ctx context.Context, tenantID, userID int, role string) error
	RemoveMember(ctx context.Context, tenantID, userID int) error
	GetMembers(ctx context.Context, tenantID int) ([]*models.User, error)
	GetMemberRole(ctx context.Context, tenantID, userID int) (string, error)
	
	// Community Settings
	GetSettings(ctx context.Context, tenantID int) (map[string]interface{}, error)
	UpdateSettings(ctx context.Context, tenantID int, settings map[string]interface{}) error
	
	// Feature Toggles
	IsFeatureEnabled(ctx context.Context, tenantID int, feature string) (bool, error)
	EnableFeature(ctx context.Context, tenantID int, feature string) error
	DisableFeature(ctx context.Context, tenantID int, feature string) error
}

// ItemManagementService handles CRUD operations for typed content objects
// including classes, events, articles, and their associated metadata
type ItemManagementService interface {
	// Generic Item Operations
	CreateItem(ctx context.Context, tenantID int, itemType string, data interface{}) (int, error)
	GetItem(ctx context.Context, itemID int) (interface{}, error)
	UpdateItem(ctx context.Context, itemID int, data interface{}) error
	DeleteItem(ctx context.Context, itemID int) error
	
	// Content Search and Discovery
	SearchItems(ctx context.Context, tenantID int, itemType string, filters map[string]interface{}) ([]interface{}, error)
	ListItems(ctx context.Context, tenantID int, itemType string, limit, offset int) ([]interface{}, error)
	
	// Metadata and Categorization
	AddTag(ctx context.Context, itemID int, tag string) error
	RemoveTag(ctx context.Context, itemID int, tag string) error
	GetTags(ctx context.Context, itemID int) ([]string, error)
	
	// Class-Specific Operations (backwards compatibility)
	CreateClass(ctx context.Context, class *models.Class) error
	GetClass(ctx context.Context, classID int) (*models.Class, error)
	UpdateClass(ctx context.Context, classID int, class *models.Class) error
	ListClasses(ctx context.Context, tenantID int, filters map[string]interface{}) ([]*models.Class, error)
	
	// Booking Operations
	CreateBooking(ctx context.Context, booking *models.Booking) error
	GetBooking(ctx context.Context, bookingID int) (*models.Booking, error)
	CancelBooking(ctx context.Context, bookingID int) error
	ListBookings(ctx context.Context, userID int) ([]*models.Booking, error)
}

// PaymentService provides abstraction over payment processing and billing
type PaymentService interface {
	// Payment Processing
	ProcessPayment(ctx context.Context, userID, tenantID int, amount int, currency, source string) (*models.Payment, error)
	GetPayment(ctx context.Context, paymentID string) (*models.Payment, error)
	RefundPayment(ctx context.Context, paymentID string, amount int) error
	
	// Subscription Management
	CreateSubscription(ctx context.Context, userID, tenantID int, planID string) (*models.Membership, error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
	GetSubscription(ctx context.Context, userID, tenantID int) (*models.Membership, error)
	
	// Credit/Klippekort System
	PurchaseKlippekort(ctx context.Context, userID, tenantID int, categoryID string, packageIndex int) (*models.Klippekort, error)
	UseKlipp(ctx context.Context, userID, tenantID int, categoryID string) error
	GetKlippekortBalance(ctx context.Context, userID, tenantID int, categoryID string) (int, error)
	
	// Webhook Handling
	HandleWebhook(ctx context.Context, provider string, payload []byte) error
	
	// Invoice Generation
	GenerateInvoice(ctx context.Context, paymentID string) ([]byte, error)
	GetInvoice(ctx context.Context, invoiceID string) ([]byte, error)
}

// EventBusService provides asynchronous messaging and decoupled communication
type EventBusService interface {
	// Event Publishing
	Publish(ctx context.Context, event *Event) error
	PublishAsync(ctx context.Context, event *Event) error
	
	// Event Subscription
	Subscribe(ctx context.Context, eventType string, handler EventHandler) error
	Unsubscribe(ctx context.Context, eventType string, handler EventHandler) error
	
	// Notification System
	SendNotification(ctx context.Context, userID int, notification *Notification) error
	SendEmail(ctx context.Context, to, subject, body string) error
	SendSMS(ctx context.Context, to, message string) error
	
	// Event Logging and Analytics
	LogEvent(ctx context.Context, event *Event) error
	GetEventHistory(ctx context.Context, filters map[string]interface{}) ([]*Event, error)
}

// Supporting Types for EventBusService

// Event represents a system event
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
	TenantID  int                    `json:"tenant_id"`
	UserID    int                    `json:"user_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventHandler defines the signature for event handlers
type EventHandler func(ctx context.Context, event *Event) error

// Notification represents a user notification
type Notification struct {
	Type     string                 `json:"type"`
	Title    string                 `json:"title"`
	Message  string                 `json:"message"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Priority int                    `json:"priority"`
}

// ServiceContainer provides dependency injection for all Core Services
// This will be used by the Application Logic Layer and Plugin System
type ServiceContainer struct {
	UserProfile         UserProfileService
	CommunityManagement CommunityManagementService
	ItemManagement      ItemManagementService
	Payment             PaymentService
	EventBus            EventBusService
	
	// Plugin system (Phase 2)
	PluginHost          PluginHostService
}

// PluginHostService manages the plugin system lifecycle
type PluginHostService interface {
	// Plugin Management
	LoadPlugin(ctx context.Context, name, path string) error
	UnloadPlugin(ctx context.Context, name string) error
	GetLoadedPlugins(ctx context.Context) []string
	
	// Plugin Execution
	ExecutePlugin(ctx context.Context, name string, params map[string]interface{}) (map[string]interface{}, error)
	
	// Lifecycle
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
}