package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertDatetimeToString(t *testing.T) {
	datetimeString := "2022-06-15T08:30:00.0Z"
	convertedDatetime, err := convertStringToDatetime(datetimeString)
	assert.Nil(t, err)
	dateString := convertDatetimeToString(convertedDatetime, "glassCosmosDB")
	expectedconvertedDatetime := "2022-06-15 08`:`30`:`00"
	assert.Equal(t, expectedconvertedDatetime, dateString)
}

func TestFormatStringDatetimeToUTCWithZone(t *testing.T) {
	datetimeString := "2022-06-15 08`:`30`:`00"
	convertedDatetime := formatTimestampToUTCWithZone(datetimeString)
	expectedconvertedDatetime := "2022-06-15T08:30:00.0Z"
	assert.Equal(t, expectedconvertedDatetime, convertedDatetime)
}
func TestCompareDatetime(t *testing.T) {
	datetimeString := "2022-06-15 08`:`30`:`00"
	datetimeStringWithTimeZone := formatTimestampToUTCWithZone(datetimeString)
	convertedDatetime, err := convertStringToDatetime(datetimeStringWithTimeZone)
	assert.Nil(t, err)
	isTimeLagged := isTimeLagged(convertedDatetime, 5)
	const expectedTimeLagged = true
	assert.Equal(t, expectedTimeLagged, isTimeLagged)
}

func TestUpdateDatetimeWhenLaggedTime(t *testing.T) {
	datetimeString := "2022-06-15 20`:`30`:`00"
	datetimeStringWithTimeZone := formatTimestampToUTCWithZone(datetimeString)
	convertedDatetime, err := convertStringToDatetime(datetimeStringWithTimeZone)
	assert.Nil(t, err)
	isTimeLagged := isTimeLagged(convertedDatetime, 5)

	if isTimeLagged {
		const laggedHours = 5
		laggedTime := lagTime(convertedDatetime, laggedHours)
		convertedDatetimeString := convertDatetimeToString(laggedTime, "glassCosmosDB")
		expectedDatetimeString := "2022-06-15 15`:`30`:`00"
		assert.Equal(t, expectedDatetimeString, convertedDatetimeString)
	}

}

func TestIsLaggedTime(t *testing.T) {
	timestamp := "2022-07-15 19`:`40`:`00"
	laggedHours := 1
	var datetimeStringWithTimeZone string = formatTimestampToUTCWithZone(timestamp)
	convertedDatetime, err := convertStringToDatetime(datetimeStringWithTimeZone)
	if err == nil {
		islagged := isTimeLagged(convertedDatetime, laggedHours)
		assert.Equal(t, islagged, false)
	}
}
