package grpc

import (
	"context"

	"github.com/contiamo/goserver/aes"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithErrorScrubbing ensures that internal or unknown errors do not leak information.
// You  may pass a function to customize the scrubbing behavior. The package provides
// two simple scrubbers: SimpleErrorScrubber and NoopErrorScrubber for convenience.
func WithErrorScrubbing(scrubber func(error) error) Option {
	if scrubber == nil {
		scrubber = SimpleErrorScrubber
	}
	return &errorScrubOption{scrubber}
}

type errorScrubOption struct {
	scrubber func(error) error
}

func (opt *errorScrubOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	scrubber := opt.scrubber
	if scrubber == nil {
		scrubber = SimpleErrorScrubber
	}

	si := StreamErrorInterceptor(scrubber)
	ui := UnaryErrorInterceptor(scrubber)
	return nil, si, ui, nil
}

func (opt *errorScrubOption) PostProcess(s *grpc.Server) error {
	return nil
}

// UnaryErrorInterceptor ensures internal or unknown errors do not leak information
func UnaryErrorInterceptor(scrubber func(error) error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)

		if err == nil {
			return resp, err
		}

		return resp, scrubber(err)
	}
}

// StreamErrorInterceptor ensure internal or unknown errors do not leak information
func StreamErrorInterceptor(scrubber func(error) error) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, stream)

		if err == nil {
			return err
		}
		return scrubber(err)
	}
}

// SimpleErrorScrubber replaces unknown or internal errors with a generic error object
func SimpleErrorScrubber(err error) error {
	switch status.Code(err) {
	case codes.Unknown, codes.Internal:
		st := status.New(codes.Internal, "Internal Server Error")
		return st.Err()
	default:
		return err
	}
}

// NewEncryptedErrorScrubber replaces unknown or internal errors with a generic error object. The
// original error message is encrypted and added to the error message in the status.Details
func NewEncryptedErrorScrubber(key string) func(error) error {
	return func(err error) error {
		switch status.Code(err) {
		case codes.Unknown, codes.Internal:

			st := status.New(codes.Internal, "Internal Server Error")

			msg, e := aes.Encrypt(err.Error(), key)
			if e != nil {
				return st.Err()
			}

			// can only error if st.Code == OK
			st, _ = st.WithDetails(&errdetails.DebugInfo{
				Detail: msg,
			})
			return st.Err()
		default:
			return err
		}
	}
}

// NoopErrorScrubber returns the error unmodified, this can be useful in debug or development
// environments
func NoopErrorScrubber(err error) error {
	return err
}
