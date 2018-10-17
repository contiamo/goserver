package grpc

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/contiamo/goserver/grpc/test"
)

var _ = Describe("Error Encryption", func() {
	It("should recover from panic when recovery option is set and encrypt the error message", func() {
		srv, err := createServerWithOptions([]Option{WithErrorEncryption("foobar"), WithRecovery()})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ListenAndServe(ctx, ":3005", srv)
		// it takes some time to run the server, can't be accessed immediately
		time.Sleep(100 * time.Millisecond)

		cli, err := createPlaintextTestClient(ctx, "localhost:3005")
		Expect(err).NotTo(HaveOccurred())
		_, err = cli.Panic(ctx, &test.PingReq{"Very bad panic"})

		Expect(err).To(HaveOccurred())
		Expect(grpc.Code(err)).To(Equal(codes.Internal))
	})
})
