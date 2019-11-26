package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// GetPersonResponse structure
type GetPersonResponse struct {
	Person person.Person `json:"person"`
}

var (
	functionGetPerson = debug.NewFunction(pkg, "GetPerson")
)

// GetPerson method
func GetPerson(rw http.ResponseWriter, req *http.Request) {
	f := functionGetPerson

	_, err := checkAuthToken(req)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	p, err := model.GetPerson(id)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetPersonResponse{
		Person: *p,
	})
}
