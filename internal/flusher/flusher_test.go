package flusher_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-offer-api/internal/flusher"
	"github.com/ozoncp/ocp-offer-api/internal/mocks"
	"github.com/ozoncp/ocp-offer-api/internal/models"
)

var _ = Describe("Flusher", func() {
	var (
		ctrl   *gomock.Controller
		m      *mocks.MockRepo
		f      flusher.Flusher
		source []models.Offer
		result []models.Offer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		m = mocks.NewMockRepo(ctrl)
		source = []models.Offer{
			{Id: 10, UserId: 20, Grade: 30, TeamId: 40},
			{Id: 11, UserId: 21, Grade: 31, TeamId: 41},
			{Id: 12, UserId: 22, Grade: 32, TeamId: 42},
			{Id: 13, UserId: 23, Grade: 33, TeamId: 43},
			{Id: 14, UserId: 24, Grade: 34, TeamId: 44},
			{Id: 15, UserId: 25, Grade: 35, TeamId: 45},
			{Id: 16, UserId: 26, Grade: 36, TeamId: 46},
			{Id: 17, UserId: 27, Grade: 37, TeamId: 47},
			{Id: 18, UserId: 28, Grade: 38, TeamId: 48},
			{Id: 19, UserId: 29, Grade: 39, TeamId: 49},
		}
		result = make([]models.Offer, 0)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Empty slice", func() {
		BeforeEach(func() {
			source = make([]models.Offer, 0)
			result = nil
		})

		It("3 chunkSize", func() {
			f = flusher.NewFlusher(3, m)

			m.EXPECT().
				AddOffers(gomock.Any()).
				Times(0)

			res, err := f.Flush(source)

			Ω(err).Should(HaveOccurred())
			Ω(res).Should(Equal(result))
		})
	})

	Context("batches for", func() {
		It("1 chunkSize", func() {
			f := flusher.NewFlusher(1, m)

			m.EXPECT().
				AddOffers(gomock.Any()).
				Times(1)

			res, err := f.Flush(source)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(res).Should(Equal(result))
		})

		It("2 chunkSize", func() {
			f := flusher.NewFlusher(2, m)

			m.EXPECT().
				AddOffers(gomock.Any()).
				Times(2)

			res, err := f.Flush(source)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(res).Should(Equal(result))
		})

		It("5 chunkSize", func() {
			f := flusher.NewFlusher(5, m)

			m.EXPECT().
				AddOffers(gomock.Any()).
				Times(5)

			res, err := f.Flush(source)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(res).Should(Equal(result))
		})
	})

})
