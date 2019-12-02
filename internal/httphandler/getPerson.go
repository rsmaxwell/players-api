package httphandler

import (
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

	_, err := checkAccessToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	p, err := model.GetPerson(id)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseObject(rw, req, http.StatusOK, "", GetPersonResponse{
		Person: *p,
	})
}
