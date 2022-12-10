package database

import (
	"fmt"
	"testing"

	_ "github.com/AntonioJanael/gocosmos"
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGivenValidStatementCaptureDataFromDatabase(t *testing.T) {
	databaseConfiguration := entities.Database{
		Driver:           "gocosmos",
		ConnectionString: "AccountEndpoint=;AccountKey=;DefaultDb=",
	}
	applicationConfiguration := entities.Application{
		NumberParallelTags: 10,
		Context:            "cosmosdb",
	}
	connection := NewSQLConnection(databaseConfiguration, applicationConfiguration)
	err := connection.Create()
	assert.Nil(t, err)
	defer connection.Destroy()
	queries := entities.Query{Mapping: map[int]string{1: ""}}
	sql := SQL{
		Connection: connection,
		Queries:    queries}

	repository := CosmosDBRepository{sql}

	statement := entities.Statement{
		ID:        1,
		Timestamp: "2022-09-30 11`:`00`:`00",
	}

	rows, err := repository.Get(statement)
	assert.Nil(t, err)
	processedRows, err := repository.ProcessData(rows)
	assert.Nil(t, err)
	fmt.Printf("Rows: %v\n", processedRows)
}
