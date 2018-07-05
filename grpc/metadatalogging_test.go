package grpc

import (
	"bytes"
	"context"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/contiamo/goserver/grpc/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sirupsen/logrus"
)

var _ = Describe("Logging", func() {
	It("should be possible to setup logging option", func() {
		buf := &bytes.Buffer{}
		level := logrus.GetLevel()
		logrus.SetFormatter(&logrus.TextFormatter{DisableColors: false})
		logrus.SetOutput(buf)
		logrus.SetLevel(logrus.DebugLevel)

		defer func() {
			logrus.SetOutput(os.Stdout)
			logrus.SetLevel(level)
		}()
		srv, err := createServerWithOptions([]Option{WithMDLogging()})
		Expect(err).NotTo(HaveOccurred())
		Expect(srv).NotTo(BeNil())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ListenAndServe(ctx, ":3003", srv)
		time.Sleep(time.Second)
		cli, err := createPlaintextTestClient(ctx, ":3003")
		Expect(err).NotTo(HaveOccurred())

		md := metadata.New(map[string]string{"test": "value"})
		ctx = metadata.NewOutgoingContext(ctx, md)
		resp, err := cli.Ping(ctx, &test.PingReq{Msg: "test"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Msg).To(Equal("test"))

		Expect(strings.Contains(buf.String(), "test=value")).To(BeTrue(), "This should be the value in logs but got %s", buf.String())
	})
})
