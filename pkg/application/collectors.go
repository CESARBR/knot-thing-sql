package application

import (
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories"
	"github.com/sirupsen/logrus"
)

type Collector interface {
	Collect()
	Transmit(entities.CapturedData)
	GetDatabaseRepository() repositories.DatabaseRepository
	GetDatesRepository() repositories.DatesPersistenceRepository
}

// Defines the basic structure of an SQL collector that serves to compose other collectors.
type SQLCollector struct {
	databaseRepository         repositories.DatabaseRepository
	datesPersistenceRepository repositories.DatesPersistenceRepository
	configuration              entities.Application
	logger                     *logrus.Entry
	transmissionChannel        chan entities.CapturedData
}

func (sql SQLCollector) Collect() {
	requestDatabaseChannel := make(chan requestData)
	// creates a pool of goroutines that manages the processing of requests to the database.
	createGoroutinesPool(sql.configuration.NumberParallelTags, requestDatabaseChannel, collectDataFromDatabase)
	ticker := NewTicker(sql.configuration.IntervalBetweenRequestInSeconds)
	for range ticker.C {
		latestRecordTimestampPerTag, _ := sql.datesPersistenceRepository.Read()
		data := NewRequestData(sql.configuration.IntervalBetweenRequestInSeconds, sql.configuration.PertinentTags, latestRecordTimestampPerTag, sql.configuration.DataRecoveryPeriodInHours, sql, sql.configuration.Context, requestDatabaseChannel, sql.logger)
		for id := range sql.configuration.PertinentTags {
			data.mu.Lock()
			data.id = id
			data.mu.Unlock()
			requestDatabaseChannel <- data
		}
	}
}

func (sql SQLCollector) Transmit(data entities.CapturedData) {
	sql.transmissionChannel <- data
}

func (sql SQLCollector) GetDatabaseRepository() repositories.DatabaseRepository {
	return sql.databaseRepository
}

func (sql SQLCollector) GetDatesRepository() repositories.DatesPersistenceRepository {
	return sql.datesPersistenceRepository
}

// Defines the basic structure for the collectors' strategy.
type DataStrategy struct {
	collector Collector
}

func NewDataStrategy(collector Collector) *DataStrategy {
	return &DataStrategy{
		collector: collector,
	}
}

func (data *DataStrategy) SetCollectorStrategy(collector Collector) error {
	data.collector = collector
	return nil
}

func (data *DataStrategy) Collect() {
	data.collector.Collect()
}
