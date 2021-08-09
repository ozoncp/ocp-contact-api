package saver

import (
	"github.com/ozoncp/ocp-contact-api/internal/flusher"
	"github.com/ozoncp/ocp-contact-api/internal/models"
	"time"
)

type Saver interface {
	Save(entity models.Contact)
	Close()
}

type saver struct {
	flusher flusher.Flusher
	ticker *time.Ticker
	contact chan models.Contact
	close chan struct{}
}

func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
	timeout time.Duration,
) Saver {
	s := &saver {
		flusher: flusher,
		ticker: time.NewTicker(timeout),
		contact: make(chan models.Contact, capacity),
		close: make(chan struct{}),
	}
	s.init()
	return s
}

func (s *saver)init() {
	contacts := make([]models.Contact, 0, cap(s.contact))
	go func() {
		defer func() {
			s.ticker.Stop()
			close(s.contact)
			close(s.close)
			s.flusher.Flush(contacts)
		}()

		for {
			select {
			case contact := <- s.contact:
				contacts = append(contacts, contact)
			case <- s.ticker.C:
				contacts, _ = s.flusher.Flush(contacts)
			case <-s.close:
				return
			}
		}
	}()
}

func (s *saver)Save(contact models.Contact) {
	s.contact <- contact
}

func (s *saver) Close() {
	s.close <- struct{}{}
}