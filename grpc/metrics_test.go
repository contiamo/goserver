package server_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/contiamo/goserver"
	. "github.com/contiamo/goserver/grpc"
	"github.com/contiamo/goserver/grpc/test"
)

var _ = Describe("Metrics", func() {
	It("should be possible to collect metrics", func() {
		srv, err := createServerWithOptions([]Option{WithMetrics()})
		Expect(err).NotTo(HaveOccurred())
		Expect(srv).NotTo(BeNil())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":3003", srv)
		go goserver.ListenAndServeMetricsAndHealth(":8080", nil)
		cli, err := createPlaintextTestClient(ctx, "localhost:3003")
		Expect(err).NotTo(HaveOccurred())
		_, err = cli.Ping(ctx, &test.PingReq{Msg: "test"})
		Expect(err).NotTo(HaveOccurred())
		resp, err := http.Get("http://localhost:8080/metrics")
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()
		bs, _ := ioutil.ReadAll(resp.Body)
		Expect(strings.Contains(string(bs), `grpc_server_handled_total{grpc_code="OK",grpc_method="Ping",grpc_service="test.PingPong",grpc_type="unary"} 1`))
	})
})
