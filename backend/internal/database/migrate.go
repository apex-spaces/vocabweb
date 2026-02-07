package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Migration represents a single database migration
type Migration struct {
	Version string
	UpSQL   string
	DownSQL string
}

// Migrator handles database migrations
type Migrator struct {
	db             *sql.DB
	migrationsPath string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, migrationsPath string) *Migrator {
	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
	}
}

// ensureMigrationsTable creates the schema_migrations table if it doesn't exist
func (m *Migrator) ensureMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	_, err := m.db.Exec(query)
	return err
}

// getAppliedMigrations returns a map of already applied migration versions
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// loadMigrations reads all migration files from the migrations directory
func (m *Migrator) loadMigrations() ([]Migration, error) {
	files, err := os.ReadDir(m.migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Group files by version
	migrationFiles := make(map[string]struct {
		up   string
		down string
	})

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		// Extract version from filename (e.g., "001_initial_schema.up.sql" -> "001")
		parts := strings.Split(name, "_")
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		entry := migrationFiles[version]

		if strings.Contains(name, ".up.sql") {
			entry.up = name
		} else if strings.Contains(name, ".down.sql") {
			entry.down = name
		}

		migrationFiles[version] = entry
	}

	// Convert to slice and sort by version
	var migrations []Migration
	for version, files := range migrationFiles {
		if files.up == "" {
			continue // Skip if no up migration
		}

		upSQL, err := os.ReadFile(filepath.Join(m.migrationsPath, files.up))
		if err != nil {
			return nil, fmt.Errorf("failed to read up migration %s: %w", files.up, err)
		}

		var downSQL []byte
		if files.down != "" {
			downSQL, err = os.ReadFile(filepath.Join(m.migrationsPath, files.down))
			if err != nil {
				return nil, fmt.Errorf("failed to read down migration %s: %w", files.down, err)
			}
		}

		migrations = append(migrations, Migration{
			Version: version,
			UpSQL:   string(upSQL),
			DownSQL: string(downSQL),
		})
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	if err := m.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	migrations, err := m.loadMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if applied[migration.Version] {
			fmt.Printf("Migration %s already applied, skipping\n", migration.Version)
			continue
		}

		fmt.Printf("Applying migration %s...\n", migration.Version)

		// Start transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Execute migration
		if _, err := tx.Exec(migration.UpSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
		}

		// Record migration
		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version, applied_at) VALUES ($1, $2)",
			migration.Version,
			time.Now(),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Version, err)
		}

		fmt.Printf("Migration %s applied successfully\n", migration.Version)
	}

	fmt.Println("All migrations applied successfully")
	return nil
}

// Down rolls back the last applied migration
func (m *Migrator) Down() error {
	if err := m.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	// Get the last applied migration
	var lastVersion string
	err := m.db.QueryRow(
		"SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1",
	).Scan(&lastVersion)
	if err == sql.ErrNoRows {
		fmt.Println("No migrations to roll back")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Load migrations
	migrations, err := m.loadMigrations()
	if err != nil {
		return err
	}

	// Find the migration to roll back
	var targetMigration *Migration
	for _, migration := range migrations {
		if migration.Version == lastVersion {
			targetMigration = &migration
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %s not found in files", lastVersion)
	}

	if targetMigration.DownSQL == "" {
		return fmt.Errorf("no down migration available for %s", lastVersion)
	}

	fmt.Printf("Rolling back migration %s...\n", lastVersion)

	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute down migration
	if _, err := tx.Exec(targetMigration.DownSQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute down migration %s: %w", lastVersion, err)
	}

	// Remove migration record
	if _, err := tx.Exec(
		"DELETE FROM schema_migrations WHERE version = $1",
		lastVersion,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %s: %w", lastVersion, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %s: %w", lastVersion, err)
	}

	fmt.Printf("Migration %s rolled back successfully\n", lastVersion)
	return nil
}

// Status shows the current migration status
func (m *Migrator) Status() error {
	if err := m.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	migrations, err := m.loadMigrations()
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	if len(migrations) == 0 {
		fmt.Println("No migrations found")
		return nil
	}

	for _, migration := range migrations {
		status := "[ ]"
		if applied[migration.Version] {
			status = "[✓]"
		}
		fmt.Printf("%s %s\n", status, migration.Version)
	}

	fmt.Printf("\nTotal: %d migrations, %d applied\n", len(migrations), len(applied))
	return nil
}

