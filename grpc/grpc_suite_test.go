package grpc

import (
	"context"
	"crypto/tls"
	"testing"

	"github.com/contiamo/goserver/grpc/test"
	utils "github.com/contiamo/goserver/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func TestGrpc(t *testing.T) {
	RegisterFailHandler(Fail)
	restore := utils.DiscardLogging()
	defer restore()
	RunSpecs(t, "Grpc Suite")
}

func createServerWithOptions(opts []Option) (*grpc.Server, error) {
	return New(&Config{
		Options: opts,
		Register: func(srv *grpc.Server) {
			test.RegisterPingPongServer(srv, test.NewPingPongServer())
		},
	})
}

func createTestClient(ctx context.Context, addr string) (test.PingPongClient, error) {
	crt, _ := tls.LoadX509KeyPair("test/pki/test-client.crt", "test/pki/test-client.key")
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{crt},
	})
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	return test.NewPingPongClient(conn), nil
}

func createPlaintextTestClient(ctx context.Context, addr string) (test.PingPongClient, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return test.NewPingPongClient(conn), nil
}
