package server

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/utils"
	"net/http"
	"time"

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

var jwtMiddleware = jwt.New(jwt.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.RootSecret), nil
	},
	Expiration:    true,
	Extractor:     jwt.FromParameter("token"),
	SigningMethod: jwt.SigningMethodHS256,
})

func createToken(userId string) (string, error) {
	now := time.Now()
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"iat":    now.Unix(),
		"exp":    now.Add(15 * time.Minute).Unix(),
	})

	return token.SignedString([]byte(utils.RootSecret))
}

func checkToken(ctx iris.Context) {
	path := ctx.Path()
	if path == "/login" || path == "/healthz" || path == "/status" {
		ctx.Next()
	}

	if err := jwtMiddleware.CheckJWT(ctx); err != nil {
		jwtMiddleware.Config.ErrorHandler(ctx, err)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.Values().Set("msg", "Wrong token")
		return
	}

	token := ctx.Values().Get("jwt").(*jwt.Token)
	tokenInfo := token.Claims.(jwt.MapClaims)
	userId := ctx.URLParam("id")
	if userId != tokenInfo["userId"] {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.Values().Set("msg", "User does not have permission.")
		return
	}
	if time.Now().Unix() > int64(tokenInfo["exp"].(float64)) {
		//token timeout
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.Values().Set("msg", "Token timeout.")
		return
	} else {
		//update token
		if int64(tokenInfo["exp"].(float64))-time.Now().Unix() < 30 {
			newToken, _ := createToken(userId)
			ctx.Header("token", newToken)
			ctx.Next()
		}
	}
	oldToken, _ := token.SignedString([]byte(utils.RootSecret))
	ctx.Header("token", oldToken)
	ctx.Next()
}
