package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// WithCredentials configures the server to use the given cert/key combination.
// If a ca file is supplied it is used to verify clients which now are required to have a certificate.
func WithCredentials(cert, key, ca string) Option {
	return &credentialOption{cert, key, ca}
}

type credentialOption struct {
	cert, key, ca string
}

func (opt *credentialOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	if opt.cert == "" && opt.key == "" {
		// not specifying keys is a noop
		return nil, nil, nil, nil
	}
	certificate, err := tls.LoadX509KeyPair(opt.cert, opt.key)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not load server key pair")
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	if opt.ca != "" {
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(opt.ca)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "could not read ca certificate")
		}
		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, nil, nil, errors.New("failed to append ca cert")
		}
		tlsConf.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConf.ClientCAs = certPool
	}
	creds := grpc.Creds(credentials.NewTLS(tlsConf))
	return creds, nil, nil, nil
}

func (opt *credentialOption) PostProcess(s *grpc.Server) error {
	return nil
}
