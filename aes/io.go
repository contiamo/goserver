// Temporarily copied with permission from github.com/trusch/streamstore

package aes

import (
	"errors"
	"io"
)

// IOCoppler copples reader, writer and closer
type IOCoppler struct {
	upper interface{}
	lower interface{}
}

// NewIOCoppler returns a new IOCoppler instance
func NewIOCoppler(upper, lower interface{}) *IOCoppler {
	return &IOCoppler{upper, lower}
}

// Write writes data to the upper layer
func (coppler *IOCoppler) Write(data []byte) (int, error) {
	if writer, ok := coppler.upper.(io.Writer); ok {
		return writer.Write(data)
	}
	return -1, errors.New("Write() not supported")
}

// Read reads data from the upper layer
func (coppler *IOCoppler) Read(data []byte) (int, error) {
	if reader, ok := coppler.upper.(io.Reader); ok {
		return reader.Read(data)
	}
	return -1, errors.New("Read() not supported")
}

// Close closes the upper layer first, then the lower layer
func (coppler *IOCoppler) Close() error {
	if closer, ok := coppler.upper.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			return err
		}
	}
	if closer, ok := coppler.lower.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
