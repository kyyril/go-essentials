package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"task-management-api/pkg/middleware"
	"task-management-api/pkg/models"
	"task-management-api/pkg/services"
)

// TaskController handles task endpoints
type TaskController struct {
	taskService *services.TaskService
}

// NewTaskController creates a new task controller
func NewTaskController(taskService *services.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task in a project
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param projectId path string true "Project ID"
// @Param request body models.CreateTaskRequest true "Create task request"
// @Success 201 {object} models.SuccessResponse{data=models.Task}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/v1/projects/{projectId}/tasks [post]
func (ctrl *TaskController) CreateTask(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projectID := c.Param("projectId")

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	task, err := ctrl.taskService.CreateTask(projectID, userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "project not found" {
			status = http.StatusNotFound
			errorType = "not_found"
		} else {
			status = http.StatusBadRequest
			errorType = "create_error"
		}

		c.JSON(status, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Data:    task,
		Message: "Task created successfully",
		Code:    http.StatusCreated,
	})
}

// GetProjectTasks godoc
// @Summary Get project tasks
// @Description Get all tasks in a project with optional filtering
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param projectId path string true "Project ID"
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Param status query string false "Filter by status (todo, in_progress, done)"
// @Param priority query string false "Filter by priority (low, medium, high)"
// @Success 200 {object} models.SuccessResponse{data=[]models.Task}
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/projects/{projectId}/tasks [get]
func (ctrl *TaskController) GetProjectTasks(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projectID := c.Param("projectId")

	offset := 0
	limit := 10
	status := c.Query("status")
	priority := c.Query("priority")

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

	var tasks []*models.Task
	var err error

	if status != "" || priority != "" {
		tasks, err = ctrl.taskService.GetProjectTasksWithFilters(projectID, userID, status, priority, offset, limit)
	} else {
		tasks, err = ctrl.taskService.GetProjectTasks(projectID, userID, offset, limit)
	}

	if err != nil {
		status_code := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "project not found" {
			status_code = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "you do not have permission to access this project" {
			status_code = http.StatusForbidden
			errorType = "forbidden"
		}

		c.JSON(status_code, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status_code,
		})
		return
	}

	if tasks == nil {
		tasks = make([]*models.Task, 0)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    tasks,
		Message: "Tasks retrieved successfully",
		Code:    http.StatusOK,
	})
}

// GetTask godoc
// @Summary Get a specific task
// @Description Get details of a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param taskId path string true "Task ID"
// @Success 200 {object} models.SuccessResponse{data=models.Task}
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/tasks/{taskId} [get]
func (ctrl *TaskController) GetTask(c *gin.Context) {
	userID := middleware.GetUserID(c)
	taskID := c.Param("taskId")

	task, err := ctrl.taskService.GetTask(taskID, userID)
	if err != nil {
		status_code := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "task not found" {
			status_code = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "you do not have permission to access this task" {
			status_code = http.StatusForbidden
			errorType = "forbidden"
		}

		c.JSON(status_code, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status_code,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    task,
		Message: "Task retrieved successfully",
		Code:    http.StatusOK,
	})
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update task details
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param taskId path string true "Task ID"
// @Param request body models.UpdateTaskRequest true "Update request"
// @Success 200 {object} models.SuccessResponse{data=models.Task}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/tasks/{taskId} [put]
func (ctrl *TaskController) UpdateTask(c *gin.Context) {
	userID := middleware.GetUserID(c)
	taskID := c.Param("taskId")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	task, err := ctrl.taskService.UpdateTask(taskID, userID, &req)
	if err != nil {
		status_code := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "task not found" {
			status_code = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "you do not have permission to update this task" {
			status_code = http.StatusForbidden
			errorType = "forbidden"
		} else {
			status_code = http.StatusBadRequest
			errorType = "update_error"
		}

		c.JSON(status_code, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status_code,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    task,
		Message: "Task updated successfully",
		Code:    http.StatusOK,
	})
}

// UpdateTaskStatus godoc
// @Summary Update task status
// @Description Update only the status of a task
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param taskId path string true "Task ID"
// @Param request body models.UpdateTaskStatusRequest true "Update status request"
// @Success 200 {object} models.SuccessResponse{data=models.Task}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/tasks/{taskId}/status [patch]
func (ctrl *TaskController) UpdateTaskStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	taskID := c.Param("taskId")

	var req models.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	task, err := ctrl.taskService.UpdateTaskStatus(taskID, userID, req.Status)
	if err != nil {
		status_code := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "task not found" {
			status_code = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "you do not have permission to update this task" {
			status_code = http.StatusForbidden
			errorType = "forbidden"
		} else {
			status_code = http.StatusBadRequest
			errorType = "update_error"
		}

		c.JSON(status_code, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status_code,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    task,
		Message: "Task status updated successfully",
		Code:    http.StatusOK,
	})
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task from a project
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param taskId path string true "Task ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/tasks/{taskId} [delete]
func (ctrl *TaskController) DeleteTask(c *gin.Context) {
	userID := middleware.GetUserID(c)
	taskID := c.Param("taskId")

	err := ctrl.taskService.DeleteTask(taskID, userID)
	if err != nil {
		status_code := http.StatusInternalServerError
		errorType := "error"
		if err.Error() == "task not found" {
			status_code = http.StatusNotFound
			errorType = "not_found"
		} else if err.Error() == "you do not have permission to delete this task" {
			status_code = http.StatusForbidden
			errorType = "forbidden"
		}

		c.JSON(status_code, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
			Code:    status_code,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    nil,
		Message: "Task deleted successfully",
		Code:    http.StatusOK,
	})
}
