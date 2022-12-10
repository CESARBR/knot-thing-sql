package repositories

import (
	"database/sql"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
)

type DatabaseRepository interface {
	List() ([]string, error)
	Get(entities.Statement) (*sql.Rows, error)
	ProcessData(*sql.Rows) ([]entities.Row, error)
}

type DatesPersistenceRepository interface {
	Read() (map[int]string, error)
	Update(map[int]string) error
}
