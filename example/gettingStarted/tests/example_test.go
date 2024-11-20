package testing

import (
	"context"
	"encoding/json"

	"github.com/gothunder/thunder/pkg/events"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

type testEvent struct {
	Hello string `json:"hello"`
}

var _ = Describe("Example", func() {
	It("should pass", func() {
		Expect(true).To(BeTrue())
	})

	It("publishes an event", func() {
		done := make(chan struct{})

		topic := "test12345"

		handler.Mock.On("Handle",
			mock.Anything,
			topic,
			mock.Anything,
		).Return(events.Success).Run(func(args mock.Arguments) {
			defer GinkgoRecover()

			var event testEvent
			json.Unmarshal(args.Get(2).([]byte), &event)

			Expect(event.Hello).To(Equal("world"))

			close(done)
		})

		publisher.Publish(context.Background(), topic, testEvent{
			Hello: "world",
		})

		Eventually(done).Should(BeClosed())
	})
})
