package logger

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
)

var (
	// Logger is the common instance of the logger
	Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

// LogHandler logs an http reposonse
func LogHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		Logger.Printf("%q", x)
		rec := httptest.NewRecorder()
		fn(rec, r)
		Logger.Printf("status: %d", rec.Code)
		Logger.Printf("%q", rec.Body)

		w.WriteHeader(rec.Code)

		h := w.Header()
		for k, v := range rec.HeaderMap {
			h[k] = v
		}

		w.Write(rec.Body.Bytes())
	}
}
