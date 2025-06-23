package handlers

import (
	"TaskManagementService/pkg/tasks"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"net/http"
)

type TasksHandler struct {
	store tasksStore
}

func NewTasksHandler(s tasksStore) *TasksHandler {
	return &TasksHandler{
		store: s,
	}
}

type tasksStore interface {
	Add(id string, task tasks.Task) error
	Get(id string) (tasks.Task, error)
	List() (map[string]tasks.Task, error)
	Update(id string, task tasks.Task) error
	Remove(id string) error
}

func (h TasksHandler) CreateTask(c *gin.Context) {
	// Get request body and convert it to task.Task
	var task tasks.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create a url friendly name
	id := slug.Make(task.Title)
	task.Id = id

	// add to the store
	err := h.store.Add(id, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed"})
	}

	// return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h TasksHandler) ListTasks(c *gin.Context) {
	r, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, r)
}

func (h TasksHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.store.Get(id)
	fmt.Println(task.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(200, task)
}

func (h TasksHandler) UpdateTask(c *gin.Context) {
	// Get request body and convert it to tasks.Task
	var task tasks.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	task.Id = slug.Make(task.Title)

	err := h.store.Update(id, task)

	if err != nil {
		if err == tasks.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h TasksHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := h.store.Remove(id)
	if err != nil {
		if err == tasks.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})

}
