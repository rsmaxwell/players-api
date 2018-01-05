package logger

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

var (
	// Logger is the common instance of the logger
	Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) {

	Logger.Printf("Request:\n")
	Logger.Printf("  url : %v %v %v\n", r.Method, r.URL, r.Proto)
	Logger.Printf("  host: %v\n", r.Host)
	Logger.Println("  headers:")
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			Logger.Printf("    %v: %v\n", name, h)
		}
	}

	if r.Method == "POST" {
		r.ParseForm()
		Logger.Printf("  post data:\n")
		Logger.Printf("%s\n", r.Form.Encode())
	}
}

// LogHandler logs an http reposonse
func LogHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formatRequest(r)

		recorder := httptest.NewRecorder()
		fn(recorder, r)

		resp := recorder.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		Logger.Printf("Response: %d", resp.StatusCode)
		Logger.Println("  Headers:")
		for k, v := range recorder.HeaderMap {
			Logger.Printf("    %s : %s", k, v)
		}
		Logger.Println("  Body : " + string(body))

		for k, values := range recorder.HeaderMap {
			for _, v := range values {
				w.Header().Set(k, v)
			}
		}

		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	}
}
