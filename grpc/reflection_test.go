package server_test

import (
	"context"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/contiamo/goserver/grpc"
)

var _ = Describe("Reflection", func() {
	It("should be possible to enable reflection", func() {
		srv, err := createServerWithOptions([]Option{WithReflection()})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":3005", srv)
		cmd := exec.Command("grpcurl", "-plaintext", "localhost:3005", "describe")
		Expect(cmd.Run()).To(Succeed())
	})
})