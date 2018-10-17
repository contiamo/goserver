// Temporarily copied with permission from github.com/trusch/streamstore

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"
)

// Encrypt AES encrypts the msg using the provided key
func Encrypt(msg, key string) (string, error) {
	buf := &bytes.Buffer{}
	encoder, err := NewWriter(buf, key)
	if err != nil {
		return "", err
	}
	in := strings.NewReader(msg)
	if _, err = io.Copy(encoder, in); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// NewWriter returns a new aes writer
func NewWriter(base io.Writer, key string) (io.WriteCloser, error) {
	k := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	_, err = base.Write(iv)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewOFB(block, iv[:])
	writer := &cipher.StreamWriter{S: stream, W: base}
	return writer, nil
}
