package notifier_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-offer-api/internal/notifier"
)

var _ = Describe("Notifier", func() {

	const (
		duration time.Duration = 100 * time.Millisecond
	)

	Context("creating new saver", func() {
		When("invalid data", func() {
			It("duration = 0", func() {
				n, err := notifier.NewNotifier(0)
				Expect(n).Should(BeNil())
				Expect(err).Should(BeEquivalentTo(notifier.ErrorDurationLessOrEqualZero))
			})
		})

		When("correct data", func() {
			It("successful creating", func() {
				n, err := notifier.NewNotifier(duration)
				Expect(n).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("initialization", func() {
		When("correct data", func() {
			It("successful initialization", func() {
				n, err := notifier.NewNotifier(duration)
				Expect(n).ShouldNot(BeNil())
				Expect(err).Should(BeNil())

				Expect(n.Init()).Should(BeNil())
				n.Close()
			})
		})

		When("duplicate initialization returns error", func() {
			It("already initialized", func() {
				n, err := notifier.NewNotifier(duration)
				Expect(n).ShouldNot(BeNil())
				Expect(err).Should(BeNil())

				// first initialization
				Expect(n.Init()).Should(BeNil())

				// second initialization
				Expect(n.Init()).Should(BeEquivalentTo(notifier.ErrorAlreadyInitialized))
			})
		})
	})
})
