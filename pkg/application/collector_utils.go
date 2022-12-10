package application

import (
	"sync"
	"time"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/sirupsen/logrus"
)

//Specifies a function to handle the SQL placeholder for each DBMS,
//which allows to apply the appropriate function depending on the DBMS from which the data is collected.
type sqlPlaceholderHandlerType func(string) string

func newSQLPlaceholderHandlerMapping() map[string]sqlPlaceholderHandlerType {
	mapping := make(map[string]sqlPlaceholderHandlerType)
	mapping[entities.CosmosDB] = removeSQLPlaceholder
	mapping[entities.SQLServer] = keepSQLPlaceholder
	mapping[entities.Oracle] = keepSQLPlaceholder
	return mapping
}
func getSQLPlaceholderMapping(context string) sqlPlaceholderHandlerType {
	mapping := newSQLPlaceholderHandlerMapping()
	return mapping[context]
}

type requestData struct {
	collector                   Collector
	latestRecordTimestampPerTag map[int]string
	id                          int
	mu                          *sync.Mutex
	context                     string
	requestChannel              chan requestData
	logger                      *logrus.Entry
}

func NewRequestData(intervalBetweenRequestInSeconds int, pertinentTags map[int]string, latestRecordTimestampPerTag map[int]string, dataRecoveryPeriodInHours int, collector Collector, context string, channel chan requestData, logger *logrus.Entry) requestData {
	var mutex sync.Mutex
	data := requestData{
		collector:                   collector,
		latestRecordTimestampPerTag: updateLaggedLatestTimestamp(intervalBetweenRequestInSeconds, pertinentTags, latestRecordTimestampPerTag, dataRecoveryPeriodInHours, context),
		mu:                          &mutex,
		context:                     context,
		requestChannel:              channel,
		logger:                      logger,
	}
	return data
}

func NewTicker(intervalBetweenRequestInSeconds int) *time.Ticker {
	return time.NewTicker(time.Duration(intervalBetweenRequestInSeconds) * time.Second)

}

func createGoroutinesPool(numberGoroutines int, dataFromDatabaseChannel chan requestData, collectDataFromDatabase func(requestData)) {
	for i := 0; i < numberGoroutines; i++ {
		go processDatabaseRequest(dataFromDatabaseChannel, collectDataFromDatabase)
	}
}

func hasData(rows []entities.Row, err error) bool {
	const minimumNumberRows = 1
	return err == nil && len(rows) >= minimumNumberRows
}

func processDatabaseRequest(requestDatabaseChannel chan requestData, collectDataFromDatabase func(requestData)) {
	for dataReceivedFromDatabase := range requestDatabaseChannel {
		collectDataFromDatabase(dataReceivedFromDatabase)
	}
}

func collectDataFromDatabase(dataFromDatabase requestData) {
	/*
		Performs the data capture from the database through the specified repository.
	*/
	dataFromDatabase.mu.Lock()
	statement := entities.Statement{
		Timestamp: dataFromDatabase.latestRecordTimestampPerTag[dataFromDatabase.id],
		ID:        dataFromDatabase.id,
	}
	dataFromDatabase.mu.Unlock()
	sqlRows, err := dataFromDatabase.collector.GetDatabaseRepository().Get(statement)
	wasOperationSuccessful := err == nil
	if !wasOperationSuccessful {
		go resendRequest(dataFromDatabase, "Error getting database data", err)
		return
	} else {
		processedRows, err := dataFromDatabase.collector.GetDatabaseRepository().ProcessData(sqlRows)
		if err != nil {
			go resendRequest(dataFromDatabase, "Error processing data", err)
			return
		}
		// The request can be successful and the data processed correctly,
		// even if there are no records in the database for the query performed.
		if hasData(processedRows, err) {
			latestRecordIndex := len(processedRows) - 1
			dataFromDatabase.mu.Lock()
			sqlPlaceholderHandler := getSQLPlaceholderMapping(dataFromDatabase.context)
			dataFromDatabase.latestRecordTimestampPerTag[dataFromDatabase.id] = sqlPlaceholderHandler(processedRows[latestRecordIndex].Timestamp)
			dataFromDatabase.collector.GetDatesRepository().Update(dataFromDatabase.latestRecordTimestampPerTag)
			dataFromDatabase.mu.Unlock()
			capturedData := entities.CapturedData{ID: dataFromDatabase.id, Rows: processedRows}
			go dataFromDatabase.collector.Transmit(capturedData)
		}
	}
}

// Sends the unsuccessful request back to the pool so it can be processed again.
func resendRequest(dataFromDatabase requestData, message string, err error) {
	dataFromDatabase.logger.Errorf("%s: %v", message, err)
	dataFromDatabase.requestChannel <- dataFromDatabase
}

func VerifyError(err error) {
	if err != nil {
		panic(err)
	}
}
