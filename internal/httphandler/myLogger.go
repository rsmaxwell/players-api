package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionHandlerFunc = debug.NewFunction(pkg, "HandlerFunc")
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		f := functionHandlerFunc

		DebugRequest(f, request)

		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		h.ServeHTTP(recorder, request)

		DebugResponse(f, request, recorder.Status)
	})
}
