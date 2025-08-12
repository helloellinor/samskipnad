package impl

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"samskipnad/internal/auth"
	"samskipnad/internal/models"
	"samskipnad/internal/services"

	"golang.org/x/crypto/bcrypt"
)

// UserProfileServiceImpl provides a concrete implementation of UserProfileService
// that wraps the existing auth.Service to maintain backward compatibility
// while establishing the Core Services Layer abstraction boundary
type UserProfileServiceImpl struct {
	db        *sql.DB
	authSvc   *auth.Service
	sessions  map[string]*models.User // Simple in-memory session store for now
	secretKey []byte
}

// NewUserProfileService creates a new UserProfileService implementation
func NewUserProfileService(db *sql.DB) services.UserProfileService {
	authSvc := auth.NewService(db)
	return &UserProfileServiceImpl{
		db:        db,
		authSvc:   authSvc,
		sessions:  make(map[string]*models.User),
		secretKey: []byte("super-secret-key-change-in-production"), // TODO: Move to config
	}
}

// Authenticate implements the UserProfileService interface
func (s *UserProfileServiceImpl) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	return s.authSvc.Login(email, password)
}

// Register implements the UserProfileService interface
func (s *UserProfileServiceImpl) Register(ctx context.Context, user *models.User) error {
	// Extract password from context or user (for now, assume it's in a special field)
	// In a real implementation, this would be handled more securely
	password := "defaultpassword" // TODO: This needs to be handled properly
	
	_, err := s.authSvc.Register(user.Email, password, user.FirstName, user.LastName, user.TenantID)
	return err
}

// GetProfile implements the UserProfileService interface
func (s *UserProfileServiceImpl) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	return s.authSvc.GetUserByID(userID)
}

// UpdateProfile implements the UserProfileService interface
func (s *UserProfileServiceImpl) UpdateProfile(ctx context.Context, userID int, updates *models.User) error {
	// Build dynamic SQL update based on non-empty fields
	query := "UPDATE users SET updated_at = ?"
	args := []interface{}{time.Now()}
	
	if updates.FirstName != "" {
		query += ", first_name = ?"
		args = append(args, updates.FirstName)
	}
	if updates.LastName != "" {
		query += ", last_name = ?"
		args = append(args, updates.LastName)
	}
	if updates.Email != "" {
		query += ", email = ?"
		args = append(args, updates.Email)
	}
	if updates.Phone != "" {
		query += ", phone = ?"
		args = append(args, updates.Phone)
	}
	if updates.Role != "" {
		query += ", role = ?"
		args = append(args, updates.Role)
	}
	
	query += " WHERE id = ?"
	args = append(args, userID)
	
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteProfile implements the UserProfileService interface
func (s *UserProfileServiceImpl) DeleteProfile(ctx context.Context, userID int) error {
	// Soft delete by setting active = false
	_, err := s.db.ExecContext(ctx, "UPDATE users SET active = false, updated_at = ? WHERE id = ?", time.Now(), userID)
	return err
}

// AssignRole implements the UserProfileService interface
func (s *UserProfileServiceImpl) AssignRole(ctx context.Context, userID int, role string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE users SET role = ?, updated_at = ? WHERE id = ?", role, time.Now(), userID)
	return err
}

// GetUserRoles implements the UserProfileService interface
func (s *UserProfileServiceImpl) GetUserRoles(ctx context.Context, userID int) ([]string, error) {
	user, err := s.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// For now, return the single role as an array
	// In the future, this could support multiple roles
	return []string{user.Role}, nil
}

// HasPermission implements the UserProfileService interface
func (s *UserProfileServiceImpl) HasPermission(ctx context.Context, userID int, permission string) (bool, error) {
	user, err := s.GetProfile(ctx, userID)
	if err != nil {
		return false, err
	}
	
	// Admin has all permissions
	if user.Role == "admin" {
		return true, nil
	}
	
	// Permission logic from existing auth service
	switch permission {
	case "manage_classes":
		return user.Role == "admin" || user.Role == "instructor", nil
	case "view_students":
		return user.Role == "admin" || user.Role == "instructor", nil
	case "manage_users":
		return user.Role == "admin", nil
	case "manage_payments":
		return user.Role == "admin", nil
	case "book_classes":
		return true, nil // All authenticated users can book classes
	default:
		return false, nil
	}
}

// CreateSession implements the UserProfileService interface
func (s *UserProfileServiceImpl) CreateSession(ctx context.Context, userID int) (string, error) {
	user, err := s.GetProfile(ctx, userID)
	if err != nil {
		return "", err
	}
	
	// Generate a simple session ID (in production, use crypto/rand)
	sessionID := fmt.Sprintf("session_%d_%d", userID, time.Now().Unix())
	
	// Store session (in production, use Redis or database)
	s.sessions[sessionID] = user
	
	return sessionID, nil
}

// ValidateSession implements the UserProfileService interface
func (s *UserProfileServiceImpl) ValidateSession(ctx context.Context, sessionID string) (*models.User, error) {
	user, exists := s.sessions[sessionID]
	if !exists {
		return nil, auth.ErrUnauthorized
	}
	
	return user, nil
}

// RevokeSession implements the UserProfileService interface
func (s *UserProfileServiceImpl) RevokeSession(ctx context.Context, sessionID string) error {
	delete(s.sessions, sessionID)
	return nil
}

// ChangePassword implements the UserProfileService interface
func (s *UserProfileServiceImpl) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	// Get current user to verify old password
	user, err := s.GetProfile(ctx, userID)
	if err != nil {
		return err
	}
	
	// Verify old password
	if !s.authSvc.CheckPassword(oldPassword, user.PasswordHash) {
		return auth.ErrInvalidCredentials
	}
	
	// Hash new password
	hashedPassword, err := s.hashPassword(newPassword)
	if err != nil {
		return err
	}
	
	// Update password in database
	_, err = s.db.ExecContext(ctx, "UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?", 
		hashedPassword, time.Now(), userID)
	return err
}

// ResetPassword implements the UserProfileService interface
func (s *UserProfileServiceImpl) ResetPassword(ctx context.Context, email string) error {
	// TODO: Implement password reset logic
	// This would typically involve:
	// 1. Generate a secure reset token
	// 2. Store token with expiration
	// 3. Send email with reset link
	// 4. Allow password reset with valid token
	
	// For now, just check if user exists
	_, err := s.authSvc.GetUserByEmail(email)
	return err
}

// Helper method for password hashing
func (s *UserProfileServiceImpl) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// GetAuthService returns the underlying auth service for backward compatibility
// This method should be removed once all handlers are refactored to use the interface
func (s *UserProfileServiceImpl) GetAuthService() *auth.Service {
	return s.authSvc
}