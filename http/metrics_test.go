package server_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/contiamo/goserver"
	. "github.com/contiamo/goserver/http"
)

var _ = Describe("Metrics", func() {
	It("should be possible to configure metrics collection", func() {
		srv, err := createServer([]Option{WithMetrics("test")})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":4002", srv)
		_, err = http.Get("http://localhost:4002/metrics_test")
		Expect(err).NotTo(HaveOccurred())
		go goserver.ListenAndServeMetricsAndHealth(":8080", nil)
		resp, err := http.Get("http://localhost:8080/metrics")
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()
		bs, _ := ioutil.ReadAll(resp.Body)
		Expect(strings.Contains(string(bs), `latencies_bucket{code="200",method="get",service="test",le="+Inf"} 1`))
	})
})
