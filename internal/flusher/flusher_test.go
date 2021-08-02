package flusher_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-contact-api/internal/flusher"
	"github.com/ozoncp/ocp-contact-api/internal/mocks"
	"github.com/ozoncp/ocp-contact-api/internal/models"
)

var _ = Describe("Flusher", func() {
	var (
		mockCtrl *gomock.Controller
		mockRepo *mocks.MockRepo
		f flusher.Flusher
		contacts []models.Contact
		chunkSize int
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(mockCtrl)

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
	Context("Flush", func() {
		When("flusher has invalid chunk size", func() {
			BeforeEach(func() {
				chunkSize = 0
				f = flusher.NewFlusher(chunkSize, mockRepo)
			})
			It("returns original slice", func() {
				result, _ := f.Flush(contacts)
				Expect(result).Should(Equal(contacts))
			})
			It("returns error", func() {
				_, err := f.Flush(contacts)
				Expect(err).Should(HaveOccurred())
			})
			It("doesn't called AddContacts from Repo", func() {
				mockRepo.EXPECT().AddContacts(gomock.Any()).Times(0)
				f.Flush(contacts)
			})
		})
		When("repo AddContacts failed", func() {
			BeforeEach(func() {
				chunkSize = 2
				mockRepo.EXPECT().AddContacts(gomock.Any()).Return(errors.New("error"))
				f = flusher.NewFlusher(chunkSize, mockRepo)
			})
			It("returns error and original slice", func() {
				result, err := f.Flush(contacts)
				Expect(result).Should(Equal(contacts))
				Expect(err).Should(HaveOccurred())
			})
		})
		When("repo AddContacts failed in the middle", func() {
			BeforeEach(func() {
				chunkSize = 2
				mockRepo.EXPECT().AddContacts(gomock.Any()).Return(nil).Times(2)
				mockRepo.EXPECT().AddContacts(gomock.Any()).Return(errors.New("error")).Times(1)
				f = flusher.NewFlusher(chunkSize, mockRepo)
			})
			It("returns error and the rest slice", func() {
				result, err := f.Flush(contacts)
				Expect(result).Should(BeEquivalentTo(contacts[2*chunkSize:]))
				Expect(err).Should(HaveOccurred())
			})
		})
		When("AddContacts has no errors", func() {
			BeforeEach(func() {
				chunkSize = 2
				mockRepo.EXPECT().AddContacts(gomock.Any()).Return(nil).AnyTimes()
				f = flusher.NewFlusher(chunkSize, mockRepo)
			})
			It("returns both nil", func() {
				result, err := f.Flush(contacts)
				Expect(result).Should(BeNil())
				Expect(err).Should(BeNil())
			})
		})
	})
})
