# Database Migrations with Atlas

This project uses [Atlas](https://atlasgo.io) for database schema migrations. Atlas provides a modern approach to database migrations with features like automatic migration generation, schema diffing, and migration linting.

## Why Atlas?

- **Automatic Migration Generation**: Generate migrations from your Go structs
- **Schema Diffing**: Compare schema with database state
- **Migration Linting**: Detect dangerous operations before applying
- **Rollback Support**: Safely rollback migrations
- **Declarative Migrations**: Define schema as code
- **Multi-Environment**: Separate configs for dev, staging, prod

## Installation

### Using Makefile (Recommended)

```bash
make install-tools
```

### Manual Installation

**Linux/macOS:**
```bash
go install ariga.io/atlas/cmd/atlas@latest
```

**Windows:**
```bash
go install ariga.io/atlas/cmd/atlas@latest
```

Or download from [Atlas Releases](https://github.com/ariga/atlas/releases)

## Quick Start

### 1. Generate Initial Migration

After setting up your database, generate the initial migration from your GORM models:

```bash
make atlas-generate
```

This creates migration files in the `migrations/` directory based on your domain entities.

### 2. Review Migration

Check the generated SQL in `migrations/`:

```bash
ls -la migrations/
cat migrations/20240204_initial.sql
```

### 3. Apply Migration

Apply the migration to your database:

```bash
make atlas-apply
```

### 4. Check Status

Verify migration status:

```bash
make atlas-status
```

## Configuration

Atlas configuration is in `atlas.hcl` with three environments:

### Local Environment (Default)
```bash
make atlas-status                    # Uses local env
ATLAS_ENV=local make atlas-status    # Explicit
```

### Development Environment
```bash
ATLAS_ENV=dev make atlas-apply
```

### Production Environment
```bash
ATLAS_ENV=prod make atlas-apply
```

Production environment includes:
- Migration history tracking
- Destructive change detection
- Additional safety checks

## Common Workflows

### Making Schema Changes

1. **Modify your domain entities** (e.g., `internal/domain/user.go`):
   ```go
   type User struct {
       ID        uint           `gorm:"primarykey" json:"id"`
       Email     string         `gorm:"uniqueIndex;not null" json:"email"`
       Password  string         `gorm:"not null" json:"-"`
       Name      string         `gorm:"not null" json:"name"`
       Role      string         `gorm:"default:'user'" json:"role"` // NEW FIELD
       CreatedAt time.Time      `json:"created_at"`
       UpdatedAt time.Time      `json:"updated_at"`
       DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
   }
   ```

2. **Generate migration**:
   ```bash
   make atlas-generate
   ```

3. **Review the generated SQL**:
   ```bash
   cat migrations/20240204_add_role_to_users.sql
   ```

4. **Apply migration**:
   ```bash
   make atlas-apply
   ```

5. **Verify**:
   ```bash
   make atlas-status
   ```

### Creating Manual Migrations

For custom SQL operations:

```bash
make atlas-new
# Enter migration name: add_indexes

# Edit migrations/20240204_add_indexes.sql
```

Example migration file:
```sql
-- Add indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

### Inspecting Database Schema

View current database schema:

```bash
make atlas-inspect
```

### Comparing Schema vs Database

See differences between your code and database:

```bash
make atlas-diff
```

### Validating Migrations

Check migration files for issues:

```bash
make atlas-validate
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make atlas-help` | Show all Atlas commands |
| `make atlas-generate` | Generate migration from schema changes |
| `make atlas-apply` | Apply pending migrations |
| `make atlas-status` | Check migration status |
| `make atlas-validate` | Validate migration files |
| `make atlas-inspect` | Inspect database schema |
| `make atlas-diff` | Show schema differences |
| `make atlas-hash` | Rehash migration directory |
| `make atlas-new` | Create new empty migration |
| `make atlas-clean` | Clean database (dev only) |

## Migration Files

Migration files are stored in `migrations/` directory:

```
migrations/
├── 20240204000001_initial.sql           # Initial schema
├── 20240204000002_add_role.sql          # Add role field
└── atlas.sum                             # Migration checksum
```

### Migration File Format

```sql
-- Create "users" table
CREATE TABLE "users" (
  "id" serial NOT NULL,
  "email" character varying NOT NULL,
  "password" character varying NOT NULL,
  "name" character varying NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
```

## Best Practices

### 1. Always Review Generated Migrations

Before applying:
```bash
make atlas-generate
cat migrations/latest_migration.sql  # Review changes
make atlas-apply
```

### 2. Use Transactions

Atlas automatically wraps migrations in transactions, but be aware of operations that can't be rolled back.

### 3. Test Migrations in Development First

```bash
# Development
ATLAS_ENV=dev make atlas-apply

# Production (after testing)
ATLAS_ENV=prod make atlas-apply
```

### 4. Backup Before Production Migrations

```bash
# Backup database
pg_dump -U postgres -d auth_service > backup.sql

# Apply migration
ATLAS_ENV=prod make atlas-apply

# Restore if needed
psql -U postgres -d auth_service < backup.sql
```

### 5. Never Modify Applied Migrations

Once a migration is applied, never modify it. Create a new migration instead:

```bash
# ❌ Don't do this
vim migrations/20240204_initial.sql  # Modify applied migration

# ✅ Do this instead
make atlas-new                       # Create new migration
```

### 6. Keep Migrations in Version Control

```bash
git add migrations/
git commit -m "Add user role migration"
```

## Migration Status

Check which migrations are applied:

```bash
$ make atlas-status

Migration Status:
┌───────────────────────┬──────────────────┬────────────┐
│ Version               │ Description      │ Status     │
├───────────────────────┼──────────────────┼────────────┤
│ 20240204000001        │ initial          │ Applied    │
│ 20240204000002        │ add_role         │ Applied    │
│ 20240204000003        │ add_indexes      │ Pending    │
└───────────────────────┴──────────────────┴────────────┘
```

## Rollback

Atlas supports rollback through versioned migrations:

```bash
# Apply specific version
atlas migrate apply --env local --version 20240204000001

# Or create down migration manually
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Database Migration

on:
  push:
    branches: [main]

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Atlas
        run: |
          curl -sSf https://atlasgo.sh | sh
      
      - name: Validate Migrations
        run: atlas migrate validate --env prod
      
      - name: Apply Migrations
        run: atlas migrate apply --env prod
        env:
          DB_HOST: ${{ secrets.DB_HOST }}
          DB_USER: ${{ secrets.DB_USER }}
          DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
```

## Troubleshooting

### Migration Checksum Mismatch

If you see checksum errors:

```bash
# Rehash migrations
make atlas-hash
```

### Migration Already Applied

If Atlas thinks a migration is already applied:

```bash
# Check status
make atlas-status

# Inspect actual database
make atlas-inspect
```

### Clean Development Database

To start fresh in development:

```bash
make atlas-clean  # Drops all tables
make atlas-apply  # Reapply migrations
```

### Connection Issues

Check your `.env` file:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_service
```

## Migration vs GORM AutoMigrate

### GORM AutoMigrate (Development)
- ✅ Easy to use
- ✅ No migration files needed
- ❌ No version control
- ❌ No rollback
- ❌ Limited control

### Atlas Migrations (Production)
- ✅ Version controlled
- ✅ Rollback support
- ✅ Migration linting
- ✅ Full SQL control
- ✅ Audit trail
- ✅ Team collaboration

## Transition from GORM AutoMigrate

If you're currently using GORM AutoMigrate:

1. **Generate initial schema snapshot**:
   ```bash
   make atlas-generate
   ```

2. **Review generated migration**:
   ```bash
   cat migrations/20240204_initial.sql
   ```

3. **Apply to create tracking table**:
   ```bash
   make atlas-apply
   ```

4. **Future changes use Atlas**:
   ```bash
   # Modify domain/user.go
   make atlas-generate
   make atlas-apply
   ```

The application detects Atlas migrations automatically and skips GORM AutoMigrate.

## Advanced Features

### Schema Diffing

Compare two database states:

```bash
atlas schema diff \
  --from "postgres://localhost:5432/db1" \
  --to "postgres://localhost:5432/db2"
```

### Migration Linting

Detect dangerous operations:

```bash
atlas migrate lint --env prod
```

This detects:
- Dropping tables/columns
- Changing column types
- Adding non-nullable columns without defaults

### Declarative Migrations

Instead of SQL, you can use HCL:

```hcl
table "users" {
  schema = schema.public
  column "id" {
    type = serial
  }
  column "email" {
    type = varchar(255)
    null = false
  }
}
```

## Resources

- **Atlas Documentation**: https://atlasgo.io/docs
- **Atlas CLI Reference**: https://atlasgo.io/cli-reference
- **Migration Guides**: https://atlasgo.io/guides
- **Discord Community**: https://discord.gg/zZ6sWVg6NT

## Summary

Atlas provides professional-grade database migrations with:
- Automatic generation from Go structs
- Version control integration
- Rollback support
- Migration validation
- Production safety features

Use `make atlas-help` to see all available commands!
