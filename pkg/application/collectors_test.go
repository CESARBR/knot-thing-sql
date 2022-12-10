package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCollector(t *testing.T) {
	fakeCollector := fakeCollector{}
	data := NewDataStrategy(fakeCollector)
	assert.NotNil(t, data)
}

func TestSetCollector(t *testing.T) {
	fakeCollector := fakeCollector{}
	data := NewDataStrategy(fakeCollector)
	err := data.SetCollectorStrategy(fakeCollector)
	assert.Nil(t, err)
}

func TestCollectCollector(t *testing.T) {
	fakeCollector := fakeCollector{}
	data := NewDataStrategy(fakeCollector)
	data.Collect()
	assert.NotNil(t, data)
}
