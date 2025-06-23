package tasks

import (
	"time"
)

type TaskTable struct {
	Id          string     `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	DueDate     time.Time  `db:"due_date"`
	Status      TaskStatus `db:"status"`
}
