package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"task-management-api/pkg/config"
	"task-management-api/pkg/controllers"
	"task-management-api/pkg/database"
	"task-management-api/pkg/middleware"
	"task-management-api/pkg/repositories"
	"task-management-api/pkg/services"
)

// @title Task Management API
// @version 1.0
// @description A production-ready REST API for managing tasks and projects
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.example.com/support
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT token with Bearer prefix

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Task Management API on port %s", cfg.App.Port)

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	taskRepo := repositories.NewTaskRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	projectService := services.NewProjectService(projectRepo, userRepo)
	taskService := services.NewTaskService(taskRepo, projectRepo, userRepo)

	// Initialize controllers
	authCtrl := controllers.NewAuthController(authService)
	projectCtrl := controllers.NewProjectController(projectService)
	taskCtrl := controllers.NewTaskController(taskService)

	// Setup router
	router := setupRouter(cfg, authCtrl, projectCtrl, taskCtrl)

	// Start server
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(cfg *config.Config, authCtrl *controllers.AuthController, projectCtrl *controllers.ProjectController, taskCtrl *controllers.TaskController) *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.RequestLoggingMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
			auth.POST("/refresh", authCtrl.Refresh)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Auth protected routes
			auth := protected.Group("/auth")
			{
				auth.GET("/me", authCtrl.GetProfile)
				auth.PUT("/me", authCtrl.UpdateProfile)
			}

			// Project routes
			projects := protected.Group("/projects")
			{
				projects.POST("", projectCtrl.CreateProject)
				projects.GET("", projectCtrl.GetProjects)
				projects.GET("/:id", projectCtrl.GetProject)
				projects.PUT("/:id", projectCtrl.UpdateProject)
				projects.DELETE("/:id", projectCtrl.DeleteProject)
				projects.GET("/:id/members", projectCtrl.GetMembers)
				projects.POST("/:id/members", projectCtrl.AddMember)
				projects.DELETE("/:id/members/:userId", projectCtrl.RemoveMember)

				// Task routes nested under projects
				tasks := projects.Group("/:projectId/tasks")
				{
					tasks.POST("", taskCtrl.CreateTask)
					tasks.GET("", taskCtrl.GetProjectTasks)
				}
			}

			// Task routes (flat)
			tasks := protected.Group("/tasks")
			{
				tasks.GET("/:taskId", taskCtrl.GetTask)
				tasks.PUT("/:taskId", taskCtrl.UpdateTask)
				tasks.PATCH("/:taskId/status", taskCtrl.UpdateTaskStatus)
				tasks.DELETE("/:taskId", taskCtrl.DeleteTask)
			}
		}
	}

	return router
}
