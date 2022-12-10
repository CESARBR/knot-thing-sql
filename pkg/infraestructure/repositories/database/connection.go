package database

import (
	"database/sql"
	//import only by the side effects of the libraries for connection with the DBMSs.
	_ "github.com/AntonioJanael/gocosmos"
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
	go_ora "github.com/sijms/go-ora/v2"
)

type Connection interface {
	Configure(entities.Database, entities.Application) error
	Create() error
	Destroy() error
	GetClient() *sql.DB
}

//General SQL connection
type SQLConnection struct {
	driver             string
	dsn                string
	client             *sql.DB
	maxOpenConnections int
	maxIdleConnections int
}

func (sqlConnection *SQLConnection) Create() error {
	var err error
	sqlConnection.client, err = sql.Open(sqlConnection.driver, sqlConnection.dsn)
	sqlConnection.client.SetMaxOpenConns(sqlConnection.maxOpenConnections)
	sqlConnection.client.SetMaxIdleConns(sqlConnection.maxIdleConnections)
	return err
}

func (sqlConnection *SQLConnection) Destroy() error {
	err := sqlConnection.client.Close()
	return err
}

func (sqlConnection *SQLConnection) Configure(databaseConfiguration entities.Database, applicationConfiguration entities.Application) error {
	sqlConnection.driver = databaseConfiguration.Driver
	sqlConnection.dsn = databaseConfiguration.ConnectionString
	sqlConnection.maxOpenConnections = applicationConfiguration.NumberParallelTags
	sqlConnection.maxIdleConnections = applicationConfiguration.NumberParallelTags
	return nil
}

func (sqlConnection *SQLConnection) GetClient() *sql.DB {
	return sqlConnection.client
}

// SQL connection specific for Oracle databases.
type OracleConnection struct {
	SQLConnection
}

func (oracleConnection *OracleConnection) Configure(databaseConfiguration entities.Database, applicationConfiguration entities.Application) error {
	dsn := go_ora.BuildJDBC(databaseConfiguration.Username, databaseConfiguration.Password, databaseConfiguration.ConnectionString, make(map[string]string))
	oracleConnection.dsn = dsn
	oracleConnection.driver = databaseConfiguration.Driver
	oracleConnection.maxOpenConnections = applicationConfiguration.NumberParallelTags
	oracleConnection.maxIdleConnections = applicationConfiguration.NumberParallelTags
	return nil
}

func NewSQLConnection(databaseConfiguration entities.Database, applicationConfiguration entities.Application) Connection {
	connectionTypeMapping := connectionTypeMappingFactory()
	connection := connectionTypeMapping[applicationConfiguration.Context]
	connection.Configure(databaseConfiguration, applicationConfiguration)
	return connection
}

func connectionTypeMappingFactory() map[string]Connection {
	mapping := make(map[string]Connection)
	mapping[entities.Oracle] = new(OracleConnection)
	mapping[entities.CosmosDB] = new(SQLConnection)
	mapping[entities.SQLServer] = new(SQLConnection)

	return mapping
}
