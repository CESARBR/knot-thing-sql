package application

import (
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories"
	"github.com/sirupsen/logrus"
)

//Defines a structure for CosmosDB composed of the SQL collector.
type cosmosDB struct {
	SQLCollector
}

func NewCosmosDB(databaseRepository repositories.DatabaseRepository, datesPersistenceRepository repositories.DatesPersistenceRepository, configuration entities.Application, logger *logrus.Entry, transmissionChannel chan entities.CapturedData) *cosmosDB {
	cosmos := new(cosmosDB)
	cosmos.databaseRepository = databaseRepository
	cosmos.datesPersistenceRepository = datesPersistenceRepository
	cosmos.configuration = configuration
	cosmos.logger = logger
	cosmos.transmissionChannel = transmissionChannel
	return cosmos
}
