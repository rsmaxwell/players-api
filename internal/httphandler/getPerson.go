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
func GetPerson(w http.ResponseWriter, r *http.Request) {
	f := functionGetPerson

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	id := mux.Vars(r)["id"]
	f.DebugVerbose("ID: %s", id)

	p, err := model.GetPerson(id)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", GetPersonResponse{
		Person: *p,
	})
}
