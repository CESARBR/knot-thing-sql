package application

import (
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories"
	"github.com/sirupsen/logrus"
)

//Defines a structure for SQL Server composed of the SQL collector.
type sqlServer struct {
	SQLCollector
}

func NewSQLServer(databaseRepository repositories.DatabaseRepository, datesPersistenceRepository repositories.DatesPersistenceRepository, configuration entities.Application, logger *logrus.Entry, transmissionChannel chan entities.CapturedData) *sqlServer {
	sqlServerCollector := new(sqlServer)
	sqlServerCollector.databaseRepository = databaseRepository
	sqlServerCollector.datesPersistenceRepository = datesPersistenceRepository
	sqlServerCollector.configuration = configuration
	sqlServerCollector.logger = logger
	sqlServerCollector.transmissionChannel = transmissionChannel
	return sqlServerCollector
}
