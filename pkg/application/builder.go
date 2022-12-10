package application

import (
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/CESARBR/knot-thing-sql/internal/utils"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories/database"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories/filesystem"
	"github.com/CESARBR/knot-thing-sql/pkg/logging"
)

//Defines the interface for building the data collection infrastructure through the Builder design pattern.
type builder interface {
	SetConnection(database.Connection)
	setDatabaseRepository()
	setDatesRepository()
	setCollector()
	setQueries(queries entities.Query)
	SetDatabaseConfiguration(databaseConfiguration entities.Database)
	setProperties(properties BuilderProperties)
	getApplicationHandler() *DataStrategy
}

//Instantiates a builder based on the context specified in the application configuration.
func NewBuilder(buildersMapping map[string]builder, properties BuilderProperties) builder {
	if foundBuilder, ok := buildersMapping[properties.ApplicationConfiguration.Context]; ok {
		foundBuilder.setProperties(properties)
		return foundBuilder
	}
	return nil
}

//Builds a mapping between string and a builder.
func NewBuilderMapping() map[string]builder {
	builders := make(map[string]builder)
	builders[entities.CosmosDB] = new(glassCosmosDBBuilder)
	builders[entities.SQLServer] = new(glassSQLServerBuilder)
	builders[entities.Oracle] = new(restroomOracleBuilder)
	return builders
}

//Basic properties shared by the builders.
type BuilderProperties struct {
	ApplicationConfiguration entities.Application
	Logger                   *logging.Logrus
	TransmissionChannel      chan entities.CapturedData
}

//Base builder composes other builders.
type sqlBuilder struct {
	BuilderProperties
	connection                 database.Connection
	databaseRepository         repositories.DatabaseRepository
	datesPersistenceRepository repositories.DatesPersistenceRepository
	collector                  Collector
	databaseConfiguration      entities.Database
	queries                    entities.Query
}

func (sqlb *sqlBuilder) SetConnection(connection database.Connection) {
	sqlb.connection = connection
}
func (sqlb *sqlBuilder) setDatabaseRepository() {
	sqlRepository := database.SQL{
		Connection: sqlb.connection,
		Queries:    sqlb.queries,
	}
	sqlb.databaseRepository = &database.CosmosDBRepository{sqlRepository}
}
func (sqlb *sqlBuilder) setDatesRepository() {
	sqlb.datesPersistenceRepository = &filesystem.DatesPersistenceFileRepository{
		Filepath: sqlb.ApplicationConfiguration.DatesPersistenceFilepath,
		Logger:   sqlb.Logger.Get("Dates persistence"),
	}
}
func (sqlb *sqlBuilder) setCollector() {
	sqlb.collector = NewCosmosDB(sqlb.databaseRepository, sqlb.datesPersistenceRepository, sqlb.ApplicationConfiguration, sqlb.Logger.Get("KNoT SQL CosmosDB"), sqlb.TransmissionChannel)
}

func (sqlb *sqlBuilder) SetApplicationConfiguration(appliactionConfiguration entities.Application) {
	sqlb.ApplicationConfiguration = appliactionConfiguration
}

func (sqlb *sqlBuilder) SetLogger(logger *logging.Logrus) {
	sqlb.Logger = logger
}

func (sqlb *sqlBuilder) SetTransmissionChannel(transmissionChannel chan entities.CapturedData) {
	sqlb.TransmissionChannel = transmissionChannel
}

func (sqlb *sqlBuilder) setQueries(queries entities.Query) {
	sqlb.queries = queries
}

func (sqlb *sqlBuilder) SetDatabaseConfiguration(databaseConfiguration entities.Database) {
	sqlb.databaseConfiguration = databaseConfiguration
}

func (sqlb *sqlBuilder) setProperties(properties BuilderProperties) {
	sqlb.ApplicationConfiguration = properties.ApplicationConfiguration
	sqlb.Logger = properties.Logger
	sqlb.TransmissionChannel = properties.TransmissionChannel

}

func (sqlb *sqlBuilder) getApplicationHandler() *DataStrategy {
	return NewDataStrategy(sqlb.collector)

}

//Builder for a specific cosmosDB data collection scenario.
//Some methods are overridden, while others are leveraged from the composite structure.
type glassCosmosDBBuilder struct {
	sqlBuilder
}

func (gcb *glassCosmosDBBuilder) setDatabaseRepository() {
	sqlRepository := database.SQL{
		Connection: gcb.connection,
		Queries:    gcb.queries,
	}
	gcb.databaseRepository = &database.CosmosDBRepository{sqlRepository}
}
func (gcb *glassCosmosDBBuilder) setCollector() {
	gcb.collector = NewSQLServer(gcb.databaseRepository, gcb.datesPersistenceRepository, gcb.ApplicationConfiguration, gcb.Logger.Get("KNoT SQL CosmosDB"), gcb.TransmissionChannel)
}

//Builder for a specific SQL Server data collection scenario.
type glassSQLServerBuilder struct {
	sqlBuilder
}

func (gsb *glassSQLServerBuilder) setDatabaseRepository() {
	sqlRepository := database.SQL{
		Connection: gsb.connection,
		Queries:    gsb.queries,
	}
	gsb.databaseRepository = &database.MSSQLServerRepository{sqlRepository}
}
func (gsb *glassSQLServerBuilder) setCollector() {
	gsb.collector = NewSQLServer(gsb.databaseRepository, gsb.datesPersistenceRepository, gsb.ApplicationConfiguration, gsb.Logger.Get("KNoT SQL SQL Server"), gsb.TransmissionChannel)
}

//Builder for a specific Oracle data collection scenario.
type restroomOracleBuilder struct {
	sqlBuilder
}

func (rob *restroomOracleBuilder) setDatabaseRepository() {
	sqlRepository := database.SQL{
		Connection: rob.connection,
		Queries:    rob.queries,
	}
	rob.databaseRepository = &database.OracleRepository{sqlRepository}
}
func (rob *restroomOracleBuilder) setCollector() {
	rob.collector = NewOracle(rob.databaseRepository, rob.datesPersistenceRepository, rob.ApplicationConfiguration, rob.Logger.Get("KNoT SQL Oracle"), rob.TransmissionChannel)
}

//Structure that represents the director who builds the data collection infrastructure.
type director struct {
	builder builder
}

func NewDirector(b builder) *director {
	return &director{builder: b}
}

func (d *director) SetBuilder(b builder) {
	d.builder = b
}

//Builds the data collection infrastructure.
func (d *director) BuildDataHandler() *DataStrategy {
	queriesConfiguration, err := utils.ConfigurationParser("internal/scripts/queries.yaml", entities.Query{})
	VerifyError(err)
	d.builder.setQueries(queriesConfiguration)
	d.builder.setDatabaseRepository()
	d.builder.setDatesRepository()
	d.builder.setCollector()
	return d.builder.getApplicationHandler()
}
