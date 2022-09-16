package timer_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/steeringwaves/go-timer"
)

var _ = Describe("Timer", func() {
	Context("Timer", func() {
		It("Should send to channel", func() {
			t := timer.NewTimer(500 * time.Millisecond)
			now := time.Now()
			when := <-t.C
			t.Stop()

			now.Sub(when)
			Expect(when.Sub(now)).Should(BeNumerically("~", 500*time.Millisecond, 100*time.Millisecond))
		})
	})
})

func TestTimers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timer")
}
