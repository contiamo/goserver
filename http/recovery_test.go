package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Recovery", func() {
	It("should be possible to configure panic recovery", func() {
		logrus.SetOutput(ioutil.Discard)
		defer func() {
			logrus.SetOutput(os.Stdout)
		}()

		srv, err := createServer([]Option{WithRecovery(ioutil.Discard, true)})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		resp, err := http.Get(ts.URL + "/panic")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("should support websockets and tracing", func() {
		logrus.SetOutput(ioutil.Discard)
		srv, err := createServer([]Option{WithRecovery(ioutil.Discard, true)})
		Expect(err).NotTo(HaveOccurred())
		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()

		err = testWebsocketEcho(ts.URL)
		Expect(err).NotTo(HaveOccurred())
	})
})
