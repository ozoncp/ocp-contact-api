package repo

import (
	"github.com/ozoncp/ocp-contact-api/internal/models"
)

type Repo interface {
	AddContacts(entities []models.Contact) error
	ListContacts(limit, offset uint64) ([]models.Contact, error)
	DescribeContact(entityId uint64) (*models.Contact, error)
}
