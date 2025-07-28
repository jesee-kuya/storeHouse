package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_"github.com/mattn/go-sqlite3"
)

// Database connection instance
var DB *sql.DB

// InitializeDatabase initializes the database connection and runs migrations
func InitializeDatabase(config *Config) error {
	var err error

	// Ensure database directory exists
	dbDir := filepath.Dir(config.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	DB, err = sql.Open("sqlite3", config.Database.Path+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(config.Database.MaxOpenConns)
	DB.SetMaxIdleConns(config.Database.MaxIdleConns)
	DB.SetConnMaxLifetime(config.Database.ConnMaxLifetime)

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := runMigrations(config.Database.MigrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Database initialized successfully at: %s", config.Database.Path)
	return nil
}

// GetDB returns the database connection instance
func GetDB() *sql.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call InitializeDatabase() first.")
	}
	return DB
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		log.Println("Closing database connection...")
		return DB.Close()
	}
	return nil
}

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	SQL         string
}

// runMigrations executes database migrations
func runMigrations(migrationsPath string) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(); err != nil {
		return err
	}

	// Get current migration version
	currentVersion, err := getCurrentMigrationVersion()
	if err != nil {
		return err
	}

	log.Printf("Current migration version: %d", currentVersion)

	// Get all available migrations
	migrations := getAllMigrations()

	// Run pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			log.Printf("Running migration %d: %s", migration.Version, migration.Description)

			if err := runSingleMigration(migration); err != nil {
				return fmt.Errorf("migration %d failed: %w", migration.Version, err)
			}

			log.Printf("Migration %d completed successfully", migration.Version)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := DB.Exec(query)
	return err
}

// getCurrentMigrationVersion gets the latest applied migration version
func getCurrentMigrationVersion() (int, error) {
	var version int
	query := "SELECT COALESCE(MAX(version), 0) FROM migrations"
	err := DB.QueryRow(query).Scan(&version)
	return version, err
}

// runSingleMigration executes a single migration
func runSingleMigration(migration Migration) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	_, err = tx.Exec(migration.SQL)
	if err != nil {
		return err
	}

	// Record migration in migrations table
	_, err = tx.Exec(
		"INSERT INTO migrations (version, description) VALUES (?, ?)",
		migration.Version,
		migration.Description,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// getAllMigrations returns all available migrations
func getAllMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Description: "Initial schema - create users, accounts tables",
			SQL: `
				-- Users table
				CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL,
					username TEXT UNIQUE NOT NULL,
					password TEXT NOT NULL,
					role TEXT NOT NULL CHECK (role IN ('admin', 'cashier')),
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);

				-- Accounts table
				CREATE TABLE IF NOT EXISTS accounts (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT UNIQUE NOT NULL,
					description TEXT,
					is_active BOOLEAN DEFAULT 1,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);

				-- Create default admin user (password: admin123)
				INSERT OR IGNORE INTO users (name, username, password, role) 
				VALUES ('Administrator', 'admin', '$2a$10$8K1p/a/3w2x6n5N5N5N5N5N5N5N5N5N5N5N5N5N5N5N5N5N5N5N5Ne', 'admin');

				-- Create default accounts
				INSERT OR IGNORE INTO accounts (name, description) VALUES 
				('Tithe', 'Regular tithe offerings'),
				('Church Building', 'Building fund contributions'),
				('Missions', 'Missionary support fund'),
				('Youth Ministry', 'Youth programs and activities'),
				('General Expenses', 'General church expenses');
			`,
		},
		{
			Version:     2,
			Description: "Create payments table",
			SQL: `
				-- Payments table
				CREATE TABLE IF NOT EXISTS payments (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					member_name TEXT NOT NULL,
					phone_number TEXT NOT NULL,
					amount REAL NOT NULL CHECK (amount > 0),
					account_id INTEGER NOT NULL,
					date TEXT NOT NULL,
					receipt_id TEXT UNIQUE NOT NULL,
					sms_sent BOOLEAN DEFAULT 0,
					notes TEXT,
					created_by INTEGER,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (account_id) REFERENCES accounts(id),
					FOREIGN KEY (created_by) REFERENCES users(id)
				);

				-- Create indexes for better performance
				CREATE INDEX IF NOT EXISTS idx_payments_date ON payments(date);
				CREATE INDEX IF NOT EXISTS idx_payments_account ON payments(account_id);
				CREATE INDEX IF NOT EXISTS idx_payments_member ON payments(member_name);
				CREATE INDEX IF NOT EXISTS idx_payments_receipt ON payments(receipt_id);
			`,
		},
		{
			Version:     3,
			Description: "Create expenses table",
			SQL: `
				-- Expenses table
				CREATE TABLE IF NOT EXISTS expenses (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					description TEXT NOT NULL,
					amount REAL NOT NULL CHECK (amount > 0),
					account_id INTEGER NOT NULL,
					date TEXT NOT NULL,
					receipt_number TEXT,
					vendor TEXT,
					notes TEXT,
					created_by INTEGER,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (account_id) REFERENCES accounts(id),
					FOREIGN KEY (created_by) REFERENCES users(id)
				);

				-- Create indexes for better performance
				CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(date);
				CREATE INDEX IF NOT EXISTS idx_expenses_account ON expenses(account_id);
				CREATE INDEX IF NOT EXISTS idx_expenses_description ON expenses(description);
			`,
		},
		{
			Version:     4,
			Description: "Add audit triggers for updated_at timestamps",
			SQL: `
				-- Trigger to update updated_at for users
				CREATE TRIGGER IF NOT EXISTS users_updated_at
				AFTER UPDATE ON users
				BEGIN
					UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
				END;

				-- Trigger to update updated_at for accounts
				CREATE TRIGGER IF NOT EXISTS accounts_updated_at
				AFTER UPDATE ON accounts
				BEGIN
					UPDATE accounts SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
				END;

				-- Trigger to update updated_at for payments
				CREATE TRIGGER IF NOT EXISTS payments_updated_at
				AFTER UPDATE ON payments
				BEGIN
					UPDATE payments SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
				END;

				-- Trigger to update updated_at for expenses
				CREATE TRIGGER IF NOT EXISTS expenses_updated_at
				AFTER UPDATE ON expenses
				BEGIN
					UPDATE expenses SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
				END;
			`,
		},
		{
			Version:     5,
			Description: "Create sessions table for JWT token management",
			SQL: `
				-- Sessions table for token management
				CREATE TABLE IF NOT EXISTS sessions (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER NOT NULL,
					token_hash TEXT UNIQUE NOT NULL,
					expires_at DATETIME NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
				);

				-- Create index for token lookups
				CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token_hash);
				CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);
				CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);

				-- Clean up expired sessions trigger
				CREATE TRIGGER IF NOT EXISTS cleanup_expired_sessions
				AFTER INSERT ON sessions
				BEGIN
					DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
				END;
			`,
		},
	}
}

