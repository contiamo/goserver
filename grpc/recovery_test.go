package grpc

import (
	"context"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/contiamo/goserver/grpc/test"
)

var _ = Describe("Recovery", func() {
	It("should recover from panic when recovery option is set", func() {
		logrus.SetOutput(ioutil.Discard)
		defer func() {
			logrus.SetOutput(os.Stdout)
		}()
		srv, err := createServerWithOptions([]Option{WithRecovery()})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":3005", srv)
		cli, err := createPlaintextTestClient(ctx, "localhost:3005")
		Expect(err).NotTo(HaveOccurred())
		_, err = cli.Panic(ctx, &test.PingReq{})
		Expect(err).To(HaveOccurred())
		Expect(grpc.Code(err)).To(Equal(codes.Internal))
	})
})
