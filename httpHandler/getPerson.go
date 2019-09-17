package httpHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/person"
	"github.com/rsmaxwell/players-api/session"
)

// GetPersonRequest structure
type GetPersonRequest struct {
	Token string `json:"token"`
}

// People details Response
type personDetailsResponse struct {
	Person person.Person `json:"person"`
}

// GetPerson method
func GetPerson(rw http.ResponseWriter, req *http.Request, id string) {

	var r GetPersonRequest

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse data"))
		clientError++
		return
	}

	session := session.LookupToken(r.Token)
	if session == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		clientError++
		return
	}

	person, err := person.Get(id)
	if err != nil {
		WriteResponse(rw, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	setHeaders(rw, req)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(personDetailsResponse{
		Person: *person,
	})
}
