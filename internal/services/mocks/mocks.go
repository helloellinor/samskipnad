package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"samskipnad/internal/config"
	"samskipnad/internal/models"
	"samskipnad/internal/services"
)

// MockUserProfileService provides a mock implementation of UserProfileService
type MockUserProfileService struct {
	mock.Mock
}

func (m *MockUserProfileService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserProfileService) Register(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserProfileService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserProfileService) UpdateProfile(ctx context.Context, userID int, updates *models.User) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *MockUserProfileService) DeleteProfile(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserProfileService) AssignRole(ctx context.Context, userID int, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockUserProfileService) GetUserRoles(ctx context.Context, userID int) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserProfileService) HasPermission(ctx context.Context, userID int, permission string) (bool, error) {
	args := m.Called(ctx, userID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserProfileService) CreateSession(ctx context.Context, userID int) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockUserProfileService) ValidateSession(ctx context.Context, sessionID string) (*models.User, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserProfileService) RevokeSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockUserProfileService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserProfileService) ResetPassword(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

// MockCommunityManagementService provides a mock implementation of CommunityManagementService
type MockCommunityManagementService struct {
	mock.Mock
}

func (m *MockCommunityManagementService) GetCommunity(ctx context.Context, tenantID int) (*models.Tenant, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockCommunityManagementService) LoadConfiguration(ctx context.Context, communitySlug string) (*config.Community, error) {
	args := m.Called(ctx, communitySlug)
	return args.Get(0).(*config.Community), args.Error(1)
}

func (m *MockCommunityManagementService) UpdateConfiguration(ctx context.Context, tenantID int, config *config.Community) error {
	args := m.Called(ctx, tenantID, config)
	return args.Error(0)
}

func (m *MockCommunityManagementService) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockCommunityManagementService) GetTenantBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockCommunityManagementService) ListTenants(ctx context.Context) ([]*models.Tenant, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Tenant), args.Error(1)
}

func (m *MockCommunityManagementService) AddMember(ctx context.Context, tenantID, userID int, role string) error {
	args := m.Called(ctx, tenantID, userID, role)
	return args.Error(0)
}

func (m *MockCommunityManagementService) RemoveMember(ctx context.Context, tenantID, userID int) error {
	args := m.Called(ctx, tenantID, userID)
	return args.Error(0)
}

func (m *MockCommunityManagementService) GetMembers(ctx context.Context, tenantID int) ([]*models.User, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockCommunityManagementService) GetMemberRole(ctx context.Context, tenantID, userID int) (string, error) {
	args := m.Called(ctx, tenantID, userID)
	return args.String(0), args.Error(1)
}

func (m *MockCommunityManagementService) GetSettings(ctx context.Context, tenantID int) (map[string]interface{}, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockCommunityManagementService) UpdateSettings(ctx context.Context, tenantID int, settings map[string]interface{}) error {
	args := m.Called(ctx, tenantID, settings)
	return args.Error(0)
}

func (m *MockCommunityManagementService) IsFeatureEnabled(ctx context.Context, tenantID int, feature string) (bool, error) {
	args := m.Called(ctx, tenantID, feature)
	return args.Bool(0), args.Error(1)
}

func (m *MockCommunityManagementService) EnableFeature(ctx context.Context, tenantID int, feature string) error {
	args := m.Called(ctx, tenantID, feature)
	return args.Error(0)
}

func (m *MockCommunityManagementService) DisableFeature(ctx context.Context, tenantID int, feature string) error {
	args := m.Called(ctx, tenantID, feature)
	return args.Error(0)
}

// MockItemManagementService provides a mock implementation of ItemManagementService
type MockItemManagementService struct {
	mock.Mock
}

func (m *MockItemManagementService) CreateItem(ctx context.Context, tenantID int, itemType string, data interface{}) (int, error) {
	args := m.Called(ctx, tenantID, itemType, data)
	return args.Int(0), args.Error(1)
}

func (m *MockItemManagementService) GetItem(ctx context.Context, itemID int) (interface{}, error) {
	args := m.Called(ctx, itemID)
	return args.Get(0), args.Error(1)
}

func (m *MockItemManagementService) UpdateItem(ctx context.Context, itemID int, data interface{}) error {
	args := m.Called(ctx, itemID, data)
	return args.Error(0)
}

func (m *MockItemManagementService) DeleteItem(ctx context.Context, itemID int) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockItemManagementService) SearchItems(ctx context.Context, tenantID int, itemType string, filters map[string]interface{}) ([]interface{}, error) {
	args := m.Called(ctx, tenantID, itemType, filters)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockItemManagementService) ListItems(ctx context.Context, tenantID int, itemType string, limit, offset int) ([]interface{}, error) {
	args := m.Called(ctx, tenantID, itemType, limit, offset)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockItemManagementService) AddTag(ctx context.Context, itemID int, tag string) error {
	args := m.Called(ctx, itemID, tag)
	return args.Error(0)
}

func (m *MockItemManagementService) RemoveTag(ctx context.Context, itemID int, tag string) error {
	args := m.Called(ctx, itemID, tag)
	return args.Error(0)
}

