package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Logging", func() {

	It("should be possible to configure logging", func() {
		buf := &bytes.Buffer{}
		logrus.SetOutput(buf)
		srv, err := createServer([]Option{WithLogging("test")})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		_, err = http.Get(ts.URL + "/logging")
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Contains(buf.String(), "successfully handled request")).To(BeTrue())
	})

	It("should support websockets", func() {
		buf := &Buffer{}
		logrus.SetOutput(buf)
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

type Buffer struct {
	data []byte
}

func (b *Buffer) Write(data []byte) (int, error) {
	b.data = append(b.data, data...)
	return len(data), nil
}

func (b *Buffer) String() string {
	return string(b.data)
}
