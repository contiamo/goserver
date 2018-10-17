// Temporarily copied with permission from github.com/trusch/streamstore

package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"io"
)

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
