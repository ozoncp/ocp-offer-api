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
				f := flusher.NewFlusher(chunkSize, m)
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
