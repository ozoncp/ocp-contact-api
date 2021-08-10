package saver

import (
	"github.com/ozoncp/ocp-contact-api/internal/flusher"
	"github.com/ozoncp/ocp-contact-api/internal/models"
	"sync"
	"time"
)

type Saver interface {
	Save(entity models.Contact)
	Close()
}

type saver struct {
	capacity uint
	flusher flusher.Flusher
	ticker *time.Ticker
	contact chan models.Contact
	close chan struct{}
	wait *sync.WaitGroup
}

func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
	timeout time.Duration,
) Saver {
	s := &saver {
		capacity: capacity,
		flusher: flusher,
		ticker: time.NewTicker(timeout),
		contact: make(chan models.Contact),
		close: make(chan struct{}),
		wait: &sync.WaitGroup{},
	}
	s.init()
	return s
}

func (s *saver)init() {
	contacts := make([]models.Contact, 0, s.capacity)
	s.wait.Add(1)
	go func() {
		defer func() {
			s.ticker.Stop()
			close(s.contact)
			s.flusher.Flush(contacts)
			s.wait.Done()
		}()

		for {
			select {
			case <-s.close:
				return
			case <- s.ticker.C:
				contacts, _ = s.flusher.Flush(contacts)
			case contact := <- s.contact:
				contacts = append(contacts, contact)
			}
		}
	}()
}

func (s *saver)Save(contact models.Contact) {
	select {
	case <- s.close:
		return
	default:
		s.contact <- contact
	}
}

func (s *saver) Close() {
	close(s.close)
	s.wait.Wait()
}