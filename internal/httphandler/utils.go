package httphandler

import (
	"encoding/json"
	"fmt"
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
func writeResponseMessage(w http.ResponseWriter, r *http.Request, statusCode int, qualifier string, message string) {
	writeResponse(w, r, statusCode, qualifier)
	json.NewEncoder(w).Encode(messageResponse{
		Message: message,
	})
}

// writeResponseObject method
func writeResponseObject(w http.ResponseWriter, r *http.Request, statusCode int, qualifier string, object interface{}) {
	writeResponse(w, r, statusCode, qualifier)
	json.NewEncoder(w).Encode(object)
}

// writeResponse method
func writeResponse(w http.ResponseWriter, r *http.Request, statusCode int, qualifier string) {

	common.MetricsData.StatusCodes[statusCode]++

	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "http://localhost:4200"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Origin, XMLHttpRequest, Content-Type, X-Auth-Token, Accept, Content-Length, Accept-Encoding, X-CSRF-Token, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Headers, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if statusCode == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"players-api: "+qualifier+"\"")
	}

	w.WriteHeader(statusCode)
}

// writeResponseError function
func writeResponseError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		if serr, ok := err.(*codeerror.CodeError); ok {
			writeResponseMessage(w, r, serr.Code(), serr.Qualifier(), serr.Error())
			return
		}

		writeResponseMessage(w, r, http.StatusInternalServerError, "", "error")
		return
	}
}

// checkAuthenticated method
func checkAuthenticated(r *http.Request) (*sessions.Session, error) {
	f := functionCheckAuthenticated

	cookieHeader := r.Header.Get("Cookie")
	fmt.Printf("cookieHeader: %s\n", cookieHeader)

	sess, err := store.Get(r, "players-api")
	if err != nil {
		f.Dump("could not get the 'players-api' cookie")
		return nil, codeerror.NewInternalServerError(err.Error())
	}
	if sess.IsNew {
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	auth, ok := sess.Values["authenticated"].(bool)
	if !ok {
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}
	if !auth {
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	return sess, nil
}

// SetupHandlers Handlers for REST API routes
func SetupHandlers(w *mux.Router) {

	s := w.PathPrefix("/players-api").Subrouter()

	s.HandleFunc("/users/authenticate", Authenticate).Methods(http.MethodOptions, http.MethodPost)
	s.HandleFunc("/users/register", Register).Methods(http.MethodPost)
	s.HandleFunc("/users", ListPeople).Methods(http.MethodPost)
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

	w.NotFoundHandler = http.HandlerFunc(NotFound)
}

// Middleware method
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f := functionMiddleware

		w2 := response.New(w)

		f.DebugRequest(r)
		h.ServeHTTP(w2, r)
	})
}
