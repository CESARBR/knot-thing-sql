package application

import (
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"math/rand"
	"strconv"
	"time"
)

// DataProducer Generate data
func DataProducer(ch chan<- entities.CapturedData) {
	for {

		rows := []entities.Row{
			{Value: getRandomValue(), Timestamp: time.Now().Format(datetimeLayout)},
		}

		data := entities.CapturedData{
			ID:   rand.Intn(5) + 1,
			Rows: rows,
		}
		ch <- data // To send data through a channel

		// Pause execution for 10 milliseconds before generating the next data.
		time.Sleep(10 * time.Millisecond)
	}
}

// Function to randomly generate value
func getRandomValue() string {
	value := rand.Float64() * 10000
	converted := strconv.FormatFloat(value, 'f', 2, 64)
	return converted
}
