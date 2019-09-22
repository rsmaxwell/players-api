package httpHandler

import (
	"encoding/json"
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

// GetPersonResponse structure
type GetPersonResponse struct {
	Person person.Person `json:"person"`
}

// GetPerson method
func GetPerson(rw http.ResponseWriter, req *http.Request, id string) {

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	var r GetPersonRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	session := session.LookupToken(r.Token)
	if session == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		clientError++
		return
	}

	person, err := person.Load(id)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}
	if person == nil {
		WriteResponse(rw, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetPersonResponse{
		Person: *person,
	})
}
