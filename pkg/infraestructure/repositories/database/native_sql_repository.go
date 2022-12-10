package database

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
)

type SQL struct {
	Connection Connection
	Queries    entities.Query
}

func (s *SQL) List() ([]string, error) {
	return []string{}, nil

}

func (s *SQL) Get(statement entities.Statement) (*sql.Rows, error) {
	query := fmt.Sprintf(s.Queries.Mapping[statement.ID], statement.Timestamp)
	rows, err := s.Connection.GetClient().Query(query)
	if err != nil {
		return &sql.Rows{}, err
	}
	return rows, nil
}

func (s *SQL) ProcessData(rows *sql.Rows) ([]entities.Row, error) {
	var value, timestamp string
	const length = 0
	processedRows := make([]entities.Row, length)
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&value, &timestamp); err != nil {
			return []entities.Row{}, err
		}

		if err := VerifyTimestamp(timestamp); err != nil {
			return []entities.Row{}, err
		}

		processedRows = append(processedRows, entities.Row{
			Value:     value,
			Timestamp: timestamp,
		})
	}
	err := rows.Err()
	if err != nil {
		return []entities.Row{}, err
	}
	return processedRows, nil
}

func VerifyTimestamp(timeStamp string) error {
	/*
		Expected format for the timestamp: "2021-12-31 23:59:59".
	*/
	formatOfTimeSTamp := "[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) (2[0-3]|[0-1][0-9]):[0-5][0-9]:[0-5][0-9]"
	regularEx, _ := regexp.Compile(formatOfTimeSTamp)
	if regularEx.MatchString(timeStamp) {
		return nil
	} else {
		return fmt.Errorf("timestamp format did not expected")
	}
}
