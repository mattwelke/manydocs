package postgres

import (
	"database/sql"
)

// DocService performs operations using Postgres as a storage engine.
type DocService struct {
	db         *sql.DB
}

// NewDocService creates a new DocService.
func NewDocService(db *sql.DB) DocService {
	return DocService{
		db,
	}
}
