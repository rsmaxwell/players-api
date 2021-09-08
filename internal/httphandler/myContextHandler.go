package httphandler

import (
	"context"
	"net/http"
	"time"
)

var (
	requestCounter = 0
)

func nextRequestID() int {
	requestCounter++
	return requestCounter
}

type MyContext struct {
	handler http.Handler
}

func AddRequestContext(handlerToWrap http.Handler) *MyContext {
	return &MyContext{handler: handlerToWrap}
}

func (h *MyContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx1, cancel := context.WithTimeout(r.Context(), time.Duration(60*time.Second))
	defer cancel()

	ctx2 := context.WithValue(ctx1, ContextRequestIdKey, nextRequestID())
	r3 := r.WithContext(ctx2)

	h.handler.ServeHTTP(w, r3)
}
