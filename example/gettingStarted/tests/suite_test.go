package testing

import (
	"context"
	"os"
	"testing"

	"github.com/gothunder/thunder/pkg/events"
	"github.com/gothunder/thunder/pkg/events/mocks"
	"github.com/gothunder/thunder/pkg/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var handler *mocks.Handler
var publisher events.EventPublisher
var app *fx.App
var logger *zerolog.Logger

func TestCase(t *testing.T) {
	os.Setenv("TZ", "UTC")

	RegisterFailHandler(Fail)
	RunSpecs(t, "Payments Suite")
}

var _ = BeforeSuite(func() {
	app = fx.New(
		fx.Populate(&publisher, &logger),
		fx.Provide(
			func() events.Handler {
				handler = mocks.NewHandler()
				return handler
			},
		),
		log.Module,
		mocks.Module,
	)

	Expect(app.Start(context.Background())).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(app.Stop(context.Background())).To(Succeed())
})

var _ = BeforeEach(func() {
	t := GinkgoT()
	handler.Mock.Test(t)

	DeferCleanup(func() {
		Expect(handler.Mock.AssertExpectations(t)).To(BeTrue())
		handler.ResetMock()
	})
})
