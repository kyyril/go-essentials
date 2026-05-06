package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"task-management-api/pkg/config"
)

// InitDB initializes the database connection pool
func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// RunMigrations executes all database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsersTable,
		createProjectsTable,
		createTasksTable,
		createProjectMembersTable,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("error running migration %d: %w", i+1, err)
		}
		log.Printf("Migration %d completed successfully\n", i+1)
	}

	return nil
}

const (
	createUsersTable = `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	createProjectsTable = `
		CREATE TABLE IF NOT EXISTS projects (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_projects_owner_id ON projects(owner_id);
		CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
	`

	createTasksTable = `
		CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'todo',
			priority VARCHAR(50) NOT NULL DEFAULT 'medium',
			due_date TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks(assigned_to);
		CREATE INDEX IF NOT EXISTS idx_tasks_title ON tasks(title);
	`

	createProjectMembersTable = `
		CREATE TABLE IF NOT EXISTS project_members (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL DEFAULT 'editor',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(project_id, user_id)
		);
		CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id);
		CREATE INDEX IF NOT EXISTS idx_project_members_user_id ON project_members(user_id);
	`
)
