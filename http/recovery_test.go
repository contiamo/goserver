package server_test

import (
	"context"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	. "github.com/contiamo/goserver/http"
)

var _ = Describe("Recovery", func() {
	It("should be possible to configure panic recovery", func() {
		logrus.SetOutput(ioutil.Discard)
		srv, err := createServer([]Option{WithRecovery(ioutil.Discard, true)})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":4003", srv)
		resp, err := http.Get("http://localhost:4003/panic")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})
})
