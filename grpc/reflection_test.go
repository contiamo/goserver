package grpc

import (
	"context"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reflection", func() {
	It("should be possible to enable reflection", func() {
		srv, err := createServerWithOptions([]Option{WithReflection()})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ListenAndServe(ctx, ":3006", srv)
		// it takes some time to run the server, can't be accessed immediately
		time.Sleep(100 * time.Millisecond)

		cmd := exec.Command("grpcurl", "-plaintext", "localhost:3006", "describe")
		Expect(cmd.Run()).To(Succeed())
	})
})
