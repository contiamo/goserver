package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/trusch/streamstore/filter/encryption/aes"
)

var msg = flag.String("msg", "", "base64 encoded message")
var key = flag.String("key", "", "aes key")

func main() {
	flag.Parse()
	bs, err := base64.StdEncoding.DecodeString(*msg)
	if err != nil {
		logrus.Fatal(err)
	}
	decrypter, err := aes.NewReader(bytes.NewReader(bs), *key)
	if err != nil {
		logrus.Fatal(err)
	}
	io.Copy(os.Stdout, decrypter)
}
