package http

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

var _ = Describe("Tracing", func() {
	tracer := mocktracer.New()
	opentracing.SetGlobalTracer(tracer)

	It("should be possible to setup tracing", func() {
		srv, err := createServer([]Option{WithTracing("localhost:test", "test", nil, nil)})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		_, err = http.Get(ts.URL + "/tracing/")
		Expect(err).NotTo(HaveOccurred())

		Expect(len(tracer.FinishedSpans())).To(Equal(1))
		tracer.Reset()
	})

	It("should be possible to set additional tags", func() {
		tagName := "testTag"
		tagValue := "something to find"
		srv, err := createServer([]Option{WithTracing("localhost:test", "test", map[string]string{tagName: tagValue}, nil)})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		_, err = http.Get(ts.URL + "/tracing/")
		Expect(err).NotTo(HaveOccurred())

		Expect(len(tracer.FinishedSpans())).To(Equal(1))

		span := tracer.FinishedSpans()[0]
		Expect(span.Tags()[tagName]).To(Equal(tagValue))
		tracer.Reset()
	})

	It("should replace uuid values with *", func() {
		srv, err := createServer([]Option{WithTracing("localhost:test", "test", nil, nil)})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		_, err = http.Get(ts.URL + "/tracing/2f6f97f2-5f44-476d-bc0c-180b2eaa36ca/2f6f97f2-5f44-476d-bc0c-180b2eaa36cb")
		Expect(err).NotTo(HaveOccurred())

		Expect(len(tracer.FinishedSpans())).To(Equal(1))
		span := tracer.FinishedSpans()[0]
		Expect(span.OperationName).To(Equal("HTTP GET /tracing/*/*"))
		tracer.Reset()
	})

	It("should allow websockets", func() {
		srv, err := createServer([]Option{WithTracing("localhost:test", "test", nil, nil)})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		err = testWebsocketEcho(ts.URL)
		Expect(err).NotTo(HaveOccurred())
	})
})
