package httphandler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/response"
)

// messageResponse structure
type messageResponse struct {
	Message string `json:"message"`
}

// ContextKey type
type ContextKey string

const (
	contextPath = "/players-api"

	// ContextDatabaseKey constant
	ContextDatabaseKey ContextKey = "database"
)

var (
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

	model.MetricsData.StatusCodes[statusCode]++

	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "http://localhost:4200"
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers",
		"Origin, XMLHttpRequest, Content-Type, X-Auth-Token, Accept, Content-Length, Accept-Encoding, X-CSRF-Token, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Headers, Authorization")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	if statusCode == http.StatusUnauthorized {
		w.Header().Add("WWW-Authenticate", "Basic realm=\"players-api: "+qualifier+"\"")
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

	s.HandleFunc("/users/register", Register).Methods(http.MethodPost)
	s.HandleFunc("/users/authenticate", Authenticate).Methods(http.MethodOptions, http.MethodPost)
	s.HandleFunc("/users", ListPeople).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", DeletePerson).Methods(http.MethodDelete)
	s.HandleFunc("/users/logout", Logout).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", GetPerson).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", UpdatePerson).Methods(http.MethodPut)
	s.HandleFunc("/users/toplaying/{id1}/{id2}", MakePlaying).Methods(http.MethodPut)
	s.HandleFunc("/users/towaiting/{id}", MakeWaiting).Methods(http.MethodPut)
	s.HandleFunc("/users/toinactive/{id}", MakeInactive).Methods(http.MethodPut)

	s.HandleFunc("/court", ListCourts).Methods(http.MethodGet)
	s.HandleFunc("/court/{id}", GetCourt).Methods(http.MethodGet)
	s.HandleFunc("/court", CreateCourt).Methods(http.MethodPost)
	s.HandleFunc("/court/{id}", UpdateCourt).Methods(http.MethodPut)
	s.HandleFunc("/court/{id}", DeleteCourt).Methods(http.MethodDelete)

	s.HandleFunc("/metrics", GetMetrics).Methods(http.MethodGet)

	w.NotFoundHandler = http.HandlerFunc(NotFound)
}

// Middleware method
func Middleware(h http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f := functionMiddleware

		w2 := response.New(w)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60*time.Second))
		defer cancel()
		r2 := r.WithContext(ctx)

		ctx = context.WithValue(r2.Context(), ContextDatabaseKey, db)
		r3 := r.WithContext(ctx)

		f.DebugRequest(r3)
		h.ServeHTTP(w2, r3)
	})
}
