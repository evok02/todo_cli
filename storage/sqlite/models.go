package sqlite

import (
	"database/sql"
	"time"
)

type StatusType string

const (
	PendingStatus StatusType = "in-progress"
	ToDoStatus    StatusType = "todo"
	DoneStatus    StatusType = "done"
)

type Task struct {
	ID          int64
	Description string
	Status      StatusType
	CreatedAt   time.Time
	DeletedAt   time.Time
}

type NullTask struct {
	ID          sql.NullInt64
	Description sql.NullString
	Status      sql.NullString
	CreatedAt   sql.NullTime
	DeletedAt   sql.NullTime
}

func (nt NullTask) ToTask() *Task {
	if !nt.ID.Valid {
		return nil
	}

	t := &Task{
		ID:        nt.ID.Int64,
		CreatedAt: nt.CreatedAt.Time,
	}

	if nt.Description.Valid {
		t.Description = nt.Description.String
	}

	if nt.Status.Valid {
		t.Status = StatusType(nt.Status.String)
	}

	return t
}
