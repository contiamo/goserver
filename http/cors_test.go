package server_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/contiamo/goserver/http"
)

var _ = Describe("CORS", func() {
	It("should be possible to configure custom cors rules", func() {
		allowedOrigins := []string{"foo.bar"}
		allowedMethods := []string{"HEAD"}
		allowedHeaders := []string{"Content-Type"}
		allowCredentials := true
		srv, err := createServer([]Option{WithCORS(allowedOrigins, allowedMethods, allowedHeaders, allowCredentials)})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":4004", srv)
		req, _ := http.NewRequest(http.MethodOptions, "http://localhost:4004", nil)
		req.Header.Set("Access-Control-Request-Method", "HEAD")
		req.Header.Set("Origin", "foo.bar")
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal("true"))
		Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("foo.bar"))
		Expect(resp.Header.Get("Access-Control-Allow-Methods")).To(Equal("HEAD"))
	})
})
