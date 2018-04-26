package server_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	. "github.com/contiamo/goserver/http"
)

var _ = Describe("Tracing", func() {
	tracer := mocktracer.New()
	opentracing.SetGlobalTracer(tracer)

	It("should be possible to setup tracing", func() {
		srv, err := createServer([]Option{WithTracing("localhost:test", "test")})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		_, err = http.Get(ts.URL + "/tracing")
		Expect(err).NotTo(HaveOccurred())

		Expect(len(tracer.FinishedSpans())).To(Equal(1))
	})

	It("should support websockets", func() {
		srv, err := createServer([]Option{WithTracing("localhost:test", "test")})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()

		err = testWebsocketEcho(ts.URL)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(tracer.FinishedSpans())).To(Equal(1))

	})
})
