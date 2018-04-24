package server_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/contiamo/goserver/http"
)

var _ = Describe("Tracing", func() {
	It("should be possible to setup tracing", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go runMockTracingServer(ctx, ":1234")
		srv, err := createServer([]Option{WithTracing("localhost:1234", "test")})
		Expect(err).NotTo(HaveOccurred())
		go ListenAndServe(ctx, ":4004", srv)
		_, err = http.Get("http://localhost:4004")
		Expect(err).NotTo(HaveOccurred())
		time.Sleep(1 * time.Second) // wait for span to be transmitted
		Expect(receivedSomething).To(BeTrue())
	})
})

var receivedSomething bool

func runMockTracingServer(ctx context.Context, addr string) error {
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	listener, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			{
				return ctx.Err()
			}
		default:
			{
				_, _, err := listener.ReadFromUDP(buf[:])
				receivedSomething = true
				if err != nil {
					fmt.Println("Error: ", err)
				}
			}
		}
	}
}
