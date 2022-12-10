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

func TestCreateApplication(t *testing.T) {
	applicationConfiguration := entities.Application{
		IntervalBetweenRequestInSeconds: 10,
	}

	repository := &database.SQL{}
	datesPersistenceRepository := &filesystem.DatesPersistenceFileRepository{}

	logrus := logging.NewLogrus("info", os.Stdout)
	logger := logrus.Get("Main")
	transmissionChannel := make(chan entities.CapturedData)
	app := NewCosmosDB(repository, datesPersistenceRepository, applicationConfiguration, logger, transmissionChannel)
	assert.Equal(t, app.configuration.IntervalBetweenRequestInSeconds, applicationConfiguration.IntervalBetweenRequestInSeconds)
}

func TestDataTransmission(t *testing.T) {
	log := logging.NewLogrus("info", os.Stdout)
	transmissionChannel := make(chan entities.CapturedData)
	row := entities.Row{}
	testRows := make([]entities.Row, 0)
	testRows = append(testRows, row)
	capturedData := entities.CapturedData{Rows: testRows}
	cosmos := NewCosmosDB(&database.SQL{}, &filesystem.DatesPersistenceFileRepository{}, entities.Application{}, log.Get(""), transmissionChannel)
	go cosmos.Transmit(capturedData)
	<-transmissionChannel
}
