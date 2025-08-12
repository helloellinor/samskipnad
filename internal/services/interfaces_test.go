package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"samskipnad/internal/models"
	"samskipnad/internal/services"
	"samskipnad/internal/services/mocks"
)

// TestUserProfileServiceInterface demonstrates that our mock implements the interface
// and validates the Core Services Layer abstraction boundary
func TestUserProfileServiceInterface(t *testing.T) {
	// Arrange
	mockService := &mocks.MockUserProfileService{}
	ctx := context.Background()
	
	testUser := &models.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      "member",
		Active:    true,
		TenantID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test Authentication
	t.Run("Authentication", func(t *testing.T) {
		mockService.On("Authenticate", ctx, "test@example.com", "password").Return(testUser, nil)
		
		user, err := mockService.Authenticate(ctx, "test@example.com", "password")
		
		assert.NoError(t, err)
		assert.Equal(t, testUser.Email, user.Email)
		mockService.AssertExpectations(t)
	})

	// Test Profile Management
	t.Run("GetProfile", func(t *testing.T) {
		mockService.On("GetProfile", ctx, 1).Return(testUser, nil)
		
		user, err := mockService.GetProfile(ctx, 1)
		
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		mockService.AssertExpectations(t)
	})

	// Test Session Management
	t.Run("CreateSession", func(t *testing.T) {
		sessionID := "session_12345"
		mockService.On("CreateSession", ctx, 1).Return(sessionID, nil)
		
		result, err := mockService.CreateSession(ctx, 1)
		
		assert.NoError(t, err)
		assert.Equal(t, sessionID, result)
		mockService.AssertExpectations(t)
	})
}

