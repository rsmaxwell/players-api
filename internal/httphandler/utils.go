package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// messageResponse structure
type MessageResponse struct {
	Message string `json:"message"`
}

// ContextKey type
type ContextKey string

const (
	contextPath = "/players-api"

	// Context Keys
	ContextDatabaseKey  ContextKey = "database"
	ContextRequestIdKey ContextKey = "requestID"
	ContextConfigKey    ContextKey = "config"
)

var (
	key   = []byte("<SESSION_SECRET_KEY>")
	store = sessions.NewCookieStore(key)
)

var (
	pkg = debug.NewPackage("httphandler")

	functionGetRequestID       = debug.NewFunction(pkg, "getRequestID")
	functionHidePasswords      = debug.NewFunction(pkg, "hidePasswords")
	functionWriteResponseError = debug.NewFunction(pkg, "writeResponseError")
)

// writeResponseMessage method
func writeResponseMessage(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	writeResponse(w, r, statusCode)
	json.NewEncoder(w).Encode(MessageResponse{
		Message: message,
	})
}

// writeResponseObject method
func writeResponseObject(w http.ResponseWriter, r *http.Request, statusCode int, object interface{}) {
	writeResponse(w, r, statusCode)
	json.NewEncoder(w).Encode(object)
}

// writeResponse method
func writeResponse(w http.ResponseWriter, r *http.Request, statusCode int) {
	model.MetricsData.StatusCodes[statusCode]++
	w.WriteHeader(statusCode)
}

// writeResponseError function
func writeResponseError(writer http.ResponseWriter, request *http.Request, err error) {
	f := functionWriteResponseError
	DebugVerbose(f, request, err.Error())

	if serr, ok := err.(*codeerror.CodeError); ok {
		writeResponseMessage(writer, request, serr.Code(), serr.Error())
		return
	}

	writeResponseMessage(writer, request, http.StatusInternalServerError, "error")
}

// SetupHandlers Handlers for REST API routes
func SetupHandlers(w *mux.Router) {

	s := w.PathPrefix("/players-api").Subrouter()

	s.HandleFunc("/register", Register).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("/signin", Signin).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("/signout", Signout).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/refresh", RefreshToken).Methods(http.MethodPost, http.MethodOptions)

	s.HandleFunc("/waiters", ListWaiters).Methods(http.MethodGet, http.MethodOptions)

	s.HandleFunc("/people", ListPeople).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("/people/{id}", DeletePerson).Methods(http.MethodDelete, http.MethodOptions)
	s.HandleFunc("/people/{id}", GetPerson).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/people/{id}", UpdatePerson).Methods(http.MethodPut, http.MethodOptions)

	s.HandleFunc("/people/toplayer/{id1}", MakePersonPlayer).Methods(http.MethodPut, http.MethodOptions)
	s.HandleFunc("/people/toinactive/{id}", MakePersonInactive).Methods(http.MethodPut, http.MethodOptions)

	s.HandleFunc("/people/toplaying/{id1}/{id2}/{id3}", MakePlayerPlay).Methods(http.MethodPut, http.MethodOptions)
	s.HandleFunc("/people/towaiting/{id}", MakePlayerWait).Methods(http.MethodPut, http.MethodOptions)

	s.HandleFunc("/courts", ListCourts).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/courts/{id}", GetCourt).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/courts", CreateCourt).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("/courts/{id}", UpdateCourt).Methods(http.MethodPut, http.MethodOptions)
	s.HandleFunc("/courts/{id}", DeleteCourt).Methods(http.MethodDelete, http.MethodOptions)
	s.HandleFunc("/courts/fill/{id}", FillCourt).Methods(http.MethodPut, http.MethodOptions)
	s.HandleFunc("/courts/clear/{id}", ClearCourt).Methods(http.MethodPut, http.MethodOptions)

	s.HandleFunc("/metrics", GetMetrics).Methods(http.MethodGet, http.MethodOptions)

	w.NotFoundHandler = http.HandlerFunc(NotFound)
}

func DebugRequest(f *debug.Function, request *http.Request) error {

	if f.Level() >= debug.APILevel {

		requestID := getRequestID(request)
		f.DebugAPI("%s %s %s %s ------------------------------------------------------", requestID, request.Method, request.Host, request.URL)

		for name, headers := range request.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				f.DebugVerbose("          %v: %v", name, h)
			}
		}
	}

	return nil
}

func DebugResponse(f *debug.Function, request *http.Request, status int) {
	if f.Level() >= debug.APILevel {
		requestID := getRequestID(request)
		f.DebugAPI("%s %s %s %s --> statusCode: %d", requestID, request.Method, request.Host, request.URL, status)
	}
}

// DebugRequestBody traces the http request body
func DebugRequestBody(f *debug.Function, req *http.Request, data []byte) {
	if f.Level() >= debug.APILevel {
		requestID := getRequestID(req)
		data2, _ := hidePasswords(data)
		f.DebugVerbose("%s %s", requestID, string(data2))
	}
}

func DebugInfo(f *debug.Function, req *http.Request, format string, a ...interface{}) {
	if f.Level() >= debug.VerboseLevel {
		requestID := getRequestID(req)
		message := fmt.Sprintf(format, a...)
		f.DebugInfo("%s %s", requestID, message)
	}
}

func DebugError(f *debug.Function, req *http.Request, format string, a ...interface{}) {
	if f.Level() >= debug.VerboseLevel {
		requestID := getRequestID(req)
		message := fmt.Sprintf(format, a...)
		f.DebugError("%s %s", requestID, message)
	}
}

func DebugVerbose(f *debug.Function, req *http.Request, format string, a ...interface{}) {
	if f.Level() >= debug.VerboseLevel {
		requestID := getRequestID(req)
		message := fmt.Sprintf(format, a...)
		f.DebugVerbose("%s %s", requestID, message)
	}
}

// DebugRequestBody traces the http request body
func getRequestID(req *http.Request) string {
	f := functionGetRequestID

	ctx := req.Context()
	object := ctx.Value(ContextRequestIdKey)
	id, ok := object.(int)
	if !ok {
		message := fmt.Sprintf("unexpected context type: %#v", id)
		f.Dump(message)
		return ""
	}

	return fmt.Sprintf("[request:%d]", id)
}

func hidePasswords(data []byte) ([]byte, error) {
	f := functionHidePasswords

	var input map[string]interface{}
	err := json.Unmarshal(data, &input)
	if err != nil {
		d := f.DumpError(err, "Could not parse data")
		d.AddByteArray("data", data)
		return nil, err
	}

	output := walk(input)

	var array []byte
	array, err = json.Marshal(output)
	if err != nil {
		d := f.DumpError(err, "Could not parse data")
		d.AddByteArray("array", array)
		return nil, err
	}
	return array, nil
}

func walk(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range input {
		z, ok := v.(map[string]interface{})
		if ok {
			output[k] = walk(z)
		} else if k == "password" {
			output[k] = "********"
		} else {
			output[k] = v
		}
	}
	return output
}

func Dump(f *debug.Function, request *http.Request, format string, a ...interface{}) *debug.Dump {
	d := f.Dump(format, a...)
	d.AddString("RequestID", getRequestID(request))
	return d
}

func DumpError(f *debug.Function, request *http.Request, err error, format string, a ...interface{}) *debug.Dump {
	d := f.DumpError(err, format, a...)
	d.AddString("RequestID", getRequestID(request))
	return d
}
