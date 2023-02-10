package vncalendar_test

import (
	"testing"

	. "github.com/openhoangnc/vncalendar"
	"github.com/stretchr/testify/assert"
)

func TestSolar2Lunar(t *testing.T) {
	var result LunarDate

	result = Solar2lunar(2023, 2, 10, 7)
	assert.Equal(t, 20, result.Day)
	assert.Equal(t, 1, result.Month)
	assert.Equal(t, 2023, result.Year)
	assert.Equal(t, false, result.Leap)

	result = Solar2lunar(2023, 2, 11, 7)
	assert.Equal(t, 21, result.Day)
	assert.Equal(t, 1, result.Month)
	assert.Equal(t, 2023, result.Year)
	assert.Equal(t, false, result.Leap)

	result = Solar2lunar(2023, 2, 20, 7)
	assert.Equal(t, 1, result.Day)
	assert.Equal(t, 2, result.Month)
	assert.Equal(t, 2023, result.Year)
	assert.Equal(t, false, result.Leap)

	result = Solar2lunar(2023, 3, 22, 7)
	assert.Equal(t, 1, result.Day)
	assert.Equal(t, 2, result.Month)
	assert.Equal(t, 2023, result.Year)
	assert.Equal(t, true, result.Leap)

	result = Solar2lunar(1903, 1, 1, 7)
	assert.Equal(t, 3, result.Day)
	assert.Equal(t, 12, result.Month)
	assert.Equal(t, 1902, result.Year)
	assert.Equal(t, false, result.Leap)
}
