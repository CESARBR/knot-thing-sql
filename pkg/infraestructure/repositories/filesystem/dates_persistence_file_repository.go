package filesystem

import (
	"os"
	"sync"

	"github.com/CESARBR/knot-thing-sql/internal/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type DatesPersistenceFileRepository struct {
	Filepath string
	Logger   *logrus.Entry
}

var lock sync.Mutex

func (persistence *DatesPersistenceFileRepository) Read() (map[int]string, error) {
	lock.Lock()
	defer lock.Unlock()
	dates, err := utils.ConfigurationParser(persistence.Filepath, make(map[int]string))
	return dates, err
}

func (persistence *DatesPersistenceFileRepository) Update(currentDates map[int]string) error {
	return writeToFile(currentDates, persistence.Filepath, persistence.Logger)
}

func writeToFile(currentDates map[int]string, filepath string, logger *logrus.Entry) error {
	lock.Lock()
	defer lock.Unlock()
	data, err := yaml.Marshal(currentDates)
	if err == nil {
		err = os.WriteFile(filepath, data, 0600)
	}
	if err != nil {
		logger.Printf("Dates persistence error: %v", err)
	}
	return err
}