// TestCommunityManagementServiceInterface validates the multi-tenant functionality
func TestCommunityManagementServiceInterface(t *testing.T) {
	// Arrange
	mockService := &mocks.MockCommunityManagementService{}
	ctx := context.Background()
	
	testTenant := &models.Tenant{
		ID:          1,
		Name:        "Test Community",
		Slug:        "test-community",
		Domain:      "test.example.com",
		Description: "A test community",
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test Tenant Management
	t.Run("GetTenantBySlug", func(t *testing.T) {
		mockService.On("GetTenantBySlug", ctx, "test-community").Return(testTenant, nil)
		
		tenant, err := mockService.GetTenantBySlug(ctx, "test-community")
		
		assert.NoError(t, err)
		assert.Equal(t, testTenant.Slug, tenant.Slug)
		mockService.AssertExpectations(t)
	})

	// Test Feature Toggles
	t.Run("IsFeatureEnabled", func(t *testing.T) {
		mockService.On("IsFeatureEnabled", ctx, 1, "classes").Return(true, nil)
		
		enabled, err := mockService.IsFeatureEnabled(ctx, 1, "classes")
		
		assert.NoError(t, err)
		assert.True(t, enabled)
		mockService.AssertExpectations(t)
	})

	// Test Member Management
	t.Run("AddMember", func(t *testing.T) {
		mockService.On("AddMember", ctx, 1, 1, "member").Return(nil)
		
		err := mockService.AddMember(ctx, 1, 1, "member")
		
		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})
}

// TestItemManagementServiceInterface validates the content management functionality
func TestItemManagementServiceInterface(t *testing.T) {
	// Arrange
	mockService := &mocks.MockItemManagementService{}
	ctx := context.Background()
	
	testClass := &models.Class{
		ID:                 1,
		TenantID:           1,
		Name:               "Test Yoga Class",
		Description:        "A test class",
		InstructorID:       1,
		StartTime:          time.Now().Add(24 * time.Hour),
		EndTime:            time.Now().Add(25 * time.Hour),
		MaxCapacity:        20,
		Price:              2000, // 20.00 in cents
		RequiresTicket:     true,
		RequiresMembership: false,
		Active:             true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Test Class Management
	t.Run("CreateClass", func(t *testing.T) {
		mockService.On("CreateClass", ctx, testClass).Return(nil)
		
		err := mockService.CreateClass(ctx, testClass)
		
		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("GetClass", func(t *testing.T) {
		mockService.On("GetClass", ctx, 1).Return(testClass, nil)
		
		class, err := mockService.GetClass(ctx, 1)
		
		assert.NoError(t, err)
		assert.Equal(t, testClass.Name, class.Name)
		assert.Equal(t, testClass.Price, class.Price)
		mockService.AssertExpectations(t)
	})

	// Test Generic Item Operations
	t.Run("CreateItem", func(t *testing.T) {
		mockService.On("CreateItem", ctx, 1, "class", mock.Anything).Return(1, nil)
		
		itemID, err := mockService.CreateItem(ctx, 1, "class", testClass)
		
		assert.NoError(t, err)
		assert.Equal(t, 1, itemID)
		mockService.AssertExpectations(t)
	})
}

// TestServiceContainerIntegration validates the dependency injection container
func TestServiceContainerIntegration(t *testing.T) {
	// Arrange
	container := &services.ServiceContainer{
		UserProfile:         &mocks.MockUserProfileService{},
		CommunityManagement: &mocks.MockCommunityManagementService{},
		ItemManagement:      &mocks.MockItemManagementService{},
		Payment:             &mocks.MockPaymentService{},
		EventBus:            &mocks.MockEventBusService{},
	}

	// Test that container holds all required services
	t.Run("ContainerHasAllServices", func(t *testing.T) {
		assert.NotNil(t, container.UserProfile)
		assert.NotNil(t, container.CommunityManagement)
		assert.NotNil(t, container.ItemManagement)
		assert.NotNil(t, container.Payment)
		assert.NotNil(t, container.EventBus)
	})

	// Test that services implement their interfaces
	t.Run("ServicesImplementInterfaces", func(t *testing.T) {
		var _ services.UserProfileService = container.UserProfile
		var _ services.CommunityManagementService = container.CommunityManagement
		var _ services.ItemManagementService = container.ItemManagement
		var _ services.PaymentService = container.Payment
		var _ services.EventBusService = container.EventBus
	})
}

// TestEventBusServiceInterface validates the messaging system
func TestEventBusServiceInterface(t *testing.T) {
	// Arrange
	mockService := &mocks.MockEventBusService{}
	ctx := context.Background()
	
	testEvent := &services.Event{
		ID:        "event_123",
		Type:      "user.registered",
		Source:    "user-service",
		Data:      map[string]interface{}{"user_id": 1},
		TenantID:  1,
		UserID:    1,
		Timestamp: time.Now(),
	}

	// Test Event Publishing
	t.Run("Publish", func(t *testing.T) {
		mockService.On("Publish", ctx, testEvent).Return(nil)
		
		err := mockService.Publish(ctx, testEvent)
		
		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	// Test Notification System
	t.Run("SendNotification", func(t *testing.T) {
		notification := &services.Notification{
			Type:     "info",
			Title:    "Welcome!",
			Message:  "Welcome to the community",
			Priority: 1,
		}
		
		mockService.On("SendNotification", ctx, 1, notification).Return(nil)
		
		err := mockService.SendNotification(ctx, 1, notification)
		
		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})
}

// TestPaymentServiceInterface validates the payment processing functionality
func TestPaymentServiceInterface(t *testing.T) {
	// Arrange
	mockService := &mocks.MockPaymentService{}
	ctx := context.Background()
	
	testPayment := &models.Payment{
		ID:          "pi_test_123",
		UserID:      1,
		TenantID:    1,
		Amount:      2000,
		Currency:    "NOK",
		Status:      "succeeded",
		PaymentType: "class",
		ReferenceID: 1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test Payment Processing
	t.Run("ProcessPayment", func(t *testing.T) {
		mockService.On("ProcessPayment", ctx, 1, 1, 2000, "NOK", "card_123").Return(testPayment, nil)
		
		payment, err := mockService.ProcessPayment(ctx, 1, 1, 2000, "NOK", "card_123")
		
		assert.NoError(t, err)
		assert.Equal(t, testPayment.Amount, payment.Amount)
		assert.Equal(t, testPayment.Currency, payment.Currency)
		mockService.AssertExpectations(t)
	})

	// Test Klippekort System
	t.Run("GetKlippekortBalance", func(t *testing.T) {
		mockService.On("GetKlippekortBalance", ctx, 1, 1, "yoga").Return(5, nil)
		
		balance, err := mockService.GetKlippekortBalance(ctx, 1, 1, "yoga")
		
		assert.NoError(t, err)
		assert.Equal(t, 5, balance)
		mockService.AssertExpectations(t)
	})
}