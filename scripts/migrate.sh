#!/bin/bash

set -e

# Database path
DB_PATH="backend/data/micro_journal.db"
MIGRATIONS_DIR="backend/migrations"

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Change to project root
cd "$PROJECT_ROOT"

# Create data directory if it doesn't exist
mkdir -p "$(dirname "$DB_PATH")"

# Initialize schema_migrations table if it doesn't exist
sqlite3 "$DB_PATH" <<EOF
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
EOF

echo "Checking for pending migrations..."

# Get all migration files sorted
MIGRATIONS=($(ls -1 "$MIGRATIONS_DIR"/*.sql 2>/dev/null | sort))

if [ ${#MIGRATIONS[@]} -eq 0 ]; then
    echo "No migration files found in $MIGRATIONS_DIR"
    exit 0
fi

APPLIED=0

# Run each migration if not already applied
for migration_file in "${MIGRATIONS[@]}"; do
    # Extract migration name (filename without path and extension)
    migration_name=$(basename "$migration_file" .sql)

    # Check if migration has been applied
    already_applied=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM schema_migrations WHERE version = '$migration_name';")

    if [ "$already_applied" -eq 0 ]; then
        echo "Applying migration: $migration_name"

        # Run the migration
        sqlite3 "$DB_PATH" < "$migration_file"

        # Record the migration
        sqlite3 "$DB_PATH" "INSERT INTO schema_migrations (version) VALUES ('$migration_name');"

        APPLIED=$((APPLIED + 1))
    fi
done

if [ $APPLIED -eq 0 ]; then
    echo "No pending migrations. Database is up to date."
else
    echo "Applied $APPLIED migration(s)."
fi

# Show current schema version
LATEST_VERSION=$(sqlite3 "$DB_PATH" "SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1;")
echo "Current schema version: $LATEST_VERSION"
