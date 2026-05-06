package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"task-management-api/pkg/models"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project in the database
func (r *ProjectRepository) Create(project *models.Project) error {
	query := `
		INSERT INTO projects (owner_id, name, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, project.OwnerID, project.Name, project.Description, project.Status).
		Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating project: %w", err)
	}

	return nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id string) (*models.Project, error) {
	query := `
		SELECT id, owner_id, name, description, status, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	project := &models.Project{}
	err := r.db.QueryRow(query, id).
		Scan(&project.ID, &project.OwnerID, &project.Name, &project.Description, &project.Status, &project.CreatedAt, &project.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving project: %w", err)
	}

	return project, nil
}

// GetByOwnerID retrieves all projects owned by a user
func (r *ProjectRepository) GetByOwnerID(ownerID string, offset, limit int) ([]*models.Project, error) {
	query := `
		SELECT id, owner_id, name, description, status, created_at, updated_at
		FROM projects
		WHERE owner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, ownerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(&project.ID, &project.OwnerID, &project.Name, &project.Description, &project.Status, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning project: %w", err)
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

// GetUserProjects retrieves all projects accessible to a user (owned or member)
func (r *ProjectRepository) GetUserProjects(userID string, offset, limit int) ([]*models.Project, error) {
	query := `
		SELECT DISTINCT p.id, p.owner_id, p.name, p.description, p.status, p.created_at, p.updated_at
		FROM projects p
		LEFT JOIN project_members pm ON p.id = pm.project_id
		WHERE p.owner_id = $1 OR pm.user_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(&project.ID, &project.OwnerID, &project.Name, &project.Description, &project.Status, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning project: %w", err)
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

// Update updates a project
func (r *ProjectRepository) Update(project *models.Project) error {
	query := `
		UPDATE projects
		SET name = $1, description = $2, status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, project.Name, project.Description, project.Status, project.ID).
		Scan(&project.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("project not found")
	}
	if err != nil {
		return fmt.Errorf("error updating project: %w", err)
	}

	return nil
}

// Delete deletes a project
func (r *ProjectRepository) Delete(id string) error {
	query := `DELETE FROM projects WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deleted rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}

// AddMember adds a user to a project
func (r *ProjectRepository) AddMember(projectID, userID, role string) error {
	query := `
		INSERT INTO project_members (project_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, user_id) DO UPDATE SET role = $3
	`

	_, err := r.db.Exec(query, projectID, userID, role)
	if err != nil {
		return fmt.Errorf("error adding member: %w", err)
	}

	return nil
}

// RemoveMember removes a user from a project
func (r *ProjectRepository) RemoveMember(projectID, userID string) error {
	query := `DELETE FROM project_members WHERE project_id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, projectID, userID)
	if err != nil {
		return fmt.Errorf("error removing member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deleted rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("member not found")
	}

	return nil
}

// GetMembers retrieves all members of a project
func (r *ProjectRepository) GetMembers(projectID string) ([]*models.ProjectMember, error) {
	query := `
		SELECT id, project_id, user_id, role, created_at
		FROM project_members
		WHERE project_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving members: %w", err)
	}
	defer rows.Close()

	var members []*models.ProjectMember
	for rows.Next() {
		member := &models.ProjectMember{}
		err := rows.Scan(&member.ID, &member.ProjectID, &member.UserID, &member.Role, &member.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning member: %w", err)
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating members: %w", err)
	}

	return members, nil
}

// GetMemberRole retrieves a member's role in a project
func (r *ProjectRepository) GetMemberRole(projectID, userID string) (string, error) {
	query := `SELECT role FROM project_members WHERE project_id = $1 AND user_id = $2`

	var role string
	err := r.db.QueryRow(query, projectID, userID).Scan(&role)

	if err == sql.ErrNoRows {
		return "", errors.New("member not found")
	}
	if err != nil {
		return "", fmt.Errorf("error retrieving member role: %w", err)
	}

	return role, nil
}

// IsMember checks if a user is a member of a project
func (r *ProjectRepository) IsMember(projectID, userID string) (bool, error) {
	query := `SELECT 1 FROM project_members WHERE project_id = $1 AND user_id = $2`

	var exists bool
	err := r.db.QueryRow(query, projectID, userID).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking member: %w", err)
	}

	return true, nil
}
