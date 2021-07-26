package utils

import (
	"fmt"
	. "github.com/ozoncp/ocp-contact-api/internal/models"
)

func SliceToMap(contacts []Contact) (map[uint64]Contact, error) {
	result := make(map[uint64]Contact, len(contacts))
	for _, contact := range contacts {
		if _, exist := result[contact.Id]; exist {
			return nil, fmt.Errorf("duplicated indexes %v", contact.Id)
		}
		result[contact.Id] = contact
	}
	return result, nil
}
