package saver_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/ozoncp/ocp-contact-api/internal/mocks"
	"github.com/ozoncp/ocp-contact-api/internal/models"
	"github.com/ozoncp/ocp-contact-api/internal/saver"
	"time"
)

var _ = Describe("Saver", func() {
	var (
		mockCtrl *gomock.Controller
		mockFlusher *mocks.MockFlusher
		s saver.Saver
		timeout time.Duration
		contacts []models.Contact
		capacity uint
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFlusher = mocks.NewMockFlusher(mockCtrl)
		timeout = 100 * time.Millisecond

		contacts = []models.Contact {
			0: {1 , 1, 41, "one"},
			1: {2 , 2, 42, "two"},
			2: {3 , 3, 43, "three"},
			3: {4 , 4, 44, "four"},
			4: {5 , 5, 45, "five"},
			5: {6 , 6, 46, "six"},
		}
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})
	Context("Save", func() {
		When("capacity is empty", func() {
			BeforeEach(func() {
				capacity = 0
				s = saver.NewSaver(capacity, mockFlusher, timeout)
				mockFlusher.EXPECT().Flush(contacts).Times(1)
			})
			It("flushes all contacts after close", func() {
				defer s.Close()
				for _, contact := range contacts {
					s.Save(contact)
				}
			})
		})
		When("capacity is less than collection size", func() {
			BeforeEach(func() {
				capacity = uint(len(contacts) / 2)
				s = saver.NewSaver(capacity, mockFlusher, timeout)
				mockFlusher.EXPECT().Flush(gomock.Any()).Times(1)
			})
			It("flushes all contacts", func() {
				defer s.Close()
				for _, contact := range contacts {
					s.Save(contact)
				}
			})
		})
		It("flushes contacts by timeout", func() {
			capacity = uint(len(contacts))
			timeout = time.Millisecond
			s = saver.NewSaver(capacity, mockFlusher, timeout)
			mockFlusher.EXPECT().Flush(gomock.Any()).MinTimes(1)

			for _, contact := range contacts {
				s.Save(contact)
			}
			time.Sleep(50 * time.Millisecond)
		})
	})
})
