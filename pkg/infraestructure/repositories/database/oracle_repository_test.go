package database

import (
	"fmt"
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	_ "github.com/sijms/go-ora/v2"
	"github.com/stretchr/testify/assert"
)

func TestGivenValidStatementCaptureDataFromOracleDatabase(t *testing.T) {
	connStr := `(DESCRIPTION=
		(ADDRESS_LIST=
			(ADDRESS=(protocol=TCP)(host=<server>)(port=1521))
		)
		(CONNECT_DATA=
			(SERVICE_NAME=<service>)
		)
		)`
	databaseConfiguration := entities.Database{
		Driver:           "oracle",
		ConnectionString: connStr,
		Username:         "",
		Password:         "",
		Port:             "",
	}
	applicationConfiguration := entities.Application{
		NumberParallelTags: 10,
		Context:            "oracle",
	}
	connection := NewSQLConnection(databaseConfiguration, applicationConfiguration)
	err := connection.Create()
	assert.Nil(t, err)
	defer connection.Destroy()
	queries := entities.Query{Mapping: map[int]string{1: ""}}
	sql := SQL{
		Connection: connection,
		Queries:    queries,
	}
	repository := OracleRepository{sql}

	statement := entities.Statement{
		ID:        1,
		Timestamp: "2022-09-01 11:00:00",
	}

	rows, err := repository.Get(statement)
	assert.Nil(t, err)
	processedRows, err := repository.ProcessData(rows)
	assert.Nil(t, err)
	fmt.Printf("Rows: %v\n", processedRows)
}
