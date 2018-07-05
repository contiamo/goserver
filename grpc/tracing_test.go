package grpc

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/contiamo/goserver/grpc/test"
)

var _ = Describe("Tracing", func() {
	tracer := mocktracer.New()
	opentracing.SetGlobalTracer(tracer)

	It("should be possible to setup tracing", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv, err := createServerWithOptions([]Option{WithTracing("localhost:test", "test")})
		Expect(err).NotTo(HaveOccurred())

		go ListenAndServe(ctx, ":1234", srv)
		cli, err := createPlaintextTestClient(ctx, "localhost:1234")
		Expect(err).NotTo(HaveOccurred())
		_, err = cli.Ping(ctx, &test.PingReq{})
		Expect(err).NotTo(HaveOccurred())

		Expect(len(tracer.FinishedSpans())).To(Equal(1))
	})
})
