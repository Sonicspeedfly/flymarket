package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
)

var ErrNoAuthentication = errors.New("no authentication")

var authenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}

type IDFunc func(ctx context.Context, token string) (int64, error)

func Authenticate(idFunc IDFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			token := request.Header.Get("Authorization")

			if token != "" {
				id, err := idFunc(request.Context(), token)
				if err != nil {
					log.Println("middleware Authenticate idFunc ERROR:", err)
					http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				ctx := context.WithValue(request.Context(), authenticationContextKey, id)
				request = request.WithContext(ctx)
			}

			handler.ServeHTTP(writer, request)
		})
	}
}

func Authentication(ctx context.Context) (int64, error) {
	if value, ok := ctx.Value(authenticationContextKey).(int64); ok {
		return value, nil
	}
	return 0, ErrNoAuthentication
}
