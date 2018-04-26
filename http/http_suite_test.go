package server_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/websocket"

	httpserver "github.com/contiamo/goserver/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Http Suite")
}

var upgrader = websocket.Upgrader{}

func echoWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func createServer(opts []httpserver.Option) (*http.Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws/", echoWS)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/panic" {
			panic("PANIC!!!")
		}
		io.Copy(w, r.Body)
	})

	return httpserver.New(&httpserver.Config{
		Handler: mux,
		Options: opts,
	})
}

func testWebsocketEcho(server string) error {
	u := "ws" + strings.TrimPrefix(server, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u+"/ws/echo", nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Send message to server, read response and check to see if it's what we expect.
	err = ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	if err != nil {
		return err
	}

	_, p, err := ws.ReadMessage()
	if err != nil {
		return err
	}

	if string(p) != "hello" {
		return fmt.Errorf("websocket echo expected \"hello\" but got \"%s\"", string(p))
	}

	return nil
}
