package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
	logrusmiddleware "github.com/trusch/logrus-middleware"
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
