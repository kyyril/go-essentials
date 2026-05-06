package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"task-management-api/pkg/middleware"
	"task-management-api/pkg/models"
	"task-management-api/pkg/services"
)

// ProjectController handles project endpoints
type ProjectController struct {
	projectService *services.ProjectService
}

// NewProjectController creates a new project controller
func NewProjectController(projectService *services.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project owned by the current user
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.CreateProjectRequest true "Create project request"
// @Success 201 {object} models.SuccessResponse{data=models.Project}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /projects [post]
func (ctrl *ProjectController) CreateProject(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	project, err := ctrl.projectService.CreateProject(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "create_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Data:    project,
		Message: "Project created successfully",
		Code:    http.StatusCreated,
	})
}

// GetProjects godoc
// @Summary Get user's projects
// @Description Get all projects accessible to the current user
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} models.SuccessResponse{data=[]models.Project}
// @Failure 401 {object} models.ErrorResponse
// @Router /projects [get]
func (ctrl *ProjectController) GetProjects(c *gin.Context) {
	userID := middleware.GetUserID(c)

	offset := 0
	limit := 10

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	projects, err := ctrl.projectService.GetUserProjects(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if projects == nil {
		projects = make([]*models.Project, 0)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    projects,
		Message: "Projects retrieved successfully",
		Code:    http.StatusOK,
	})
}

// GetProject godoc
// @Summary Get a specific project
// @Description Get details of a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Success 200 {object} models.SuccessResponse{data=models.Project}
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /projects/{id} [get]
func (ctrl *ProjectController) GetProject(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projectID := c.Param("id")

	project, err := ctrl.projectService.GetProject(projectID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "project not found" {
			status = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "you do not have permission to access this project" {
			status = http.StatusForbidden
			errorType = "forbidden"
		}

		c.JSON(status, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    project,
		Message: "Project retrieved successfully",
		Code:    http.StatusOK,
	})
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update project details
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Param request body models.UpdateProjectRequest true "Update request"
// @Success 200 {object} models.SuccessResponse{data=models.Project}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /projects/{id} [put]
func (ctrl *ProjectController) UpdateProject(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projectID := c.Param("id")

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	project, err := ctrl.projectService.UpdateProject(projectID, userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "project not found" {
			status = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "only project owner can update the project" {
			status = http.StatusForbidden
			errorType = "forbidden"
		} else {
			status = http.StatusBadRequest
			errorType = "update_error"
		}

		c.JSON(status, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    project,
		Message: "Project updated successfully",
		Code:    http.StatusOK,
	})
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project (owner only)
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /projects/{id} [delete]
func (ctrl *ProjectController) DeleteProject(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projectID := c.Param("id")

	err := ctrl.projectService.DeleteProject(projectID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "project not found" {
			status = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "only project owner can delete the project" {
			status = http.StatusForbidden
			errorType = "forbidden"
		}

		c.JSON(status, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    nil,
		Message: "Project deleted successfully",
		Code:    http.StatusOK,
	})
}

// GetMembers godoc
// @Summary Get project members
// @Description Get all members of a project
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Success 200 {object} models.SuccessResponse{data=[]models.ProjectMember}
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /projects/{id}/members [get]
func (ctrl *ProjectController) GetMembers(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projectID := c.Param("id")

	members, err := ctrl.projectService.GetProjectMembers(projectID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "project not found" {
			status = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "you do not have permission to access this project" {
			status = http.StatusForbidden
			errorType = "forbidden"
		}

		c.JSON(status, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	if members == nil {
		members = make([]*models.ProjectMember, 0)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    members,
		Message: "Members retrieved successfully",
		Code:    http.StatusOK,
	})
}

// AddMember godoc
// @Summary Add member to project
// @Description Add a user as a member to a project
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Param request body models.AddProjectMemberRequest true "Add member request"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /projects/{id}/members [post]
func (ctrl *ProjectController) AddMember(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projectID := c.Param("id")

	var req models.AddProjectMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := ctrl.projectService.AddProjectMember(projectID, userID, req.UserID, req.Role)
	if err != nil {
		status := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "project not found" {
			status = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "only project owner can add members" {
			status = http.StatusForbidden
			errorType = "forbidden"
		} else {
			status = http.StatusBadRequest
			errorType = "add_member_error"
		}

		c.JSON(status, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Data:    nil,
		Message: "Member added successfully",
		Code:    http.StatusCreated,
	})
}

// RemoveMember godoc
// @Summary Remove member from project
// @Description Remove a user from a project
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Param userId path string true "User ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /projects/{id}/members/{userId} [delete]
func (ctrl *ProjectController) RemoveMember(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projectID := c.Param("id")
	memberUserID := c.Param("userId")

	err := ctrl.projectService.RemoveProjectMember(projectID, userID, memberUserID)
	if err != nil {
		status := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "project not found" {
			status = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "only project owner can remove members" {
			status = http.StatusForbidden
			errorType = "forbidden"
		} else {
			status = http.StatusBadRequest
			errorType = "remove_error"
		}

		c.JSON(status, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    nil,
		Message: "Member removed successfully",
		Code:    http.StatusOK,
	})
}
