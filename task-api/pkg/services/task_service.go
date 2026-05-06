package services

import (
	"errors"

	"task-management-api/pkg/models"
	"task-management-api/pkg/repositories"
	"task-management-api/pkg/utils"
)

// TaskService handles task business logic
type TaskService struct {
	taskRepo    *repositories.TaskRepository
	projectRepo *repositories.ProjectRepository
	userRepo    *repositories.UserRepository
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo *repositories.TaskRepository, projectRepo *repositories.ProjectRepository, userRepo *repositories.UserRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// CreateTask creates a new task in a project
func (s *TaskService) CreateTask(projectID, userID string, req *models.CreateTaskRequest) (*models.Task, error) {
	// Check project exists and user has access
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	// Check user is owner or member
	isMember, _ := s.projectRepo.IsMember(projectID, userID)
	if project.OwnerID != userID && !isMember {
		return nil, errors.New("you do not have permission to create tasks in this project")
	}

	// Validate input
	if !utils.ValidateTaskTitle(req.Title) {
		return nil, errors.New("task title must be between 1 and 255 characters")
	}

	if !utils.IsValidTaskStatus(req.Status) {
		return nil, errors.New("invalid task status")
	}

	if !utils.IsValidTaskPriority(req.Priority) {
		return nil, errors.New("invalid task priority")
	}

	// Verify assigned user if provided
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		_, err := s.userRepo.GetByID(*req.AssignedTo)
		if err != nil {
			return nil, errors.New("assigned user not found")
		}
	}

	task := &models.Task{
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		AssignedTo:  req.AssignedTo,
		DueDate:     req.DueDate,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetTask retrieves a task with permission check
func (s *TaskService) GetTask(taskID, userID string) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	// Check user has access to project
	project, err := s.projectRepo.GetByID(task.ProjectID)
	if err != nil {
		return nil, err
	}

	isMember, _ := s.projectRepo.IsMember(task.ProjectID, userID)
	if project.OwnerID != userID && !isMember {
		return nil, errors.New("you do not have permission to access this task")
	}

	return task, nil
}

// GetProjectTasks retrieves all tasks in a project
func (s *TaskService) GetProjectTasks(projectID, userID string, offset, limit int) ([]*models.Task, error) {
	// Check project exists and user has access
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	isMember, _ := s.projectRepo.IsMember(projectID, userID)
	if project.OwnerID != userID && !isMember {
		return nil, errors.New("you do not have permission to access this project")
	}

	return s.taskRepo.GetByProjectID(projectID, offset, limit)
}

// GetProjectTasksWithFilters retrieves tasks with filtering
func (s *TaskService) GetProjectTasksWithFilters(projectID, userID, status, priority string, offset, limit int) ([]*models.Task, error) {
	// Check project exists and user has access
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	isMember, _ := s.projectRepo.IsMember(projectID, userID)
	if project.OwnerID != userID && !isMember {
		return nil, errors.New("you do not have permission to access this project")
	}

	// Validate filters
	if status != "" && !utils.IsValidTaskStatus(status) {
		return nil, errors.New("invalid task status")
	}

	if priority != "" && !utils.IsValidTaskPriority(priority) {
		return nil, errors.New("invalid task priority")
	}

	return s.taskRepo.GetByProjectIDWithFilters(projectID, status, priority, offset, limit)
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(taskID, userID string, req *models.UpdateTaskRequest) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	// Check user has permission to update
	project, err := s.projectRepo.GetByID(task.ProjectID)
	if err != nil {
		return nil, err
	}

	isMember, _ := s.projectRepo.IsMember(task.ProjectID, userID)
	if project.OwnerID != userID && !isMember {
		return nil, errors.New("you do not have permission to update this task")
	}

	// Update fields
	if req.Title != "" {
		if !utils.ValidateTaskTitle(req.Title) {
			return nil, errors.New("task title must be between 1 and 255 characters")
		}
		task.Title = req.Title
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	if req.Status != "" {
		if !utils.IsValidTaskStatus(req.Status) {
			return nil, errors.New("invalid task status")
		}
		task.Status = req.Status
	}

	if req.Priority != "" {
		if !utils.IsValidTaskPriority(req.Priority) {
			return nil, errors.New("invalid task priority")
		}
		task.Priority = req.Priority
	}

	if req.AssignedTo != nil {
		if *req.AssignedTo != "" {
			_, err := s.userRepo.GetByID(*req.AssignedTo)
			if err != nil {
				return nil, errors.New("assigned user not found")
			}
		}
		task.AssignedTo = req.AssignedTo
	}

	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTaskStatus updates only the task status
func (s *TaskService) UpdateTaskStatus(taskID, userID, status string) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	// Check user has permission
	project, err := s.projectRepo.GetByID(task.ProjectID)
	if err != nil {
		return nil, err
	}

	isMember, _ := s.projectRepo.IsMember(task.ProjectID, userID)
	if project.OwnerID != userID && !isMember {
		return nil, errors.New("you do not have permission to update this task")
	}

	// Validate status
	if !utils.IsValidTaskStatus(status) {
		return nil, errors.New("invalid task status")
	}

	if err := s.taskRepo.UpdateStatus(taskID, status); err != nil {
		return nil, err
	}

	task.Status = status
	return task, nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(taskID, userID string) error {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return err
	}

	// Check user has permission to delete
	project, err := s.projectRepo.GetByID(task.ProjectID)
	if err != nil {
		return err
	}

	isMember, _ := s.projectRepo.IsMember(task.ProjectID, userID)
	if project.OwnerID != userID && !isMember {
		return errors.New("you do not have permission to delete this task")
	}

	return s.taskRepo.Delete(taskID)
}

// GetUserTasks retrieves all tasks assigned to a user
func (s *TaskService) GetUserTasks(userID string, offset, limit int) ([]*models.Task, error) {
	return s.taskRepo.GetAssignedToUser(userID, offset, limit)
}
