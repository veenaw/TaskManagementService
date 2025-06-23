package handlers

import (
	"TaskManagementService/pkg/tasks"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

type TaskTable struct {
	Id          string           `db:"id"`
	Title       string           `db:"title"`
	Description string           `db:"description"`
	DueDate     time.Time        `db:"due_date"`
	Status      tasks.TaskStatus `db:"status"`
}

type TasksDbHandler struct {
	db *sqlx.DB
}

func NewTasksDatabaseHandler(db *sqlx.DB) *TasksDbHandler {
	return &TasksDbHandler{
		db: db,
	}
}

func (h TasksDbHandler) CreateTask(c *gin.Context) {
	var t tasks.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tbl := tasks.TaskFromRequest(t)
	tbl.Id = slug.Make(t.Title)
	stmt := `INSERT INTO task_manager.task (id,title,description,due_date,status)
             VALUES (:id,:title,:description,:due_date,:status)`
	_, err := h.db.NamedExecContext(c.Request.Context(), stmt, tbl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tasks.TaskToResponse(tbl))
}

func (h TasksDbHandler) ListTasks(c *gin.Context) {
	var tasks []TaskTable
	err := h.db.SelectContext(c.Request.Context(), &tasks,
		`SELECT id, title, description, due_date, status FROM task_manager.task ORDER BY due_date`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h TasksDbHandler) FilterTasks(c *gin.Context) {

	type Filter struct {
		Status  *tasks.TaskStatus `form:"status"`
		Before  *time.Time        `form:"before"`
		After   *time.Time        `form:"after"`
		Keyword string            `form:"q"`
	}

	var f Filter
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tasks []TaskTable
	err := h.db.SelectContext(c.Request.Context(), &tasks, `
SELECT id, title, description, due_date, status
FROM task_manager.task
WHERE
  ($1::integer IS NULL OR status   = $1)
  AND ($2::timestamptz IS NULL   OR due_date < $2)
  AND ($3::timestamptz IS NULL   OR due_date > $3)
  AND ( $4 = '' OR title ILIKE '%'||$4||'%' OR description ILIKE '%'||$4||'%' )
ORDER BY due_date
`, f.Status, f.Before, f.After, f.Keyword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h TasksDbHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	var t TaskTable
	err := h.db.GetContext(c.Request.Context(), &t,
		`SELECT id, title, description, due_date, status FROM task_manager.task WHERE id = $1`, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(200, t)
}

func (h *TasksDbHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var input tasks.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map API struct to DB struct
	tbl := tasks.TaskFromRequest(input)
	tbl.Id = id

	// Named parameter update
	res, err := h.db.NamedExecContext(c.Request.Context(), `
    UPDATE task_manager.task SET
      title       = :title,
      description = :description,
      due_date    = :due_date,
      status      = :status
    WHERE id = :id`, tbl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, tasks.TaskToResponse(tbl))
}

func (h TasksDbHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.ExecContext(c.Request.Context(),
		`DELETE FROM task_manager.task WHERE id = $1`, id)

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
