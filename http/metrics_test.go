package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/contiamo/goserver"
)

var _ = Describe("Metrics", func() {
	It("should be possible to configure metrics collection", func() {
		srv, err := createServer([]Option{WithMetrics("test")})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ListenAndServe(ctx, ":4002", srv)
		// it takes some time to run the server, can't be accessed immediately
		time.Sleep(100 * time.Millisecond)

		_, err = http.Get("http://localhost:4002/metrics_test")
		Expect(err).NotTo(HaveOccurred())

		go goserver.ListenAndServeMetricsAndHealth(":8080", nil)
		// it takes some time to run the server, can't be accessed immediately
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://localhost:8080/metrics")
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()
		bs, _ := ioutil.ReadAll(resp.Body)
		Expect(strings.Contains(string(bs), `latencies_bucket{code="200",method="get",service="test",le="+Inf"} 1`))
	})

	It("should support websockets", func() {
		srv, err := createServer([]Option{WithMetrics("test")})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()

		err = testWebsocketEcho(ts.URL)
		Expect(err).NotTo(HaveOccurred())

		go goserver.ListenAndServeMetricsAndHealth(":8080", nil)
		// it takes some time to run the server, can't be accessed immediately
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://localhost:8080/metrics")
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()
		bs, _ := ioutil.ReadAll(resp.Body)
		Expect(strings.Contains(string(bs), `latencies_bucket{code="200",method="get",service="test",le="+Inf"} 1`))
	})
})
