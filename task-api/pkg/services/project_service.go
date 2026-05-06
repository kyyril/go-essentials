package services

import (
	"errors"
	"fmt"

	"task-management-api/pkg/models"
	"task-management-api/pkg/repositories"
	"task-management-api/pkg/utils"
)

// ProjectService handles project business logic
type ProjectService struct {
	projectRepo *repositories.ProjectRepository
	userRepo    *repositories.UserRepository
}

// NewProjectService creates a new project service
func NewProjectService(projectRepo *repositories.ProjectRepository, userRepo *repositories.UserRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(userID string, req *models.CreateProjectRequest) (*models.Project, error) {
	// Validate input
	if !utils.ValidateProjectName(req.Name) {
		return nil, errors.New("project name must be between 1 and 255 characters")
	}

	project := &models.Project{
		OwnerID:     userID,
		Name:        req.Name,
		Description: req.Description,
		Status:      "active",
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetProject retrieves a project by ID with permission check
func (s *ProjectService) GetProject(projectID, userID string) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	// Check if user has access
	if !s.canAccessProject(project, userID) {
		return nil, errors.New("you do not have permission to access this project")
	}

	return project, nil
}

// GetUserProjects retrieves all projects accessible to a user
func (s *ProjectService) GetUserProjects(userID string, offset, limit int) ([]*models.Project, error) {
	return s.projectRepo.GetUserProjects(userID, offset, limit)
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(projectID, userID string, req *models.UpdateProjectRequest) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	// Check if user is project owner
	if project.OwnerID != userID {
		return nil, errors.New("only project owner can update the project")
	}

	if req.Name != "" {
		if !utils.ValidateProjectName(req.Name) {
			return nil, errors.New("project name must be between 1 and 255 characters")
		}
		project.Name = req.Name
	}

	if req.Description != "" {
		project.Description = req.Description
	}

	if req.Status != "" {
		if !utils.IsValidProjectStatus(req.Status) {
			return nil, errors.New("invalid project status")
		}
		project.Status = req.Status
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(projectID, userID string) error {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}

	// Check if user is project owner
	if project.OwnerID != userID {
		return errors.New("only project owner can delete the project")
	}

	return s.projectRepo.Delete(projectID)
}

// AddProjectMember adds a user to a project
func (s *ProjectService) AddProjectMember(projectID, userID, memberUserID, role string) error {
	// Check if project exists and user has permission
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}

	// Check if user is project owner
	if project.OwnerID != userID {
		return errors.New("only project owner can add members")
	}

	// Validate role
	if !utils.IsValidProjectMemberRole(role) {
		return errors.New("invalid member role")
	}

	// Check if member user exists
	_, err = s.userRepo.GetByID(memberUserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Add member
	return s.projectRepo.AddMember(projectID, memberUserID, role)
}

// RemoveProjectMember removes a user from a project
func (s *ProjectService) RemoveProjectMember(projectID, userID, memberUserID string) error {
	// Check if project exists and user has permission
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}

	// Check if user is project owner
	if project.OwnerID != userID {
		return errors.New("only project owner can remove members")
	}

	// Cannot remove owner
	if project.OwnerID == memberUserID {
		return errors.New("cannot remove project owner")
	}

	return s.projectRepo.RemoveMember(projectID, memberUserID)
}

// GetProjectMembers retrieves all members of a project
func (s *ProjectService) GetProjectMembers(projectID, userID string) ([]*models.ProjectMember, error) {
	// Check if user has access to project
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	if !s.canAccessProject(project, userID) {
		return nil, errors.New("you do not have permission to access this project")
	}

	return s.projectRepo.GetMembers(projectID)
}

// canAccessProject checks if a user can access a project
func (s *ProjectService) canAccessProject(project *models.Project, userID string) bool {
	// Owner has access
	if project.OwnerID == userID {
		return true
	}

	// Check if user is a member
	isMember, err := s.projectRepo.IsMember(project.ID, userID)
	if err == nil && isMember {
		return true
	}

	return false
}

// CheckProjectAccess checks if user has permission for a project
func (s *ProjectService) CheckProjectAccess(projectID, userID string) error {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}

	if !s.canAccessProject(project, userID) {
		return errors.New("you do not have permission to access this project")
	}

	return nil
}

// GetUserRole gets the role of a user in a project (returns "" if not a member and not owner)
func (s *ProjectService) GetUserRole(projectID, userID string) (string, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return "", err
	}

	if project.OwnerID == userID {
		return "owner", nil
	}

	role, err := s.projectRepo.GetMemberRole(projectID, userID)
	if err != nil {
		return "", nil // Not a member
	}

	return role, nil
}
