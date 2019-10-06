package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/codeerror"
)

var (
	clientSuccess             int
	clientError               int
	clientAuthenticationError int
	serverError               int
)

// messageResponse structure
type messageResponse struct {
	Message string `json:"message"`
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
			clientError++
			return
		}

		WriteResponse(rw, http.StatusInternalServerError, "InternalServerError")
		clientError++
		return
	}
}

// SetupHandlers Handlers for REST API routes
func SetupHandlers(r *mux.Router) {

	r.HandleFunc("/register",
		func(w http.ResponseWriter, req *http.Request) {
			Register(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/login",
		func(w http.ResponseWriter, req *http.Request) {
			Login(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/court",
		func(w http.ResponseWriter, req *http.Request) {
			ListCourts(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			GetCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodGet)

	r.HandleFunc("/court",
		func(w http.ResponseWriter, req *http.Request) {
			CreateCourt(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			UpdateCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			DeleteCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodDelete)

	r.HandleFunc("/person",
		func(w http.ResponseWriter, req *http.Request) {
			ListPeople(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			GetPerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodGet)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			UpdatePerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/personplayer/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			UpdatePersonPlayer(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/personrole/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			UpdatePersonRole(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			DeletePerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodDelete)

	r.HandleFunc("/metrics",
		func(w http.ResponseWriter, req *http.Request) {
			GetMetrics(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/queue",
		func(w http.ResponseWriter, req *http.Request) {
			GetQueue(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/move",
		func(w http.ResponseWriter, req *http.Request) {
			PostMove(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/queue",
		func(w http.ResponseWriter, req *http.Request) {
			GetQueue(w, req)
		}).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			NotFound(w, req)
		})
}
