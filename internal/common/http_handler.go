package common

import (
	"log"
	"net/http"
)

type Method func(writer http.ResponseWriter, request *http.Request)

type httpHandler struct {
	get     Method
	post    Method
	put     Method
	patch   Method
	delete  Method
	head    Method
	options Method
}

func (handler *httpHandler) handle(method Method, writer http.ResponseWriter, request *http.Request) {
	if method == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	method(writer, request)
}

func (handler *httpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("method", request.Method)
	switch request.Method {
	case "GET":
		handler.handle(handler.get, writer, request)
	case "POST":
		handler.handle(handler.post, writer, request)
	case "PUT":
		handler.handle(handler.put, writer, request)
	case "PATCH":
		handler.handle(handler.patch, writer, request)
	case "DELETE":
		handler.handle(handler.delete, writer, request)
	case "HEAD":
		handler.handle(handler.head, writer, request)
	case "OPTIONS":
		writer.WriteHeader(http.StatusNoContent)
	}
}

type httpHandlerBuilder struct {
	get     Method
	post    Method
	put     Method
	patch   Method
	delete  Method
	head    Method
	options Method
}

func NewHttpHandlerBuilder() *httpHandlerBuilder {
	return &httpHandlerBuilder{}
}

func (builder *httpHandlerBuilder) Get(get Method) *httpHandlerBuilder {
	builder.get = get
	return builder
}

func (builder *httpHandlerBuilder) Post(post Method) *httpHandlerBuilder {
	builder.post = post
	return builder
}

func (builder *httpHandlerBuilder) Put(put Method) *httpHandlerBuilder {
	builder.put = put
	return builder
}

func (builder *httpHandlerBuilder) Patch(patch Method) *httpHandlerBuilder {
	builder.patch = patch
	return builder
}

func (builder *httpHandlerBuilder) Delete(delete Method) *httpHandlerBuilder {
	builder.delete = delete
	return builder
}

func (builder *httpHandlerBuilder) Head(head Method) *httpHandlerBuilder {
	builder.head = head
	return builder
}

func (builder *httpHandlerBuilder) Options(options Method) *httpHandlerBuilder {
	builder.options = options
	return builder
}

func (builder *httpHandlerBuilder) Build() http.Handler {
	handler := &httpHandler{
		get:     builder.get,
		post:    builder.post,
		put:     builder.put,
		patch:   builder.patch,
		delete:  builder.delete,
		head:    builder.head,
		options: builder.options,
	}
	return handler
}
