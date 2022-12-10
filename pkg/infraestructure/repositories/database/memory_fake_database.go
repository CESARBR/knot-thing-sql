package database

import (
	"github.com/CESARBR/knot-thing-sql/internal/entities"
)

type FakeDatabaseMemory struct {
}

func (fake FakeDatabaseMemory) List() ([]string, error) {
	return []string{"Fake data"}, nil
}

func (fake FakeDatabaseMemory) Get(statement entities.Statement) ([]entities.Row, error) {
	const length = 0
	rows := make([]entities.Row, length)
	rows = append(rows, entities.Row{
		Timestamp: "0000-00-00",
	})
	return rows, nil
}
