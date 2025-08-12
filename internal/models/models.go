package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID          int       `json:"id" db:"id"`
	Email       string    `json:"email" db:"email"`
	PasswordHash string   `json:"-" db:"password_hash"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Phone       string    `json:"phone" db:"phone"`
	Role        string    `json:"role" db:"role"` // admin, instructor, member
	Active      bool      `json:"active" db:"active"`
	TenantID    int       `json:"tenant_id" db:"tenant_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Tenant represents a community/organization instance
type Tenant struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Domain      string    `json:"domain" db:"domain"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Class represents a yoga class or event
type Class struct {
	ID          int       `json:"id" db:"id"`
	TenantID    int       `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	InstructorID int      `json:"instructor_id" db:"instructor_id"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	MaxCapacity int       `json:"max_capacity" db:"max_capacity"`
	Price       int       `json:"price" db:"price"` // in cents
	RequiresTicket bool   `json:"requires_ticket" db:"requires_ticket"`
	RequiresMembership bool `json:"requires_membership" db:"requires_membership"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Booking represents a user's booking for a class
type Booking struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ClassID   int       `json:"class_id" db:"class_id"`
	Status    string    `json:"status" db:"status"` // confirmed, cancelled, waitlist
	PaymentID string    `json:"payment_id" db:"payment_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Membership represents a user's membership
type Membership struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	TenantID  int       `json:"tenant_id" db:"tenant_id"`
	Type      string    `json:"type" db:"type"` // monthly, yearly, unlimited
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	Active    bool      `json:"active" db:"active"`
	PaymentID string    `json:"payment_id" db:"payment_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Ticket represents a class ticket/pass
type Ticket struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	TenantID    int       `json:"tenant_id" db:"tenant_id"`
	Type        string    `json:"type" db:"type"` // single, pack_5, pack_10
	ClassesLeft int       `json:"classes_left" db:"classes_left"`
	ExpiryDate  time.Time `json:"expiry_date" db:"expiry_date"`
	PaymentID   string    `json:"payment_id" db:"payment_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Payment represents a payment transaction
type Payment struct {
	ID            string    `json:"id" db:"id"` // Stripe payment intent ID
	UserID        int       `json:"user_id" db:"user_id"`
	TenantID      int       `json:"tenant_id" db:"tenant_id"`
	Amount        int       `json:"amount" db:"amount"` // in cents
	Currency      string    `json:"currency" db:"currency"`
	Status        string    `json:"status" db:"status"`
	PaymentType   string    `json:"payment_type" db:"payment_type"` // class, membership, ticket
	ReferenceID   int       `json:"reference_id" db:"reference_id"` // ID of class, membership, or ticket
	StripeData    string    `json:"stripe_data" db:"stripe_data"` // JSON data from Stripe
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Role represents a role in the system
type Role struct {
	ID          int       `json:"id" db:"id"`
	TenantID    int       `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Permissions string    `json:"permissions" db:"permissions"` // JSON array of permissions
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}