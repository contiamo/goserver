// Copyright 2017 David Ackroyd. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/trusch/streamstore/filter/encryption/aes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func WithErrorEncryption(key string) Option {
	return &errorEncryptionOption{key}
}

type errorEncryptionOption struct {
	key string
}

func (opt *errorEncryptionOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	si := errorEncrypterStreamServerInterceptor(opt.key)
	ui := errorEncrypterUnaryServerInterceptor(opt.key)
	return nil, si, ui, nil
}

func (opt *errorEncryptionOption) PostProcess(s *grpc.Server) error {
	return nil
}

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func errorEncrypterUnaryServerInterceptor(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if err != nil {
				stat, _ := status.FromError(err)
				msg, e := encrypt(stat.Message(), key)
				fmt.Println(msg)
				if e != nil {
					logrus.Errorf("failed to encrypt error message: %v", e)
					return
				}
				fmt.Println("before: ", err)
				err = status.Error(stat.Code(), msg)
				fmt.Println("after: ", err)
			}
		}()
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for panic recovery.
func errorEncrypterStreamServerInterceptor(key string) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if err != nil {
				if stat, ok := status.FromError(err); ok {
					msg, e := encrypt(stat.Message(), key)
					if e != nil {
						logrus.Errorf("failed to encrypt error message: %v", e)
						return
					}
					err = status.Error(stat.Code(), msg)
				}
			}
		}()
		return handler(srv, stream)
	}
}

func encrypt(msg, key string) (string, error) {
	buf := &bytes.Buffer{}
	encoder, err := aes.NewWriter(buf, key)
	if err != nil {
		return "", err
	}
	in := strings.NewReader(msg)
	if _, err = io.Copy(encoder, in); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
