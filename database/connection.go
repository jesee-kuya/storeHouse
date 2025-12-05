package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL environment variable not set")
	}

	// Extract database name
	urlParts := strings.Split(dbURL, "/")
	if len(urlParts) < 4 {
		log.Fatal("❌ Invalid DATABASE_URL format")
	}
	dbName := strings.Split(urlParts[len(urlParts)-1], "?")[0]

	// Connect to the default 'postgres' database
	adminURL := strings.Replace(dbURL, dbName, "postgres", 1)
	adminDB, err := sql.Open("postgres", adminURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to admin database:", err)
	}
	defer adminDB.Close()

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = adminDB.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		log.Fatal("❌ Failed to check if database exists:", err)
	}

	// Create if not exists
	if !exists {
		fmt.Println("🛠️  Database not found. Creating:", dbName)
		_, err = adminDB.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			log.Fatal("❌ Failed to create database:", err)
		}
		fmt.Println("✅ Database created successfully.")
	}

	// Connect to the actual database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to target database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("❌ Failed to ping target database:", err)
	}

	// Connection pool tuning
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Connected to database:", dbName)
	return db
}

func ApplyMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Migrations applied successfully")
}
