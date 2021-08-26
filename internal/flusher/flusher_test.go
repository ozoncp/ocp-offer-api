package flusher_test

import (
	"context"
	"errors"
	"reflect"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-offer-api/internal/flusher"
	"github.com/ozoncp/ocp-offer-api/internal/mocks"
	"github.com/ozoncp/ocp-offer-api/internal/models"
	utils "github.com/ozoncp/ocp-offer-api/internal/utils/models"
)

var _ = Describe("Flusher", func() {
	var (
		ctrl   *gomock.Controller
		m      *mocks.MockIRepository
		f      flusher.Flusher
		source []models.Offer
		ctx    context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		m = mocks.NewMockIRepository(ctrl)
		ctx = context.Background()
		source = []models.Offer{
			{ID: 10, UserID: 20, Grade: 30, TeamID: 40},
			{ID: 11, UserID: 21, Grade: 31, TeamID: 41},
			{ID: 12, UserID: 22, Grade: 32, TeamID: 42},
			{ID: 13, UserID: 23, Grade: 33, TeamID: 43},
			{ID: 14, UserID: 24, Grade: 34, TeamID: 44},
			{ID: 15, UserID: 25, Grade: 35, TeamID: 45},
			{ID: 16, UserID: 26, Grade: 36, TeamID: 46},
			{ID: 17, UserID: 27, Grade: 37, TeamID: 47},
			{ID: 18, UserID: 28, Grade: 38, TeamID: 48},
			{ID: 19, UserID: 29, Grade: 39, TeamID: 49},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("save offers to repo with flusher", func() {
		Context("when MultiCreateOffer returns error", func() {
			It("returns error", func() {
				f = flusher.NewFlusher(3, m)

				m.EXPECT().
					MultiCreateOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(uint64(0), errors.New("error"))

				res, err := f.Flush(ctx, source)

				Expect(err).Should(BeEquivalentTo(errors.New("error")))
				Expect(res).Should(Equal(source))
			})
		})

		Context("when MultiCreateOffer returns an error in loop of Flush function", func() {
			It("returns error", func() {
				chunkSize := 3
				f = flusher.NewFlusher(chunkSize, m)
				chunks, _ := utils.SplitOffersToBatches(source, uint(chunkSize))

				m.EXPECT().
					MultiCreateOffer(gomock.Any(), gomock.Any()).
					Times(2).
					DoAndReturn(
						func(ctx context.Context, chunk []models.Offer) (uint64, error) {
							if reflect.DeepEqual(chunk, chunks[1]) {
								return 0, errors.New("error")
							}

							return 0, nil
						},
					)

				res, err := f.Flush(ctx, source)

				// проверяем что вернутась нужная ошибка
				Expect(err).Should(BeEquivalentTo(errors.New("error")))

				// проверяем что вернулась часть - только не добавленные в хранилище
				Expect(res).Should(Equal(source[len(chunks[1]):]))
			})
		})

		Context("not normal case", func() {
			Context("when chunk < 0", func() {
				It("returns error", func() {
					f = flusher.NewFlusher(-1, m)

					m.EXPECT().
						MultiCreateOffer(gomock.Any(), gomock.Any()).
						Times(0)

					res, err := f.Flush(ctx, source)

					Expect(err).Should(HaveOccurred())
					Expect(res).Should(Equal(source))
				})
			})

			Context("when source is empty", func() {
				It("returns error", func() {
					data := make([]models.Offer, 0)
					f := flusher.NewFlusher(1, m)

					m.EXPECT().
						MultiCreateOffer(gomock.Any(), gomock.Any()).
						Times(0)

					res, err := f.Flush(ctx, data)

					Expect(err).Should(HaveOccurred())
					Expect(res).Should(Equal(data))
				})
			})

		})

		Context("normal case", func() {

			Context("when chunk = 1", func() {
				It("works without errors", func() {
					f := flusher.NewFlusher(1, m)

					m.EXPECT().
						MultiCreateOffer(gomock.Any(), gomock.Any()).
						Times(1)

					res, err := f.Flush(ctx, source)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res).Should(BeNil())
				})
			})

			Context("when chunk = 2", func() {
				It("works without errors", func() {
					f := flusher.NewFlusher(2, m)

					m.EXPECT().
						MultiCreateOffer(gomock.Any(), gomock.Any()).
						Times(2)

					res, err := f.Flush(ctx, source)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res).Should(BeNil())
				})
			})

			Context("when chunk = 5", func() {
				It("works without errors", func() {
					f := flusher.NewFlusher(5, m)

					m.EXPECT().
						MultiCreateOffer(gomock.Any(), gomock.Any()).
						Times(5)

					res, err := f.Flush(ctx, source)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res).Should(BeNil())
				})
			})
		})

	})

})
