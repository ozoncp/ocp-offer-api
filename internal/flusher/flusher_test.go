package flusher_test

import (
	"errors"

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
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("save offers to repo with flusher", func() {
		Context("when AddOffers returns error", func() {
			It("returns error", func() {
				var result []models.Offer = nil
				f = flusher.NewFlusher(3, m)

				m.EXPECT().
					AddOffers(gomock.Any()).
					Return(errors.New("error")).
					Times(1)

				res, err := f.Flush(source)

				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeEquivalentTo(errors.New("error")))
				Expect(res).Should(Equal(result))
			})
		})

		Context("not normal case", func() {
			Context("when chunk < 0", func() {
				It("returns error", func() {
					var result []models.Offer = nil
					f = flusher.NewFlusher(-1, m)

					m.EXPECT().
						AddOffers(gomock.Any()).
						Times(0)

					res, err := f.Flush(source)

					Expect(err).Should(HaveOccurred())
					Expect(res).Should(Equal(result))
				})
			})

			Context("when source is empty", func() {
				It("returns error", func() {
					var result []models.Offer = nil
					f := flusher.NewFlusher(1, m)

					m.EXPECT().
						AddOffers(gomock.Any()).
						Times(0)

					res, err := f.Flush(make([]models.Offer, 0))

					Expect(err).Should(HaveOccurred())
					Expect(res).Should(Equal(result))
				})
			})
		})

		Context("normal case", func() {

			Context("when chunk = 1", func() {
				It("works without errors", func() {
					f := flusher.NewFlusher(1, m)

					m.EXPECT().
						AddOffers(gomock.Any()).
						Times(1)

					res, err := f.Flush(source)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res).Should(Equal(source))
				})
			})

			Context("when chunk = 2", func() {
				It("works without errors", func() {
					f := flusher.NewFlusher(2, m)

					m.EXPECT().
						AddOffers(gomock.Any()).
						Times(2)

					res, err := f.Flush(source)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res).Should(Equal(source))
				})
			})

			Context("when chunk = 5", func() {
				It("works without errors", func() {
					f := flusher.NewFlusher(5, m)

					m.EXPECT().
						AddOffers(gomock.Any()).
						Times(5)

					res, err := f.Flush(source)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res).Should(Equal(source))
				})
			})
		})

	})

})
