package middlewares

import (
	"log"
	"net/http"
)

type corsMiddlewareDecorator struct {
	wrapee http.Handler
}

func NewCorsMiddlewareDecorator(wrapee http.Handler) http.Handler {
	return &corsMiddlewareDecorator{
		wrapee: wrapee,
	}
}

func (decorator *corsMiddlewareDecorator) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("cors")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Timezone")
	decorator.wrapee.ServeHTTP(writer, request)
}
