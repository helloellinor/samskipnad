package impl_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"

	"samskipnad/internal/models"
	"samskipnad/internal/services/impl"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	
	// Create users table
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			phone TEXT,
			role TEXT DEFAULT 'member',
			active BOOLEAN DEFAULT TRUE,
			tenant_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`)
	require.NoError(t, err)
	
	return db
}

func TestUserProfileServiceImpl_Authenticate(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()
	
	service := impl.NewUserProfileService(db)
	ctx := context.Background()
	
	// Create a test user by registering
	testUser := &models.User{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		TenantID:  1,
	}
	
	err := service.Register(ctx, testUser)
	require.NoError(t, err)
	
	// Test successful authentication
	t.Run("ValidCredentials", func(t *testing.T) {
		user, err := service.Authenticate(ctx, "test@example.com", "defaultpassword")
		
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test", user.FirstName)
		assert.Equal(t, "User", user.LastName)
		assert.True(t, user.Active)
	})
	
	// Test invalid credentials
	t.Run("InvalidCredentials", func(t *testing.T) {
		user, err := service.Authenticate(ctx, "test@example.com", "wrongpassword")
		
		assert.Error(t, err)
		assert.Nil(t, user)
	})
	
	// Test non-existent user
	t.Run("NonExistentUser", func(t *testing.T) {
		user, err := service.Authenticate(ctx, "nonexistent@example.com", "password")
		
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserProfileServiceImpl_ProfileManagement(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()
	
	service := impl.NewUserProfileService(db)
	ctx := context.Background()
	
	// Create a test user
	testUser := &models.User{
		Email:     "profile@example.com",
		FirstName: "Profile",
		LastName:  "Test",
		TenantID:  1,
	}
	
	err := service.Register(ctx, testUser)
	require.NoError(t, err)
	
	// Get the user to find the ID
	user, err := service.Authenticate(ctx, "profile@example.com", "defaultpassword")
	require.NoError(t, err)
	userID := user.ID
	
	// Test GetProfile
	t.Run("GetProfile", func(t *testing.T) {
		profile, err := service.GetProfile(ctx, userID)
		
		assert.NoError(t, err)
		assert.Equal(t, "profile@example.com", profile.Email)
		assert.Equal(t, "Profile", profile.FirstName)
	})
	
	// Test UpdateProfile
	t.Run("UpdateProfile", func(t *testing.T) {
		updates := &models.User{
			FirstName: "Updated",
			LastName:  "Name",
			Phone:     "+47 123 45 678",
		}
		
		err := service.UpdateProfile(ctx, userID, updates)
		assert.NoError(t, err)
		
		// Verify the update
		profile, err := service.GetProfile(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated", profile.FirstName)
		assert.Equal(t, "Name", profile.LastName)
		assert.Equal(t, "+47 123 45 678", profile.Phone)
	})
	
	// Test AssignRole
	t.Run("AssignRole", func(t *testing.T) {
		err := service.AssignRole(ctx, userID, "instructor")
		assert.NoError(t, err)
		
		// Verify the role change
		profile, err := service.GetProfile(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, "instructor", profile.Role)
	})
	
	// Test GetUserRoles
	t.Run("GetUserRoles", func(t *testing.T) {
		roles, err := service.GetUserRoles(ctx, userID)
		
		assert.NoError(t, err)
		assert.Equal(t, []string{"instructor"}, roles)
	})
}

func TestUserProfileServiceImpl_Permissions(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()
	
	service := impl.NewUserProfileService(db)
	ctx := context.Background()
	
	// Create users with different roles
	adminUser := &models.User{
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		TenantID:  1,
	}
	memberUser := &models.User{
		Email:     "member@example.com",
		FirstName: "Member",
		LastName:  "User",
		TenantID:  1,
	}
	
	require.NoError(t, service.Register(ctx, adminUser))
	require.NoError(t, service.Register(ctx, memberUser))
	
	// Get user IDs
	admin, err := service.Authenticate(ctx, "admin@example.com", "defaultpassword")
	require.NoError(t, err)
	require.NoError(t, service.AssignRole(ctx, admin.ID, "admin"))
	
	member, err := service.Authenticate(ctx, "member@example.com", "defaultpassword")
	require.NoError(t, err)
	
	// Test admin permissions
	t.Run("AdminPermissions", func(t *testing.T) {
		permissions := []string{"manage_classes", "view_students", "manage_users", "manage_payments", "book_classes"}
		
		for _, permission := range permissions {
			hasPermission, err := service.HasPermission(ctx, admin.ID, permission)
			assert.NoError(t, err)
			assert.True(t, hasPermission, "Admin should have permission: %s", permission)
		}
	})
	
	// Test member permissions
	t.Run("MemberPermissions", func(t *testing.T) {
		// Member should only have book_classes permission
		hasBooking, err := service.HasPermission(ctx, member.ID, "book_classes")
		assert.NoError(t, err)
		assert.True(t, hasBooking)
		
		// Member should NOT have admin permissions
		hasManageUsers, err := service.HasPermission(ctx, member.ID, "manage_users")
		assert.NoError(t, err)
		assert.False(t, hasManageUsers)
		
		hasManageClasses, err := service.HasPermission(ctx, member.ID, "manage_classes")
		assert.NoError(t, err)
		assert.False(t, hasManageClasses)
	})
}

func TestUserProfileServiceImpl_SessionManagement(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()
	
	service := impl.NewUserProfileService(db)
	ctx := context.Background()
	
	// Create a test user
	testUser := &models.User{
		Email:     "session@example.com",
		FirstName: "Session",
		LastName:  "Test",
		TenantID:  1,
	}
	
	require.NoError(t, service.Register(ctx, testUser))
	user, err := service.Authenticate(ctx, "session@example.com", "defaultpassword")
	require.NoError(t, err)
	
	// Test session creation and validation
	t.Run("CreateAndValidateSession", func(t *testing.T) {
		// Create session
		sessionID, err := service.CreateSession(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, sessionID)
		
		// Validate session
		sessionUser, err := service.ValidateSession(ctx, sessionID)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, sessionUser.Email)
	})
	
	// Test session revocation
	t.Run("RevokeSession", func(t *testing.T) {
		// Create session
		sessionID, err := service.CreateSession(ctx, user.ID)
		require.NoError(t, err)
		
		// Revoke session
		err = service.RevokeSession(ctx, sessionID)
		assert.NoError(t, err)
		
		// Validate that session is no longer valid
		sessionUser, err := service.ValidateSession(ctx, sessionID)
		assert.Error(t, err)
		assert.Nil(t, sessionUser)
	})
}

func TestUserProfileServiceImpl_PasswordManagement(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()
	
	service := impl.NewUserProfileService(db)
	ctx := context.Background()
	
	// Create a test user
	testUser := &models.User{
		Email:     "password@example.com",
		FirstName: "Password",
		LastName:  "Test",
		TenantID:  1,
	}
	
	require.NoError(t, service.Register(ctx, testUser))
	user, err := service.Authenticate(ctx, "password@example.com", "defaultpassword")
	require.NoError(t, err)
	
	// Test password change
	t.Run("ChangePassword", func(t *testing.T) {
		err := service.ChangePassword(ctx, user.ID, "defaultpassword", "newpassword123")
		assert.NoError(t, err)
		
		// Verify old password no longer works
		_, err = service.Authenticate(ctx, "password@example.com", "defaultpassword")
		assert.Error(t, err)
		
		// Verify new password works
		updatedUser, err := service.Authenticate(ctx, "password@example.com", "newpassword123")
		assert.NoError(t, err)
		assert.Equal(t, user.Email, updatedUser.Email)
	})
	
	// Test password reset
	t.Run("ResetPassword", func(t *testing.T) {
		err := service.ResetPassword(ctx, "password@example.com")
		assert.NoError(t, err) // Should not error for existing user
		
		err = service.ResetPassword(ctx, "nonexistent@example.com")
		assert.Error(t, err) // Should error for non-existent user
	})
}