package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"task-management-api/pkg/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password, name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, user.Email, user.Password, user.Name, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return fmt.Errorf("email already exists")
		}
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, user.Name, user.Email, user.ID).
		Scan(&user.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("user not found")
	}
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deleted rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetAll retrieves all users with pagination
func (r *UserRepository) GetAll(offset, limit int) ([]*models.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// ExistsByEmail checks if a user exists with given email
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	query := `SELECT 1 FROM users WHERE email = $1`

	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %w", err)
	}

	return true, nil
}
