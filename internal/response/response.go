package response

import (
	"net/http"
	"strings"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	pkg *debug.Package
)

func init() {
	pkg = debug.NewPackage("response")
}

// Wrapper type
type Wrapper struct {
	http.ResponseWriter
	w http.ResponseWriter
}

// New function
func New(rw http.ResponseWriter) http.ResponseWriter {
	wrap := Wrapper{w: rw}
	return &wrap
}

// Header function
func (r *Wrapper) Header() http.Header {
	return r.w.Header() // pass it to the original ResponseWriter
}

// Write function
func (r *Wrapper) Write(b []byte) (int, error) {
	f := debug.NewFunction(pkg, "Write")
	f.DebugVerbose(strings.TrimSuffix(string(b), "\n")) // log it out
	return r.w.Write(b)                                 // pass it to the original ResponseWriter
}

// WriteHeader function
func (r *Wrapper) WriteHeader(statusCode int) {
	f := debug.NewFunction(pkg, "WriteHeader")
	f.DebugVerbose("statusCode: %d", statusCode)
	r.w.WriteHeader(statusCode) // pass it to the original ResponseWriter
}
