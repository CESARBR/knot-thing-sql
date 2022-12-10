package application

import (
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories"
	"github.com/sirupsen/logrus"
)

type oracle struct {
	SQLCollector
}

func NewOracle(databaseRepository repositories.DatabaseRepository, datesPersistenceRepository repositories.DatesPersistenceRepository, configuration entities.Application, logger *logrus.Entry, transmissionChannel chan entities.CapturedData) *oracle {
	oracleCollector := new(oracle)
	oracleCollector.databaseRepository = databaseRepository
	oracleCollector.datesPersistenceRepository = datesPersistenceRepository
	oracleCollector.configuration = configuration
	oracleCollector.logger = logger
	oracleCollector.transmissionChannel = transmissionChannel
	return oracleCollector
}
