package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"vocabweb/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Get database connection string from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/vocabweb?sslmode=disable"
		fmt.Println("Using default DATABASE_URL (set DATABASE_URL env var to override)")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create migrator
	migrator := database.NewMigrator(db, "./migrations")

	// Execute command
	command := os.Args[1]
	switch command {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	case "status":
		if err := migrator.Status(); err != nil {
			log.Fatalf("Status check failed: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("VocabWeb Database Migration Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  migrate <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up      Apply all pending migrations")
	fmt.Println("  down    Rollback the last applied migration")
	fmt.Println("  status  Show current migration status")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  DATABASE_URL  PostgreSQL connection string")
	fmt.Println("                (default: postgres://postgres:postgres@localhost:5432/vocabweb?sslmode=disable)")
}
