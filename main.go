package main

import (
	"TaskManagementService/pkg/database"
	"TaskManagementService/pkg/handlers"
	"embed"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	// Create Gin router
	router := gin.Default()

	// database migrations
	db, err := database.NewPostgres("myuser", "securepass", "localhost", "5432", "task_manager")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations applied!")

	// Instantiate tasks Handler and provide a data store
	//store := tasks.NewMemStore()
	//tasksHandler := handlers.NewTasksHandler(store)
	tasksHandler := handlers.NewTasksDatabaseHandler(db)

	// Register Routes
	router.GET("/", homePage)
	router.GET("/tasks", tasksHandler.ListTasks)
	router.GET("/tasks/filter", tasksHandler.FilterTasks)
	router.POST("/tasks", tasksHandler.CreateTask)
	router.GET("/tasks/:id", tasksHandler.GetTask)
	router.PUT("/tasks/:id", tasksHandler.UpdateTask)
	router.DELETE("/tasks/:id", tasksHandler.DeleteTask)

	// Start the server
	router.Run(":8080")
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is my home page")
}
