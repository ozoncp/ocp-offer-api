package saver_test

import (
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-offer-api/internal/flusher"
	"github.com/ozoncp/ocp-offer-api/internal/mocks"
	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/ozoncp/ocp-offer-api/internal/saver"
)

var _ = Describe("Saver", func() {

	const (
		chunkSize uint          = 2
		capacity  uint          = 10
		duration  time.Duration = 100 * time.Millisecond
	)

	var (
		ctrl *gomock.Controller
		m    *mocks.MockIRepository
		f    flusher.Flusher
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		m = mocks.NewMockIRepository(ctrl)
		f = flusher.NewFlusher(int(chunkSize), m)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("creating new saver", func() {

		When("invalid data", func() {
			It("capacity = 0", func() {
				s, err := saver.NewSaver(0, f, duration)

				Expect(s).Should(BeNil())
				Expect(err).Should(BeEquivalentTo(saver.ErrorCapacionLessOrEqualZero))
			})

			It("flusher is nil", func() {
				s, err := saver.NewSaver(capacity, nil, duration)

				Expect(s).Should(BeNil())
				Expect(err).Should(BeEquivalentTo(saver.ErrorFlusherIsNil))
			})

			It("duration = 0", func() {
				s, err := saver.NewSaver(capacity, f, 0)

				Expect(s).Should(BeNil())
				Expect(err).Should(BeEquivalentTo(saver.ErrorDurationLessOrEqualZero))
			})

			It("all props incorrected", func() {
				s, err := saver.NewSaver(0, nil, 0)

				Expect(s).Should(BeNil())
				Expect(err).Should(BeEquivalentTo(saver.ErrorCapacionLessOrEqualZero))
			})
		})

		When("correct data", func() {
			It("all props corrected", func() {
				s, err := saver.NewSaver(capacity, f, duration)

				Expect(s).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("invalid case", func() {
		When("save returns an error", func() {
			It("when save channel is closed ", func() {
				s, _ := saver.NewSaver(capacity, f, duration)

				m.EXPECT().
					MultiCreateOffer(gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(uint64(1), nil)

				s.Close()
				err := s.Save(models.Offer{ID: 1})

				Expect(err).Should(BeEquivalentTo(saver.ErrorChanelClosed))
			})
		})
	})

	Context("normal case", func() {
		It("when offers are less than capacity", func() {
			s, _ := saver.NewSaver(capacity, f, duration)

			m.EXPECT().
				MultiCreateOffer(gomock.Any(), gomock.Any()).
				AnyTimes().
				Return(uint64(0), nil)

			countOffers := int(capacity) - 1

			for i := 0; i < countOffers; i++ {
				err := s.Save(models.Offer{
					ID: uint64(i) + 1,
				})
				Expect(err).Should(BeNil())
			}

			s.Close()
		})
	})

})
