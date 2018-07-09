package grpc

import (
	"context"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/contiamo/goserver/grpc/test"
	utils "github.com/contiamo/goserver/test"
)

var _ = Describe("Recovery", func() {
	It("should recover from panic when recovery option is set", func() {
		buf, restore := utils.SetupLoggingBuffer()
		defer restore()

		srv, err := createServerWithOptions([]Option{WithRecovery()})
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
		Expect(strings.Contains(buf.String(), `level=error msg="Very bad panic" stacktrace=`)).To(BeTrue(), "The logs should contain the error message and stacktrace but got %s", buf.String())
	})
})
