package server_test

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	. "github.com/contiamo/goserver/http"
)

var _ = Describe("Logging", func() {
	It("should be possible to configure logging", func() {
		buf := &bytes.Buffer{}
		logrus.SetOutput(buf)
		srv, err := createServer([]Option{WithLogging("test")})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":4001", srv)
		_, err = http.Get("http://localhost:4001")
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Contains(buf.String(), "completed handling request")).To(BeTrue())
	})
})
