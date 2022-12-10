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

func TestCreateSQLServerCollector(t *testing.T) {
	applicationConfiguration := entities.Application{
		IntervalBetweenRequestInSeconds: 10,
	}
	connection := new(database.SQLConnection)
	sql := database.SQL{Connection: connection}
	repository := database.MSSQLServerRepository{sql}
	datesPersistenceRepository := &filesystem.DatesPersistenceFileRepository{}

	logrus := logging.NewLogrus("info", os.Stdout)
	logger := logrus.Get("Main")
	transmissionChannel := make(chan entities.CapturedData)
	app := NewSQLServer(&repository, datesPersistenceRepository, applicationConfiguration, logger, transmissionChannel)
	assert.Equal(t, app.configuration.IntervalBetweenRequestInSeconds, applicationConfiguration.IntervalBetweenRequestInSeconds)
}
