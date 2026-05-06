package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"task-management-api/pkg/models"
)

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create creates a new task in the database
func (r *TaskRepository) Create(task *models.Task) error {
	query := `
		INSERT INTO tasks (project_id, assigned_to, title, description, status, priority, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, task.ProjectID, task.AssignedTo, task.Title, task.Description, task.Status, task.Priority, task.DueDate).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating task: %w", err)
	}

	return nil
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(id string) (*models.Task, error) {
	query := `
		SELECT id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	task := &models.Task{}
	err := r.db.QueryRow(query, id).
		Scan(&task.ID, &task.ProjectID, &task.AssignedTo, &task.Title, &task.Description, &task.Status, &task.Priority, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving task: %w", err)
	}

	return task, nil
}

// GetByProjectID retrieves all tasks in a project
func (r *TaskRepository) GetByProjectID(projectID string, offset, limit int) ([]*models.Task, error) {
	query := `
		SELECT id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, projectID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.ID, &task.ProjectID, &task.AssignedTo, &task.Title, &task.Description, &task.Status, &task.Priority, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// GetByProjectIDWithFilters retrieves tasks in a project with filtering
func (r *TaskRepository) GetByProjectIDWithFilters(projectID, status, priority string, offset, limit int) ([]*models.Task, error) {
	query := `
		SELECT id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE project_id = $1
	`

	args := []interface{}{projectID}
	argCount := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if priority != "" {
		query += fmt.Sprintf(" AND priority = $%d", argCount)
		args = append(args, priority)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.ID, &task.ProjectID, &task.AssignedTo, &task.Title, &task.Description, &task.Status, &task.Priority, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// Update updates a task
func (r *TaskRepository) Update(task *models.Task) error {
	query := `
		UPDATE tasks
		SET project_id = $1, assigned_to = $2, title = $3, description = $4, status = $5, priority = $6, due_date = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, task.ProjectID, task.AssignedTo, task.Title, task.Description, task.Status, task.Priority, task.DueDate, task.ID).
		Scan(&task.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("task not found")
	}
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	return nil
}

// UpdateStatus updates only the status of a task
func (r *TaskRepository) UpdateStatus(id, status string) error {
	query := `
		UPDATE tasks
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("error updating task status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking updated rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// Delete deletes a task
func (r *TaskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deleted rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// GetAssignedToUser retrieves all tasks assigned to a user
func (r *TaskRepository) GetAssignedToUser(userID string, offset, limit int) ([]*models.Task, error) {
	query := `
		SELECT id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE assigned_to = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.ID, &task.ProjectID, &task.AssignedTo, &task.Title, &task.Description, &task.Status, &task.Priority, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}
