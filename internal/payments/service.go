package payments

import (
	"database/sql"
	"fmt"
	"os"
	"samskipnad/internal/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/customer"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	// Initialize Stripe with API key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		stripe.Key = "sk_test_..." // Use a test key for development
	}
	
	return &Service{
		db: db,
	}
}

// CreatePaymentIntent creates a Stripe payment intent for a class booking
func (s *Service) CreatePaymentIntent(userID int, classID int, amount int64) (*stripe.PaymentIntent, error) {
	// Get user details for customer creation
	user, err := s.getUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Create or get Stripe customer
	stripeCustomer, err := s.getOrCreateStripeCustomer(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	// Create payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Customer: stripe.String(stripeCustomer.ID),
		Metadata: map[string]string{
			"user_id":  fmt.Sprintf("%d", userID),
			"class_id": fmt.Sprintf("%d", classID),
			"type":     "class_booking",
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Store payment record in database
	err = s.storePayment(&models.Payment{
		ID:          pi.ID,
		UserID:      userID,
		TenantID:    user.TenantID,
		Amount:      int(amount),
		Currency:    "usd",
		Status:      string(pi.Status),
		PaymentType: "class",
		ReferenceID: classID,
		StripeData:  "", // We could store the full JSON here if needed
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store payment: %w", err)
	}

	return pi, nil
}

// CreateMembershipPaymentIntent creates a payment intent for membership
func (s *Service) CreateMembershipPaymentIntent(userID int, membershipType string, amount int64) (*stripe.PaymentIntent, error) {
	user, err := s.getUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	stripeCustomer, err := s.getOrCreateStripeCustomer(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Customer: stripe.String(stripeCustomer.ID),
		Metadata: map[string]string{
			"user_id":         fmt.Sprintf("%d", userID),
			"membership_type": membershipType,
			"type":           "membership",
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Store payment record
	err = s.storePayment(&models.Payment{
		ID:          pi.ID,
		UserID:      userID,
		TenantID:    user.TenantID,
		Amount:      int(amount),
		Currency:    "usd",
		Status:      string(pi.Status),
		PaymentType: "membership",
		ReferenceID: 0, // Will be set when membership is created
		StripeData:  "",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store payment: %w", err)
	}

	return pi, nil
}

// ConfirmPayment processes a successful payment
func (s *Service) ConfirmPayment(paymentIntentID string) error {
	// Get payment from database
	payment, err := s.getPaymentByID(paymentIntentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Get updated payment intent from Stripe
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return fmt.Errorf("failed to get payment intent: %w", err)
	}

	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return fmt.Errorf("payment not successful: %s", pi.Status)
	}

	// Update payment status
	err = s.updatePaymentStatus(paymentIntentID, string(pi.Status))
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Process the payment based on type
	switch payment.PaymentType {
	case "class":
		err = s.processClassBooking(payment.UserID, payment.ReferenceID)
	case "membership":
		err = s.processMembershipPayment(payment.UserID, pi.Metadata["membership_type"])
	default:
		return fmt.Errorf("unknown payment type: %s", payment.PaymentType)
	}

	if err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}

	return nil
}

// Helper methods

func (s *Service) getOrCreateStripeCustomer(user *models.User) (*stripe.Customer, error) {
	// Try to find existing customer first
	params := &stripe.CustomerListParams{
		Email: stripe.String(user.Email),
	}
	
	customers := customer.List(params)
	for customers.Next() {
		return customers.Customer(), nil
	}

	// Create new customer
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(fmt.Sprintf("%s %s", user.FirstName, user.LastName)),
		Metadata: map[string]string{
			"user_id":   fmt.Sprintf("%d", user.ID),
			"tenant_id": fmt.Sprintf("%d", user.TenantID),
		},
	}

	return customer.New(customerParams)
}

func (s *Service) getUserByID(userID int) (*models.User, error) {
	query := `SELECT id, email, first_name, last_name, phone, role, active, tenant_id, created_at, updated_at 
			  FROM users WHERE id = ?`
	
	var user models.User
	err := s.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Role, &user.Active, &user.TenantID,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (s *Service) storePayment(payment *models.Payment) error {
	query := `INSERT INTO payments (id, user_id, tenant_id, amount, currency, status, payment_type, reference_id, stripe_data, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`
	
	_, err := s.db.Exec(query, payment.ID, payment.UserID, payment.TenantID,
		payment.Amount, payment.Currency, payment.Status, payment.PaymentType,
		payment.ReferenceID, payment.StripeData)
	
	return err
}

func (s *Service) getPaymentByID(paymentID string) (*models.Payment, error) {
	query := `SELECT id, user_id, tenant_id, amount, currency, status, payment_type, reference_id, stripe_data, created_at, updated_at
			  FROM payments WHERE id = ?`
	
	var payment models.Payment
	err := s.db.QueryRow(query, paymentID).Scan(
		&payment.ID, &payment.UserID, &payment.TenantID, &payment.Amount,
		&payment.Currency, &payment.Status, &payment.PaymentType,
		&payment.ReferenceID, &payment.StripeData, &payment.CreatedAt, &payment.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &payment, nil
}

func (s *Service) updatePaymentStatus(paymentID, status string) error {
	query := `UPDATE payments SET status = ?, updated_at = datetime('now') WHERE id = ?`
	_, err := s.db.Exec(query, status, paymentID)
	return err
}

func (s *Service) processClassBooking(userID, classID int) error {
	// Create booking record
	query := `INSERT INTO bookings (user_id, class_id, status, created_at, updated_at)
			  VALUES (?, ?, 'confirmed', datetime('now'), datetime('now'))`
	
	_, err := s.db.Exec(query, userID, classID)
	return err
}

func (s *Service) processMembershipPayment(userID int, membershipType string) error {
	// Create membership record
	user, err := s.getUserByID(userID)
	if err != nil {
		return err
	}

	// Calculate membership dates based on type
	startDate := "datetime('now')"
	var endDate string
	switch membershipType {
	case "monthly":
		endDate = "datetime('now', '+1 month')"
	case "yearly":
		endDate = "datetime('now', '+1 year')"
	default:
		endDate = "datetime('now', '+1 month')" // Default to monthly
	}

	query := fmt.Sprintf(`INSERT INTO memberships (user_id, tenant_id, type, start_date, end_date, active, created_at, updated_at)
						  VALUES (?, ?, ?, %s, %s, true, datetime('now'), datetime('now'))`, startDate, endDate)
	
	_, err = s.db.Exec(query, userID, user.TenantID, membershipType)
	return err
}