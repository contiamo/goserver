package server

import (
	"net/http"

	"github.com/rs/cors"
)

// WithCORS configures CORS on the webserver
func WithCORS(allowedOrigins, allowedMethods, allowedHeaders []string, allowCredentials bool) Option {
	return &corsOption{allowedOrigins, allowedMethods, allowedHeaders, allowCredentials}
}

// WithCORSWideOpen allows requests from all origins with all methods and all headers/cookies/credentials allowed.
func WithCORSWideOpen() Option {
	return &corsOption{
		allowedOrigins:   []string{"*"},
		allowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		allowedHeaders:   []string{"*"},
		allowCredentials: true,
	}
}

type corsOption struct {
	allowedOrigins   []string
	allowedMethods   []string
	allowedHeaders   []string
	allowCredentials bool
}

func (opt *corsOption) WrapHandler(handler http.Handler) (http.Handler, error) {
	c := cors.New(cors.Options{
		AllowedOrigins:   opt.allowedOrigins,
		AllowedMethods:   opt.allowedMethods,
		AllowedHeaders:   opt.allowedHeaders,
		AllowCredentials: opt.allowCredentials,
	})
	return c.Handler(handler), nil
}