// DatabaseStats represents database statistics
type DatabaseStats struct {
	TotalUsers    int
	TotalAccounts int
	TotalPayments int
	TotalExpenses int
	DatabaseSize  int64
}

// GetDatabaseStats returns current database statistics
func GetDatabaseStats() (*DatabaseStats, error) {
	stats := &DatabaseStats{}

	// Get table counts
	queries := map[string]*int{
		"SELECT COUNT(*) FROM users":    &stats.TotalUsers,
		"SELECT COUNT(*) FROM accounts": &stats.TotalAccounts,
		"SELECT COUNT(*) FROM payments": &stats.TotalPayments,
		"SELECT COUNT(*) FROM expenses": &stats.TotalExpenses,
	}

	for query, target := range queries {
		if err := DB.QueryRow(query).Scan(target); err != nil {
			return nil, fmt.Errorf("failed to get count: %w", err)
		}
	}

	// Get database file size
	if config := GetConfig(); config != nil {
		if fileInfo, err := os.Stat(config.Database.Path); err == nil {
			stats.DatabaseSize = fileInfo.Size()
		}
	}

	return stats, nil
}

// BackupDatabase creates a backup of the database
func BackupDatabase(backupPath string) error {
	config := GetConfig()
	if config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	// Ensure backup directory exists
	backupDir := filepath.Dir(backupPath)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup using SQLite backup API
	sourceDB, err := sql.Open("sqlite3", config.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to open source database: %w", err)
	}
	defer sourceDB.Close()

	backupDB, err := sql.Open("sqlite3", backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup database: %w", err)
	}
	defer backupDB.Close()

	// Perform the backup using VACUUM INTO (SQLite 3.27.0+)
	_, err = sourceDB.Exec("VACUUM INTO ?", backupPath)
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	log.Printf("Database backup completed: %s", backupPath)
	return nil
}

// OptimizeDatabase runs VACUUM and ANALYZE to optimize the database
func OptimizeDatabase() error {
	log.Println("Optimizing database...")

	// Run VACUUM to reclaim space and defragment
	if _, err := DB.Exec("VACUUM"); err != nil {
		return fmt.Errorf("VACUUM failed: %w", err)
	}

	// Run ANALYZE to update query planner statistics
	if _, err := DB.Exec("ANALYZE"); err != nil {
		return fmt.Errorf("ANALYZE failed: %w", err)
	}

	log.Println("Database optimization completed")
	return nil
}

// HealthCheck performs a basic database health check
func DatabaseHealthCheck() error {
	// Test connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Test a simple query
	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table'").Scan(&count); err != nil {
		return fmt.Errorf("health check query failed: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("no tables found in database")
	}

	return nil
}
