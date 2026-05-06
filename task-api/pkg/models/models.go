package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password
	Name      string    `json:"name"`
	Role      string    `json:"role"` // admin, user, viewer
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Project represents a project
type Project struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // active, archived
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Task represents a task within a project
type Task struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	AssignedTo  *string   `json:"assigned_to"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // todo, in_progress, done
	Priority    string    `json:"priority"` // low, medium, high
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectMember represents membership in a project
type ProjectMember struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"` // owner, editor, viewer
	CreatedAt time.Time `json:"created_at"`
}

// Request/Response DTOs

// RegisterRequest is the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest is the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse contains the JWT token after successful login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// CreateProjectRequest is the request body for creating a project
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateProjectRequest is the request body for updating a project
type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// CreateTaskRequest is the request body for creating a task
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Status      string     `json:"status" binding:"required"`
	Priority    string     `json:"priority" binding:"required"`
	AssignedTo  *string    `json:"assigned_to"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskRequest is the request body for updating a task
type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssignedTo  *string    `json:"assigned_to"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskStatusRequest is the request body for updating task status only
type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// AddProjectMemberRequest is the request body for adding a project member
type AddProjectMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

// UpdateUserRequest is the request body for updating user profile
type UpdateUserRequest struct {
	Name string `json:"name"`
	Email string `json:"email" binding:"email"`
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Type  string `json:"type"` // access or refresh
}

// ErrorResponse is a standard error response format
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse is a standard success response format
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}
