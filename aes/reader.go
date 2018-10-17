// Temporarily copied with permission from github.com/trusch/streamstore

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

func Decrypt(msg, key string) (string, error) {
	msgBytes, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}

	decoder, err := NewReader(bytes.NewReader(msgBytes), key)
	if err != nil {
		return "", err
	}
	defer decoder.Close()

	buf := &bytes.Buffer{}
	if _, err = io.Copy(buf, decoder); err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

// NewReader returns a new aes reader
func NewReader(base io.Reader, key string) (io.ReadCloser, error) {
	k := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	bs, err := base.Read(iv[:])
	if bs != aes.BlockSize {
		return nil, errors.New("ciphertext to short")
	}
	if err != nil {
		return nil, err
	}
	stream := cipher.NewOFB(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: base}
	return NewIOCoppler(reader, base), nil
}
