package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"

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

	functionMiddleware = debug.NewFunction(pkg, "Middleware")
)

// WriteResponse method
func WriteResponse(w http.ResponseWriter, httpStatus int, message string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	json.NewEncoder(w).Encode(messageResponse{
		Message: message,
	})
}

func setHeaders(rw http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")

	if origin == "" {
		origin = "http://localhost:4200"
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", origin)
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Access-Control-Allow-Origin, Authorization")
}

// errorHandler function
func errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		setHeaders(rw, req)
		if serr, ok := err.(*codeerror.CodeError); ok {
			WriteResponse(rw, serr.Code(), serr.Error())
			common.MetricsData.ClientError++
			return
		}

		WriteResponse(rw, http.StatusInternalServerError, "InternalServerError")
		common.MetricsData.ClientError++
		return
	}
}

// checkAuthToken method
func checkAuthToken(req *http.Request) (*Claims, error) {

	cookie, err := req.Cookie("players-api")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, codeerror.NewUnauthorized(err.Error())
		}
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	tokenString := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, codeerror.NewUnauthorized(err.Error())
		}
		return nil, codeerror.NewBadRequest(err.Error())
	}
	if !token.Valid {
		return nil, codeerror.NewUnauthorized(err.Error())
	}

	return claims, nil
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
