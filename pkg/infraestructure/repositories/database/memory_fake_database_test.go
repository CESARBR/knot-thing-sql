package database

import (
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestListData(t *testing.T) {
	fake_database_memory := FakeDatabaseMemory{}
	_, err := fake_database_memory.List()
	assert.Nil(t, err)
}

func TestGetData(t *testing.T) {
	statement := entities.Statement{}
	fake_database_memory := FakeDatabaseMemory{}
	_, err := fake_database_memory.Get(statement)
	assert.Nil(t, err)
}
