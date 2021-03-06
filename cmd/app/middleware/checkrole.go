package middleware

import (
	"context"
	"log"
	"net/http"
)

type HasAnyRoleFunc func(ctx context.Context, id int64, roles ...string) bool

func CheckRole(hasAnyRoleFunc HasAnyRoleFunc, roles ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			id, err := Authentication(request.Context())
			if err != nil && err != ErrNoAuthentication {
				log.Println("managers HasAnyRole middleware.Authentication ERROR:", err)
				http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			if id != 0 && !hasAnyRoleFunc(request.Context(), id, roles...) {
				http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			handler.ServeHTTP(writer, request)
		})
	}
}