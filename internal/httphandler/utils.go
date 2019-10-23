package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
	pkg *debug.Package
)

func init() {
	pkg = debug.NewPackage("httphandler")
}

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

// SetupHandlers Handlers for REST API routes
func SetupHandlers(r *mux.Router) {

	s := r.PathPrefix("/players-api").Subrouter()

	s.Handle("/register", Middleware(http.HandlerFunc(Register))).Methods(http.MethodPost)
	s.Handle("/login", Middleware(http.HandlerFunc(Login))).Methods(http.MethodGet)
	s.Handle("/court", Middleware(http.HandlerFunc(ListCourts))).Methods(http.MethodGet)
	s.Handle("/court/{id}", Middleware(http.HandlerFunc(GetCourt))).Methods(http.MethodGet)
	s.Handle("/court", Middleware(http.HandlerFunc(CreateCourt))).Methods(http.MethodPost)
	s.Handle("/court/{id}", Middleware(http.HandlerFunc(UpdateCourt))).Methods(http.MethodPut)
	s.Handle("/court/{id}", Middleware(http.HandlerFunc(DeleteCourt))).Methods(http.MethodDelete)
	s.Handle("/person", Middleware(http.HandlerFunc(ListPeople))).Methods(http.MethodGet)
	s.Handle("/person/{id}", Middleware(http.HandlerFunc(GetPerson))).Methods(http.MethodGet)
	s.Handle("/person/{id}", Middleware(http.HandlerFunc(UpdatePerson))).Methods(http.MethodPut)
	s.Handle("/personplayer/{id}", Middleware(http.HandlerFunc(UpdatePersonPlayer))).Methods(http.MethodPut)
	s.Handle("/personrole/{id}", Middleware(http.HandlerFunc(UpdatePersonRole))).Methods(http.MethodPut)
	s.Handle("/person/{id}", Middleware(http.HandlerFunc(DeletePerson))).Methods(http.MethodDelete)
	s.Handle("/metrics", Middleware(http.HandlerFunc(GetMetrics))).Methods(http.MethodGet)
	s.Handle("/queue", Middleware(http.HandlerFunc(GetQueue))).Methods(http.MethodGet)
	s.Handle("/move", Middleware(http.HandlerFunc(PostMove))).Methods(http.MethodPost)
	s.Handle("/queue", Middleware(http.HandlerFunc(GetQueue))).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(NotFound)
}

// Middleware method
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		f := debug.NewFunction(pkg, "Middleware")

		rw2 := response.New(rw)

		f.DebugRequest(req)
		h.ServeHTTP(rw2, req)
	})
}
