package middleware

import (
	"net/http"
	"os"

	"github.com/sergios/errors"
	"github.com/sergios/render"
)

type Authorization struct {
	HeaderKey  string
	TokenCheck string
}

func NewAuthorization(headerKey string, envKey string) *Authorization {
	return &Authorization{
		HeaderKey:  headerKey,
		TokenCheck: os.Getenv(envKey),
	}
}

func (a *Authorization) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if a.isAuthorized(r) {
		next(w, r)
		return
	}
	render.WriteError(w, errors.UserNotAuthorized)
}

func (a *Authorization) isAuthorized(r *http.Request) (authorized bool) {
	// liberar tudo em caso de GET
	if r.Method == "GET" || r.RequestURI == "/healthcheck" || r.RequestURI == "/version" || r.RequestURI == "/status" {
		return true
	}

	tokenHeader := r.Header.Get(a.HeaderKey)
	if a.TokenCheck != "" && tokenHeader != "" && tokenHeader == a.TokenCheck {
		return true
	}
	return false
}
