package tasks

import (
	"fmt"
	"time"
)

type TaskStatus int

const (
	StatusToDo TaskStatus = iota
	StatusInProgress
	StatusDone
)

// String method for custom string representation
func (s TaskStatus) String() string {
	switch s {
	case StatusToDo:
		return "ToDo"
	case StatusInProgress:
		return "InProgress"
	case StatusDone:
		return "Done"
	default:
		return "Unknown"
	}
}

// Tasks Represents a Task
type Task struct {
	Id          string     `json:"id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	DueDate     time.Time  `json:"dueDate,omitempty"`
	Status      TaskStatus `json:"status,omitempty"`
}

func (t Task) String() string {
	return fmt.Sprint("Id: ", t.Id, " Title: ", t.Title, " Description: ", t.Description, " DueDate: ", t.DueDate, " Status: ", t.Status)
}
