package server_test

import (
	"io"
	"net/http"
	"testing"

	httpserver "github.com/contiamo/goserver/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Http Suite")
}

func createServer(opts []httpserver.Option) (*http.Server, error) {
	return httpserver.New(&httpserver.Config{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/panic" {
				panic("PANIC!!!")
			}
			io.Copy(w, r.Body)
		}),
		Options: opts,
	})
}
