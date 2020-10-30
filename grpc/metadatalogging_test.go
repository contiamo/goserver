package grpc

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/contiamo/goserver/grpc/test"
	utils "github.com/contiamo/goserver/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	It("should be possible to setup logging option", func() {
		buf, restore := utils.SetupLoggingBuffer()
		defer restore()

		srv, err := createServerWithOptions([]Option{WithMDLogging()})
		Expect(err).NotTo(HaveOccurred())
		Expect(srv).NotTo(BeNil())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ListenAndServe(ctx, ":3003", srv)
		// it takes some time to run the server, can't be accessed immediately
		time.Sleep(100 * time.Millisecond)

		cli, err := createPlaintextTestClient(ctx, ":3003")
		Expect(err).NotTo(HaveOccurred())

		md := metadata.New(map[string]string{"test": "value", "token": "foo", "password": "secret", "secret": "hide me"})
		ctx = metadata.NewOutgoingContext(ctx, md)
		resp, err := cli.Ping(ctx, &test.PingReq{Msg: "test"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Msg).To(Equal("test"))

		logLine := buf.String()
		Expect(strings.Contains(logLine, "test=value")).To(BeTrue(), "This should be the value in logs but got %s", logLine)
		Expect(strings.Contains(logLine, "token=\"****\"")).To(BeTrue(), "This should be the value in logs but got %s", logLine)
		Expect(strings.Contains(logLine, "password=\"****\"")).To(BeTrue(), "This should be the value in logs but got %s", logLine)
		Expect(strings.Contains(logLine, "secret=\"****\"")).To(BeTrue(), "This should be the value in logs but got %s", logLine)
	})
})
