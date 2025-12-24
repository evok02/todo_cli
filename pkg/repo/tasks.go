package repo

import (
	"database/sql"
	"fmt"
	"github.com/evok02/todo_cli/storage/sqlite"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r Repo) CreateTask(desc string) (*sqlite.Task, error) {
	q := `
		INSERT INTO tasks (description, status, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		RETURNING *
	`
	var task sqlite.NullTask
	err := r.db.QueryRow(q, desc, sqlite.ToDoStatus).Scan(
		&task.ID,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.DeletedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("couldn't find appropriate row: %w", err)
	case err != nil:
		return nil, fmt.Errorf("couldn't execute query: %w", err)
	default:
		fmt.Printf("Created task:\n%s \nID: %d\n", desc, task.ID.Int64)
	}

	return task.ToTask(), nil
}

func (r Repo) UpdateTask(id int, desc string) error {
	q := `
	UPDATE tasks
	SET description = ?
	WHERE id = ? AND deleted_at IS NULL
	`
	res, err := r.db.Exec(q, desc, id)

	if err != nil {
		return fmt.Errorf("couldn't execute query: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could't find appropriate row: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("failed to affect any rows")
	}

	fmt.Println("Updated successfuly!")
	return nil
}

func (r Repo) SoftDeleteTask(id int) error {
	q := `
	UPDATE tasks
	SET deleted_at = CURRENT_TIMESTAMP
	WHERE id = ? and deleted_at is NULL
	`
	res, err := r.db.Exec(q, id)

	if err != nil {
		return fmt.Errorf("couldn't execute query: %w", err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("could't find appropriate row: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("failed to delete any rows")
	}

	fmt.Println("Deleted successfuly!")
	return nil
}

func (r Repo) HardDeleteTask(id int) (*sqlite.Task, error) {
	q := `
	DELETE FROM tasks WHERE id = ?
	RETURNING id, description 
	`
	var task sqlite.NullTask
	err := r.db.QueryRow(q, id).Scan(
		&task.ID,
		&task.Description,
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("couldn't find appropriate row: %w", err)
	case err != nil:
		return nil, fmt.Errorf("couldn't execute query: %w", err)
	default:
		fmt.Printf("Deleted task:\n%s \nID: %d\n", task.Description.String, task.ID.Int64)
	}

	return task.ToTask(), nil
}

func (r Repo) MarkInProgress(id int) error {
	q := `
	UPDATE tasks
	SET status = ?
	WHERE id = ? and DELETED_AT IS NULL
	`
	res, err := r.db.Exec(q, sqlite.InProgressStatus, id)

	if err != nil {
		return fmt.Errorf("couldn't execute query: %w", err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("could't check affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("could't affect any rows")
	}

	fmt.Println("Edited successfuly!")
	return nil
}

func (r Repo) MarkDone(id int) error {
	q := `
	UPDATE tasks
	SET status = ?
	WHERE id = ? and DELETED_AT IS NULL
	`
	res, err := r.db.Exec(q, sqlite.DoneStatus, id)

	if err != nil {
		return fmt.Errorf("couldn't execute query: %w", err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("could't check affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("could't affect any rows")
	}

	fmt.Println("Edited successfuly!")
	return nil

}

func (r Repo) List(args ...string) ([]*sqlite.Task, error) {
	tasks := []*sqlite.Task{}
	var q string
	switch args[0] {
	case string(sqlite.InProgressStatus):
		q = `
		SELECT id, description, status, created_at FROM tasks
		WHERE deleted_at IS NULL AND status = 'in-progress'
		`
	case string(sqlite.ToDoStatus):
		q = `
		SELECT id, description, status, created_at FROM tasks
		WHERE deleted_at IS NULL AND status = 'todo'
		`
	case string(sqlite.DoneStatus):
		q = `
		SELECT id, description, status, created_at FROM tasks
		WHERE deleted_at IS NULL AND status = 'done'
		`
	default:
		q = `
		SELECT id, description, status, created_at FROM tasks
		WHERE deleted_at IS NULL
		`
	}

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("couldn't execute the query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var task sqlite.NullTask

		err := rows.Scan(
			&task.ID,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("couldn't extract the row: %w", err)
		}

		fmtTask := task.ToTask()
		tasks = append(tasks, fmtTask)
	}

	return tasks, nil
}
