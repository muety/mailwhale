package middleware

import (
	"errors"
	"fmt"
	"github.com/muety/mailwhale/util"
	"net/http"
)

type RecoverMiddleware struct {
	handler http.Handler
}

func NewRecoverMiddleware() func(handler http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &RecoverMiddleware{
			handler: h,
		}
	}
}

func (m RecoverMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			util.RespondError(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("recovered from error – %v", err)))
			return
		}
	}()

	m.handler.ServeHTTP(w, r)
}
