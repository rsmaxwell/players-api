package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/person"
	"github.com/rsmaxwell/players-api/queue"
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
		if serr, ok := err.(*codeError.CodeError); ok {
			WriteResponse(rw, serr.Code(), serr.Error())
			clientError++
			return
		}

		WriteResponse(rw, http.StatusInternalServerError, "InternalServerError")
		clientError++
		return
	}
}

// CheckConsistency checks the state on disk is consistent
func CheckConsistency() error {

	err := person.CheckConsistency()
	if err != nil {
		return err
	}

	err = court.CheckConsistency()
	if err != nil {
		return err
	}

	err = queue.CheckConsistency()
	if err != nil {
		return err
	}

	return nil
}
