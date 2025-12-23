package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func New(config Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.Path)

	if err != nil {
		return nil, fmt.Errorf("can't open the database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to the database: %w", err)
	}

	return db, nil
}

func Init() (*sql.DB, error) {
	config := Config{
		Path: "./data/sqlite/tasks.db",
	}

	db, err := New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create new connection: %w", err)
	}

	DB = db

	err = runMigrations()
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return DB, nil
}

func runMigrations() error {
	q := `
	CREATE TABLE IF NOT EXISTS tasks ( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL ,
		status TEXT NOT NULL DEFAULT 'todo',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME);

		CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
    	CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
    	CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);
	`
	_, err := DB.Exec(q)

	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
