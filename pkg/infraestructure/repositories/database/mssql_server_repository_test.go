package database

import (
	"fmt"
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/stretchr/testify/assert"
)

func TestGivenValidStatementCaptureDataFromSQLServerDatabase(t *testing.T) {
	databaseConfiguration := entities.Database{
		Driver:           "mssql",
		ConnectionString: "",
	}
	applicationConfiguration := entities.Application{
		NumberParallelTags: 10,
		Context:            "sqlserver",
	}
	connection := NewSQLConnection(databaseConfiguration, applicationConfiguration)

	connection.Create()
	defer connection.Destroy()
	queries := entities.Query{}
	sql := SQL{
		Connection: connection,
		Queries:    queries,
	}
	repository := MSSQLServerRepository{sql}

	statement := entities.Statement{
		Timestamp: "",
	}

	rows, err := repository.Get(statement)
	fmt.Println(rows)
	assert.Nil(t, err)
}
