package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"samskipnad/internal/models"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const sessionName = "samskipnad-session"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrUnauthorized       = errors.New("unauthorized")
)

type Service struct {
	db    *sql.DB
	store *sessions.CookieStore
}

func NewService(db *sql.DB) *Service {
	// In production, use a proper secret key from environment
	store := sessions.NewCookieStore([]byte("super-secret-key-change-in-production"))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	return &Service{
		db:    db,
		store: store,
	}
}

func (s *Service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *Service) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Service) Register(email, password, firstName, lastName string, tenantID int) (*models.User, error) {
	// Check if user already exists
	var exists bool
	err := s.db.QueryRow("SELECT 1 FROM users WHERE email = ?", email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Insert user
	result, err := s.db.Exec(`
		INSERT INTO users (email, password_hash, first_name, last_name, tenant_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		email, hashedPassword, firstName, lastName, tenantID, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetUserByID(int(id))
}

func (s *Service) Login(email, password string) (*models.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	if !user.Active {
		return nil, ErrUnauthorized
	}

	return user, nil
}

func (s *Service) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	var phone sql.NullString
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, first_name, last_name, phone, role, active, tenant_id, created_at, updated_at
		FROM users WHERE id = ?`, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&phone, &user.Role, &user.Active, &user.TenantID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if phone.Valid {
		user.Phone = phone.String
	}
	return user, nil
}

func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	var phone sql.NullString
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, first_name, last_name, phone, role, active, tenant_id, created_at, updated_at
		FROM users WHERE email = ?`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&phone, &user.Role, &user.Active, &user.TenantID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if phone.Valid {
		user.Phone = phone.String
	}
	return user, nil
}

func (s *Service) CreateSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Values["user_id"] = user.ID
	session.Values["user_role"] = user.Role
	session.Values["tenant_id"] = user.TenantID

	return session.Save(r, w)
}

func (s *Service) GetCurrentUser(r *http.Request) (*models.User, error) {
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return nil, ErrUnauthorized
	}

	return s.GetUserByID(userID)
}

func (s *Service) IsAuthenticated(r *http.Request) bool {
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		return false
	}

	_, ok := session.Values["user_id"].(int)
	return ok
}

func (s *Service) IsAdmin(r *http.Request) bool {
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		return false
	}

	role, ok := session.Values["user_role"].(string)
	return ok && role == "admin"
}

func (s *Service) DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1

	return session.Save(r, w)
}

func (s *Service) HasPermission(r *http.Request, permission string) bool {
	user, err := s.GetCurrentUser(r)
	if err != nil {
		return false
	}

	// Admin has all permissions
	if user.Role == "admin" {
		return true
	}

	// TODO: Implement role-based permissions from database
	// For now, simple role-based checks
	switch permission {
	case "manage_classes":
		return user.Role == "admin" || user.Role == "instructor"
	case "view_students":
		return user.Role == "admin" || user.Role == "instructor"
	case "manage_users":
		return user.Role == "admin"
	case "manage_payments":
		return user.Role == "admin"
	case "book_classes":
		return true // All authenticated users can book classes
	default:
		return false
	}
}