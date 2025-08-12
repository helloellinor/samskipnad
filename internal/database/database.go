package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Init initializes the database connection and runs migrations
func Init() (*sql.DB, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./samskipnad.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB) error {
	migrations := []string{
		createTenantsTable,
		createUsersTable,
		createClassesTable,
		createBookingsTable,
		createMembershipsTable,
		createTicketsTable,
		createPaymentsTable,
		createRolesTable,
		insertDefaultData,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			log.Printf("Migration failed: %s", migration)
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

const createTenantsTable = `
CREATE TABLE IF NOT EXISTS tenants (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	slug TEXT UNIQUE NOT NULL,
	domain TEXT,
	description TEXT,
	active BOOLEAN DEFAULT true,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	phone TEXT,
	role TEXT DEFAULT 'member' CHECK (role IN ('admin', 'instructor', 'member')),
	active BOOLEAN DEFAULT true,
	tenant_id INTEGER NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);`

const createClassesTable = `
CREATE TABLE IF NOT EXISTS classes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	tenant_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	description TEXT,
	instructor_id INTEGER NOT NULL,
	start_time DATETIME NOT NULL,
	end_time DATETIME NOT NULL,
	max_capacity INTEGER DEFAULT 20,
	price INTEGER DEFAULT 0,
	requires_ticket BOOLEAN DEFAULT false,
	requires_membership BOOLEAN DEFAULT false,
	active BOOLEAN DEFAULT true,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (tenant_id) REFERENCES tenants(id),
	FOREIGN KEY (instructor_id) REFERENCES users(id)
);`

const createBookingsTable = `
CREATE TABLE IF NOT EXISTS bookings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	class_id INTEGER NOT NULL,
	status TEXT DEFAULT 'confirmed' CHECK (status IN ('confirmed', 'cancelled', 'waitlist')),
	payment_id TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (class_id) REFERENCES classes(id),
	UNIQUE(user_id, class_id)
);`

const createMembershipsTable = `
CREATE TABLE IF NOT EXISTS memberships (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	tenant_id INTEGER NOT NULL,
	type TEXT NOT NULL CHECK (type IN ('monthly', 'yearly', 'unlimited')),
	start_date DATETIME NOT NULL,
	end_date DATETIME NOT NULL,
	active BOOLEAN DEFAULT true,
	payment_id TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);`

const createTicketsTable = `
CREATE TABLE IF NOT EXISTS tickets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	tenant_id INTEGER NOT NULL,
	type TEXT NOT NULL CHECK (type IN ('single', 'pack_5', 'pack_10')),
	classes_left INTEGER NOT NULL,
	expiry_date DATETIME NOT NULL,
	payment_id TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);`

const createPaymentsTable = `
CREATE TABLE IF NOT EXISTS payments (
	id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	tenant_id INTEGER NOT NULL,
	amount INTEGER NOT NULL,
	currency TEXT DEFAULT 'usd',
	status TEXT NOT NULL,
	payment_type TEXT NOT NULL CHECK (payment_type IN ('class', 'membership', 'ticket')),
	reference_id INTEGER NOT NULL,
	stripe_data TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);`

const createRolesTable = `
CREATE TABLE IF NOT EXISTS roles (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	tenant_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	description TEXT,
	permissions TEXT, -- JSON array
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (tenant_id) REFERENCES tenants(id),
	UNIQUE(tenant_id, name)
);`

const insertDefaultData = `
INSERT OR IGNORE INTO tenants (id, name, slug, description) VALUES 
(1, 'Kjernekraft Oslo', 'kjernekraft', 'High-end yoga gym in Oslo');

INSERT OR IGNORE INTO users (id, email, password_hash, first_name, last_name, role, tenant_id) VALUES 
(1, 'admin@kjernekraft.no', '$2a$10$728fKDj4WMprz9OZw3u0qu2JXgpbGX1Lxw1PIVc7y..juGQc8Xc4u', 'Admin', 'User', 'admin', 1);

INSERT OR IGNORE INTO roles (tenant_id, name, description, permissions) VALUES 
(1, 'admin', 'Full system access', '["read", "write", "delete", "manage_users", "manage_payments"]'),
(1, 'instructor', 'Manage classes and view students', '["read", "write_classes", "view_students"]'),
(1, 'member', 'Book classes and manage profile', '["read", "book_classes", "manage_profile"]');
`