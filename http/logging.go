package server

import (
	"net/http"

	logrusmiddleware "github.com/bakins/logrus-middleware"
	"github.com/sirupsen/logrus"
)

// WithLogging configures a logrus middleware for that server
func WithLogging(app string) Option {
	return &loggingOption{app}
}

type loggingOption struct{ app string }

func (opt *loggingOption) WrapHandler(handler http.Handler) (http.Handler, error) {
	l := logrusmiddleware.Middleware{
		Name:   opt.app,
		Logger: logrus.StandardLogger(),
	}
	return l.Handler(handler, opt.app), nil
}
