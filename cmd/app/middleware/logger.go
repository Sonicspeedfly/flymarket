package middleware

import (
	"log"
	"net/http"
)

func Logger(input http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("START %s %s\n", request.Method, request.URL.Path)
		input.ServeHTTP(writer, request)
		log.Printf("FINISH %s %s\n", request.Method, request.URL.Path)
	})
}
