package flusher

import (
	"github.com/ozoncp/ocp-contact-api/internal/models"
	"github.com/ozoncp/ocp-contact-api/internal/repo"
	"github.com/ozoncp/ocp-contact-api/internal/utils"
)

type Flusher interface {
	Flush(contacts []models.Contact) ([]models.Contact, error)
}

type flusher struct {
	chunkSize int
	contactRepo  repo.Repo
}

func NewFlusher(
	chunkSize int,
	contactRepo repo.Repo,
) Flusher {
	return &flusher{
		chunkSize: chunkSize,
		contactRepo:  contactRepo,
	}
}

func (f *flusher) Flush(contacts []models.Contact) ([]models.Contact, error) {
	chunks, err := utils.Split(contacts, f.chunkSize)
	if err != nil {
		return contacts, err
	}
	for index := range chunks {
		if err := f.contactRepo.AddContacts(chunks[index]); err != nil {
			return contacts[f.chunkSize * index:], err
		}
	}
	return nil, nil
}
