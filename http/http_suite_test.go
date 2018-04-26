package server_test

import (
	"io"
	"log"
	"net/http"
	"testing"

	"golang.org/x/net/websocket"

	httpserver "github.com/contiamo/goserver/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Http Suite")
}

func createServer(opts []httpserver.Option) (*http.Server, error) {
	mux := http.NewServeMux()
	mux.Handle("/ws/", websocket.Handler(func(ws *websocket.Conn) {
		io.Copy(ws, ws)
	}))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/panic" {
			panic("PANIC!!!")
		}
		log.Println("Got request to " + r.URL.Path)
		io.Copy(w, r.Body)
	})

	return httpserver.New(&httpserver.Config{
		Handler: mux,
		Options: opts,
	})
}
