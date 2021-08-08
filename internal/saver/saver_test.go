package saver_test

import (
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-offer-api/internal/flusher"
	"github.com/ozoncp/ocp-offer-api/internal/mocks"
	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/ozoncp/ocp-offer-api/internal/notifier"
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
		m    *mocks.MockRepo
		n    notifier.Notifier
		f    flusher.Flusher
		s    saver.Saver
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		m = mocks.NewMockRepo(ctrl)
		f = flusher.NewFlusher(int(chunkSize), m)
		n, _ = notifier.NewNotifier(duration)
		s, _ = saver.NewSaver(capacity, f, n)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("creating new saver", func() {

		When("invalid data", func() {
			It("capacity = 0", func() {
				s, err := saver.NewSaver(0, f, n)

				Expect(s).Should(BeNil())
				Expect(err).Should(BeEquivalentTo(saver.ErrorCapacionLessOrEqualZero))
			})

			It("flusher is nil", func() {
				s, err := saver.NewSaver(capacity, nil, n)

				Expect(s).Should(BeNil())
				Expect(err).Should(BeEquivalentTo(saver.ErrorFlusherIsNil))
			})

			It("notifier is nil", func() {
				s, err := saver.NewSaver(capacity, f, nil)

				Expect(s).Should(BeNil())
				Expect(err).Should(BeEquivalentTo(saver.ErrorNotifierIsNil))
			})

			It("all props incorrected", func() {
				s, err := saver.NewSaver(0, nil, nil)

				Expect(s).Should(BeNil())
				Expect(err).Should(BeEquivalentTo(saver.ErrorCapacionLessOrEqualZero))
			})
		})

		When("correct data", func() {
			It("all props corrected", func() {
				s, err := saver.NewSaver(capacity, f, n)

				Expect(s).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("initialization", func() {
		When("correct data", func() {
			It("successful initialization", func() {
				Expect(s.Init()).Should(BeNil())
				s.Close()
			})
		})

		When("duplicate initialization returns error", func() {
			It("already initialized", func() {
				// first initialization
				Expect(s.Init()).Should(BeNil())

				// second initialization
				Expect(s.Init()).
					Should(BeEquivalentTo(saver.ErrorAlreadyInitialized))
			})
		})
	})

	Context("invalid case", func() {
		When("save returns an error", func() {
			It("saver not initializated", func() {
				Expect(s.Save(models.Offer{Id: 1})).
					Should(BeEquivalentTo(saver.ErrorNotInitialized))
			})

			It("when there are more offers than capacity", func() {
				Expect(s.Init()).Should(BeNil())

				countOffers := (int(capacity) + 1) * 2

				m.EXPECT().
					AddOffers(gomock.Any()).
					AnyTimes().Return(errors.New("error"))

				for i := 0; i < countOffers; i++ {
					err := s.Save(models.Offer{
						Id: uint64(i),
					})

					if uint(i) > capacity {
						Expect(err).Should(BeEquivalentTo(saver.ErrorMaximumCapacityReached))
					} else if uint(i) < capacity {
						Expect(err).Should(BeNil())
					}
				}
				s.Close()
			})
		})
	})

	Context("normal case", func() {
		It("when offers are less than capacity", func() {
			err := s.Init()
			Expect(err).Should(BeNil())

			countOffers := int(capacity) - 1

			m.EXPECT().
				AddOffers(gomock.Any()).
				Times(int(chunkSize)).
				Return(nil)

			for i := 0; i < countOffers; i++ {
				err := s.Save(models.Offer{
					Id: uint64(i) + 1,
				})
				Expect(err).Should(BeNil())
			}

			s.Close()
		})
	})

})
