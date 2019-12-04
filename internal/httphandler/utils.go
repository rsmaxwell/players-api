package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/response"
)

// messageResponse structure
type messageResponse struct {
	Message string `json:"message"`
}

var (
	contextPath = "/players-api"

	pkg = debug.NewPackage("httphandler")

	functionMiddleware         = debug.NewFunction(pkg, "Middleware")
	functionAuthenticate       = debug.NewFunction(pkg, "Authenticate")
	functionCheckAuthenticated = debug.NewFunction(pkg, "checkAuthenticated")
)

// writeResponseMessage method
func writeResponseMessage(r http.ResponseWriter, req *http.Request, statusCode int, qualifier string, message string) {
	writeResponse(r, req, statusCode, qualifier)
	json.NewEncoder(r).Encode(messageResponse{
		Message: message,
	})
}

// writeResponseObject method
func writeResponseObject(r http.ResponseWriter, req *http.Request, statusCode int, qualifier string, object interface{}) {
	writeResponse(r, req, statusCode, qualifier)
	json.NewEncoder(r).Encode(object)
}

// writeResponse method
func writeResponse(r http.ResponseWriter, req *http.Request, statusCode int, qualifier string) {

	common.MetricsData.StatusCodes[statusCode]++

	if statusCode == http.StatusOK {

		origin := req.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:4200"
		}

		r.Header().Set("Content-Type", "application/json")
		r.Header().Set("Access-Control-Allow-Origin", origin)
		r.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		r.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Access-Control-Allow-Origin, Authorization")

	} else if statusCode == http.StatusUnauthorized {

		r.Header().Set("WWW-Authenticate", "Basic realm=\"players-api: "+qualifier+"\"")

	} else {

	}

	r.WriteHeader(statusCode)
}

// writeResponseError function
func writeResponseError(rw http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		if serr, ok := err.(*codeerror.CodeError); ok {
			writeResponseMessage(rw, req, serr.Code(), serr.Qualifier(), serr.Error())
			return
		}

		writeResponseMessage(rw, req, http.StatusInternalServerError, "", "error")
		return
	}
}

// checkAuthenticated method
func checkAuthenticated(req *http.Request) (*sessions.Session, error) {
	f := functionCheckAuthenticated

	sess, err := store.Get(req, "players-api")
	if err != nil {
		f.Dump("could not get the 'players-api' cookie")
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	auth, ok := sess.Values["authenticated"].(bool)
	if !ok {
		return nil, codeerror.NewForbidden("Forbidden")
	}
	if !auth {
		return nil, codeerror.NewForbidden("Forbidden")
	}

	return sess, nil
}

// SetupHandlers Handlers for REST API routes
func SetupHandlers(r *mux.Router) {

	s := r.PathPrefix("/players-api").Subrouter()

	s.HandleFunc("/users/authenticate", Authenticate).Methods(http.MethodPost)
	s.HandleFunc("/users/register", Register).Methods(http.MethodPost)
	s.HandleFunc("/users", ListPeople).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", DeletePerson).Methods(http.MethodDelete)
	s.HandleFunc("/users/logout", Logout).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", GetPerson).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", UpdatePerson).Methods(http.MethodPut)
	s.HandleFunc("/users/player/{id}", UpdatePersonPlayer).Methods(http.MethodPut)
	s.HandleFunc("/users/role/{id}", UpdatePersonRole).Methods(http.MethodPut)
	s.HandleFunc("/users/move", PostMove).Methods(http.MethodPost)

	s.HandleFunc("/court", ListCourts).Methods(http.MethodGet)
	s.HandleFunc("/court/{id}", GetCourt).Methods(http.MethodGet)
	s.HandleFunc("/court", CreateCourt).Methods(http.MethodPost)
	s.HandleFunc("/court/{id}", UpdateCourt).Methods(http.MethodPut)
	s.HandleFunc("/court/{id}", DeleteCourt).Methods(http.MethodDelete)

	s.HandleFunc("/queue", GetQueue).Methods(http.MethodGet)
	s.HandleFunc("/metrics", GetMetrics).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(NotFound)
}

// Middleware method
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		f := functionMiddleware

		rw2 := response.New(rw)

		f.DebugRequest(req)
		h.ServeHTTP(rw2, req)
	})
}
