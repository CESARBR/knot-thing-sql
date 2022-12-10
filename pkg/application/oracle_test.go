package application

import (
	"os"
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories/database"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories/filesystem"
	"github.com/CESARBR/knot-thing-sql/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestCreateOracleCollector(t *testing.T) {
	applicationConfiguration := entities.Application{
		IntervalBetweenRequestInSeconds: 10,
		Context:                         "oracle",
	}
	databaseConfiguration := entities.Database{}
	connection := database.NewSQLConnection(databaseConfiguration, applicationConfiguration)
	sql := database.SQL{Connection: connection}
	repository := database.OracleRepository{sql}
	datesPersistenceRepository := &filesystem.DatesPersistenceFileRepository{}

	logrus := logging.NewLogrus("info", os.Stdout)
	logger := logrus.Get("Main")
	transmissionChannel := make(chan entities.CapturedData)
	app := NewOracle(&repository, datesPersistenceRepository, applicationConfiguration, logger, transmissionChannel)
	assert.Equal(t, app.configuration.IntervalBetweenRequestInSeconds, applicationConfiguration.IntervalBetweenRequestInSeconds)
}
