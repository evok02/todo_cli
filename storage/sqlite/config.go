package sqlite

import (
	"database/sql"
)

type Config struct {
	Path          string
	MigrationPath string
}

var DB *sql.DB
