package httphandler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/response"
)

// messageResponse structure
type MessageResponse struct {
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
	key   = []byte("<SESSION_SECRET_KEY>")
	store = sessions.NewCookieStore(key)
)

var (
	pkg = debug.NewPackage("httphandler")

	functionMiddleware         = debug.NewFunction(pkg, "Middleware")
	functionSignin             = debug.NewFunction(pkg, "Signin")
	functionCheckAuthenticated = debug.NewFunction(pkg, "checkAuthenticated")
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
		w.Header().Add("WWW-Authenticate", "Basic realm=players-api")
	}

	w.WriteHeader(statusCode)
}

// writeResponseError function
func writeResponseError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		if serr, ok := err.(*codeerror.CodeError); ok {
			writeResponseMessage(w, r, serr.Code(), serr.Error())
			return
		}

		writeResponseMessage(w, r, http.StatusInternalServerError, "error")
		return
	}
}

// checkAuthenticated method
func checkAuthenticated(r *http.Request) (int, error) {
	f := functionCheckAuthenticated

	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		f.DebugError("missing Authorization header")
		return 0, fmt.Errorf("not authorized")
	}

	splitToken := strings.Split(authorizationHeader, "Bearer ")
	if splitToken == nil {
		return 0, fmt.Errorf("not authorized")
	}

	tokenString := splitToken[1]

	claims, err := basic.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.ID, nil
}

// SetupHandlers Handlers for REST API routes
func SetupHandlers(w *mux.Router) {

	s := w.PathPrefix("/players-api").Subrouter()

	s.HandleFunc("/register", Register).Methods(http.MethodOptions, http.MethodPost)
	s.HandleFunc("/signin", Signin).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("/signout", Signout).Methods(http.MethodGet)
	s.HandleFunc("/refresh", RefreshTokens).Methods(http.MethodGet)

	s.HandleFunc("/people", ListPeople).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("/people/{id}", DeletePerson).Methods(http.MethodDelete)
	s.HandleFunc("/people/{id}", GetPerson).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/people/{id}", UpdatePerson).Methods(http.MethodPut)

	s.HandleFunc("/people/toplayer/{id1}", MakePersonPlayer).Methods(http.MethodPut)
	s.HandleFunc("/people/toinactive/{id}", MakePersonInactive).Methods(http.MethodPut)

	s.HandleFunc("/people/toplaying/{id1}/{id2}/{id3}", MakePlayerPlay).Methods(http.MethodPut)
	s.HandleFunc("/people/towaiting/{id}", MakePlayerWait).Methods(http.MethodPut)

	s.HandleFunc("/courts", ListCourts).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/courts/{id}", GetCourt).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/courts", CreateCourt).Methods(http.MethodPost)
	s.HandleFunc("/courts/{id}", UpdateCourt).Methods(http.MethodPut)
	s.HandleFunc("/courts/{id}", DeleteCourt).Methods(http.MethodDelete)
	s.HandleFunc("/courts/fill/{id}", FillCourt).Methods(http.MethodPut, http.MethodOptions)
	s.HandleFunc("/courts/clear/{id}", ClearCourt).Methods(http.MethodPut, http.MethodOptions)

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
