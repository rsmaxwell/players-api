package logger

import (
	"log"
	"net/http"
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
