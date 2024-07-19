package middlewares

import (
	"log"
	"net/http"

	"github.com/gustavlsouz/documents-service/internal/common"
)

type counterMiddlewareDecorator struct {
	wrapee http.Handler
}

func NewCounterMiddlewareDecorator(wrapee http.Handler) http.Handler {
	return &counterMiddlewareDecorator{
		wrapee: wrapee,
	}
}

func (decorator *counterMiddlewareDecorator) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	common.GetSingletonApplicationUpTime().AddRequest()
	log.Println("added a request for count")
	decorator.wrapee.ServeHTTP(writer, request)
}