func (m *MockItemManagementService) GetTags(ctx context.Context, itemID int) ([]string, error) {
	args := m.Called(ctx, itemID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockItemManagementService) CreateClass(ctx context.Context, class *models.Class) error {
	args := m.Called(ctx, class)
	return args.Error(0)
}

func (m *MockItemManagementService) GetClass(ctx context.Context, classID int) (*models.Class, error) {
	args := m.Called(ctx, classID)
	return args.Get(0).(*models.Class), args.Error(1)
}

func (m *MockItemManagementService) UpdateClass(ctx context.Context, classID int, class *models.Class) error {
	args := m.Called(ctx, classID, class)
	return args.Error(0)
}

func (m *MockItemManagementService) ListClasses(ctx context.Context, tenantID int, filters map[string]interface{}) ([]*models.Class, error) {
	args := m.Called(ctx, tenantID, filters)
	return args.Get(0).([]*models.Class), args.Error(1)
}

func (m *MockItemManagementService) CreateBooking(ctx context.Context, booking *models.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockItemManagementService) GetBooking(ctx context.Context, bookingID int) (*models.Booking, error) {
	args := m.Called(ctx, bookingID)
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockItemManagementService) CancelBooking(ctx context.Context, bookingID int) error {
	args := m.Called(ctx, bookingID)
	return args.Error(0)
}

func (m *MockItemManagementService) ListBookings(ctx context.Context, userID int) ([]*models.Booking, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Booking), args.Error(1)
}

// MockPaymentService provides a mock implementation of PaymentService
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) ProcessPayment(ctx context.Context, userID, tenantID int, amount int, currency, source string) (*models.Payment, error) {
	args := m.Called(ctx, userID, tenantID, amount, currency, source)
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentService) GetPayment(ctx context.Context, paymentID string) (*models.Payment, error) {
	args := m.Called(ctx, paymentID)
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentService) RefundPayment(ctx context.Context, paymentID string, amount int) error {
	args := m.Called(ctx, paymentID, amount)
	return args.Error(0)
}

func (m *MockPaymentService) CreateSubscription(ctx context.Context, userID, tenantID int, planID string) (*models.Membership, error) {
	args := m.Called(ctx, userID, tenantID, planID)
	return args.Get(0).(*models.Membership), args.Error(1)
}

func (m *MockPaymentService) CancelSubscription(ctx context.Context, subscriptionID string) error {
	args := m.Called(ctx, subscriptionID)
	return args.Error(0)
}

func (m *MockPaymentService) GetSubscription(ctx context.Context, userID, tenantID int) (*models.Membership, error) {
	args := m.Called(ctx, userID, tenantID)
	return args.Get(0).(*models.Membership), args.Error(1)
}

func (m *MockPaymentService) PurchaseKlippekort(ctx context.Context, userID, tenantID int, categoryID string, packageIndex int) (*models.Klippekort, error) {
	args := m.Called(ctx, userID, tenantID, categoryID, packageIndex)
	return args.Get(0).(*models.Klippekort), args.Error(1)
}

func (m *MockPaymentService) UseKlipp(ctx context.Context, userID, tenantID int, categoryID string) error {
	args := m.Called(ctx, userID, tenantID, categoryID)
	return args.Error(0)
}

func (m *MockPaymentService) GetKlippekortBalance(ctx context.Context, userID, tenantID int, categoryID string) (int, error) {
	args := m.Called(ctx, userID, tenantID, categoryID)
	return args.Int(0), args.Error(1)
}

func (m *MockPaymentService) HandleWebhook(ctx context.Context, provider string, payload []byte) error {
	args := m.Called(ctx, provider, payload)
	return args.Error(0)
}

func (m *MockPaymentService) GenerateInvoice(ctx context.Context, paymentID string) ([]byte, error) {
	args := m.Called(ctx, paymentID)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPaymentService) GetInvoice(ctx context.Context, invoiceID string) ([]byte, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).([]byte), args.Error(1)
}

// MockEventBusService provides a mock implementation of EventBusService
type MockEventBusService struct {
	mock.Mock
}

func (m *MockEventBusService) Publish(ctx context.Context, event *services.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBusService) PublishAsync(ctx context.Context, event *services.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBusService) Subscribe(ctx context.Context, eventType string, handler services.EventHandler) error {
	args := m.Called(ctx, eventType, handler)
	return args.Error(0)
}

func (m *MockEventBusService) Unsubscribe(ctx context.Context, eventType string, handler services.EventHandler) error {
	args := m.Called(ctx, eventType, handler)
	return args.Error(0)
}

func (m *MockEventBusService) SendNotification(ctx context.Context, userID int, notification *services.Notification) error {
	args := m.Called(ctx, userID, notification)
	return args.Error(0)
}

func (m *MockEventBusService) SendEmail(ctx context.Context, to, subject, body string) error {
	args := m.Called(ctx, to, subject, body)
	return args.Error(0)
}

func (m *MockEventBusService) SendSMS(ctx context.Context, to, message string) error {
	args := m.Called(ctx, to, message)
	return args.Error(0)
}

func (m *MockEventBusService) LogEvent(ctx context.Context, event *services.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBusService) GetEventHistory(ctx context.Context, filters map[string]interface{}) ([]*services.Event, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*services.Event), args.Error(1)
}