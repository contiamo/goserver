package http

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CORS", func() {
	It("should be possible to configure custom cors rules", func() {
		allowedOrigins := []string{"foo.bar"}
		allowedMethods := []string{"HEAD"}
		allowedHeaders := []string{"Content-Type"}
		allowCredentials := true
		srv, err := createServer([]Option{WithCORS(allowedOrigins, allowedMethods, allowedHeaders, allowCredentials)})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		req, _ := http.NewRequest(http.MethodOptions, ts.URL+"/cors", nil)
		req.Header.Set("Access-Control-Request-Method", "HEAD")
		req.Header.Set("Origin", "foo.bar")
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal("true"))
		Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("foo.bar"))
		Expect(resp.Header.Get("Access-Control-Allow-Methods")).To(Equal("HEAD"))
	})

	It("should support websockets", func() {
		srv, err := createServer([]Option{WithCORSWideOpen()})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		err = testWebsocketEcho(ts.URL)
		Expect(err).NotTo(HaveOccurred())
	})

})
