package database

import (
	"fmt"
	"testing"

	_ "github.com/AntonioJanael/gocosmos"
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGivenValidSQLConnectionDriverConnect(t *testing.T) {
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
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := connection.Destroy()
		if err != nil {
			fmt.Println(err)
		}
	}()
	assert.Nil(t, err)
}

func TestSQLServerConnection(t *testing.T) {
	databaseConfiguration := entities.Database{
		Driver:           "mssql",
		ConnectionString: "server=<server_address>;user id=<user>;password=<password>;port=<port>;database=<database>",
	}
	applicationConfiguration := entities.Application{
		NumberParallelTags: 10,
		Context:            "sqlserver",
	}
	mssql := NewSQLConnection(databaseConfiguration, applicationConfiguration)
	openErr := mssql.Create()
	pingErr := mssql.GetClient().Ping()
	closeErr := mssql.Destroy()
	assert.Nil(t, openErr)
	assert.Nil(t, closeErr)
	assert.Nil(t, pingErr)
}

func TestOracleConnection(t *testing.T) {
	connectionString := `(DESCRIPTION=
		(ADDRESS_LIST=
			(ADDRESS=(protocol=TCP)(host=<server>)(port=1521))
		)
		(CONNECT_DATA=
			(SERVICE_NAME=<service>)
		)
		)`
	databaseConfiguration := entities.Database{
		Driver:           "oracle",
		ConnectionString: connectionString,
		Username:         "",
		Password:         "",
		Port:             "",
	}
	applicationConfiguration := entities.Application{
		NumberParallelTags: 10,
		Context:            "oracle",
	}
	mssql := NewSQLConnection(databaseConfiguration, applicationConfiguration)
	openErr := mssql.Create()
	pingErr := mssql.GetClient().Ping()
	closeErr := mssql.Destroy()
	assert.Nil(t, openErr)
	assert.Nil(t, closeErr)
	assert.Nil(t, pingErr)
}
