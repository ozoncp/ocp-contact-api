package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceToMapRegular(t *testing.T) {
	contacts := makeTestData()
	result, err := SliceToMap(contacts)
	assert.NoError(t, err, "result shouldn't have any errors")
	assert.Len(t, result, len(contacts), "slice and map lengths should be equal")

	for _, contact := range contacts {
		assert.Equal(t, contact, result[contact.Id], "contact should be equal")
	}
}

func TestSliceToMapDuplicate(t *testing.T) {
	contacts := makeTestData()
	contacts = append(contacts, contacts[len(contacts) - 1])

	result, err := SliceToMap(contacts)
	assert.Error(t, err, "result should have error")
	assert.Nil(t, result)
}
