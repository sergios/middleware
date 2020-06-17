package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/sergios/errors"
	"github.com/sergios/render"

	log "github.com/Sirupsen/logrus"
)

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	PrintStack bool
	StackAll   bool
	StackSize  int
}

// NewRecovery returns a new instance of Recovery
func NewRecovery() *Recovery {
	return &Recovery{
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			log.WithFields(log.Fields{
				"request": r.RequestURI,
				"err":     fmt.Sprintf("%s\n%s", err, stack),
			}).Error("Panic request")

			if rec.PrintStack {
				render.WriteError(rw, errors.Http500)
			}
		}
	}()

	next(rw, r)
}
