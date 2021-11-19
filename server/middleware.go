package server

import (
	"net/http"

	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/urfave/negroni"
)

// NewRequestLogger creates a RequestLogger.
func NewRequestLogger(s *Server) *RequestLogger {
	return &RequestLogger{
		s.Logger,
		s.WebAuthentication,
		s.WebUsername,
		s.WebPassword,
	}
}

// RequestLogger logs requests and their response codes
// as well as handle the basicauth on the requests
type RequestLogger struct {
	logger            logging.SimpleLogging
	WebAuthentication bool
	WebUsername       string
	WebPassword       string
}

// ServeHTTP implements the middleware function. It logs all requests at DEBUG level.
func (l *RequestLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	l.logger.Debug("%s %s – from %s", r.Method, r.URL.RequestURI(), r.RemoteAddr)
	allowed := false
	if !l.WebAuthentication ||
		r.URL.Path == "/events" ||
		r.URL.Path == "/healthz" ||
		r.URL.Path == "/status" {
		allowed = true
	} else {
		user, pass, ok := r.BasicAuth()
		if ok {
			r.SetBasicAuth(user, pass)
			l.logger.Debug("user: %s / pass: %s >> url: %s", user, pass, r.URL.RequestURI())
			if user == l.WebUsername && pass == l.WebPassword {
				l.logger.Debug("[VALID] user: %s / pass: %s >> url: %s", user, pass, r.URL.RequestURI())
				allowed = true
			} else {
				allowed = false
				l.logger.Info("[INVALID] user: %s / pass: %s >> url: %s", user, pass, r.URL.RequestURI())
			}
		}
	}
	if !allowed {
		rw.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
	} else {
		next(rw, r)
	}
	l.logger.Debug("%s %s – respond HTTP %d", r.Method, r.URL.RequestURI(), rw.(negroni.ResponseWriter).Status())
}
