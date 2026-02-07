# Database Migrations

## Overview

This directory contains SQL migration files for the VocabWeb database schema.

## Migration Files

- `001_initial_schema.up.sql` - Creates all 12 core tables
- `001_initial_schema.down.sql` - Drops all tables (rollback)

## Database Schema

### Core Tables (12)

1. **profiles** - User extended profile information
2. **words** - Global word dictionary (shared)
3. **groups** - User-defined word groups/folders
4. **user_words** - User collected words with context
5. **tags** - User-defined tags
6. **user_word_tags** - Many-to-many word-tag relationship
7. **review_logs** - Spaced repetition history (SM-2 algorithm)
8. **daily_stats** - Daily learning statistics
9. **achievements** - Achievement/badge definitions
10. **user_achievements** - User earned achievements
11. **exam_wordlists** - Exam-specific word lists
12. **study_plans** - User exam preparation plans

### Key Features

- **SM-2 Algorithm Support**: `review_logs` table includes:
  - `easiness_factor` (DECIMAL 4,2, default 2.5)
  - `interval` (INTEGER, days until next review)
  - `repetitions` (INTEGER, consecutive correct reviews)
  - `next_review_at` (TIMESTAMPTZ)
  - `quality` (INTEGER 0-5, user rating)

- **Proper Indexing**: All foreign keys and frequently queried columns are indexed
- **Constraints**: CHECK constraints for data validation
- **Timestamps**: All tables have `created_at` and `updated_at`
- **Foreign Keys**: Proper CASCADE/SET NULL behavior

## Usage

### Using the Go Migrator

```go
package main

import (
    "database/sql"
    "log"
    
    "your-module/internal/database"
    _ "github.com/lib/pq"
)

func main() {
    // Connect to database
    db, err := sql.Open("postgres", "postgres://user:pass@localhost/vocabweb?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Create migrator
    migrator := database.NewMigrator(db, "./migrations")
    
    // Run migrations
    if err := migrator.Up(); err != nil {
        log.Fatal(err)
    }
}
```

### CLI Commands

```bash
# Apply all pending migrations
go run cmd/migrate/main.go up

# Rollback last migration
go run cmd/migrate/main.go down

# Show migration status
go run cmd/migrate/main.go status
```

### Direct SQL Execution

```bash
# Apply migrations
psql -U postgres -d vocabweb -f migrations/001_initial_schema.up.sql

# Rollback
psql -U postgres -d vocabweb -f migrations/001_initial_schema.down.sql
```

## Migration Tracking

The migrator automatically creates a `schema_migrations` table to track applied migrations:

```sql
CREATE TABLE schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## Adding New Migrations

1. Create new files with incremented version number:
   - `002_add_feature.up.sql`
   - `002_add_feature.down.sql`

2. Follow naming convention: `{version}_{description}.{up|down}.sql`

3. Always provide both up and down migrations

4. Test rollback before committing

## Notes

- PostgreSQL 15+ required
- Uses UUID extension for primary keys
- All timestamps use TIMESTAMPTZ (timezone-aware)
- Foreign keys use CASCADE or SET NULL appropriately
- JSONB used for flexible schema (word definitions)
