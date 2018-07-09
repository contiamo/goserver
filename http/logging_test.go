package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	utils "github.com/contiamo/goserver/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {

	It("should be possible to configure logging", func() {
		buf, restore := utils.SetupLoggingBuffer()
		defer restore()

		srv, err := createServer([]Option{WithLogging("test")})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		_, err = http.Get(ts.URL + "/logging")
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Contains(buf.String(), "successfully handled request")).To(BeTrue())
	})

	It("should support websockets", func() {
		buf, restore := utils.SetupLoggingBuffer()
		defer restore()

		srv, err := createServer([]Option{WithLogging("test")})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		err = testWebsocketEcho(ts.URL)
		Expect(err).NotTo(HaveOccurred())
		time.Sleep(100 * time.Millisecond)
		Expect(strings.Contains(buf.String(), "successfully handled request")).To(BeTrue())
	})
})
