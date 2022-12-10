package application

import (
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories/database"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories/filesystem"
)

type fakeCollector struct{}

func (fake fakeCollector) Collect() {
	//Method used only in tests.
}

func (fake fakeCollector) Transmit(entities.CapturedData) {
	//Method used only in tests.
}

func (fake fakeCollector) GetDatabaseRepository() repositories.DatabaseRepository {
	//Method used only in tests.
	return &database.SQL{}
}

func (fake fakeCollector) GetDatesRepository() repositories.DatesPersistenceRepository {
	//Method used only in tests.
	return &filesystem.DatesPersistenceFileRepository{}
}
