package httpHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/logger"
	"github.com/rsmaxwell/players-api/person"
	"github.com/rsmaxwell/players-api/session"
)

// UpdatePerson method
func UpdatePerson(rw http.ResponseWriter, req *http.Request, id string) {

	logger.Logger.Printf("writeUpdatePersonResponse")

	var r person.UpdatePersonRequest

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse person data"))
		clientError++
		return
	}

	session := session.LookupToken(r.Token)
	if session == nil {
		WriteResponse(rw, http.StatusBadRequest, "Not Authorized")
		clientError++
		return
	}

	_, err = person.UpdatePerson(id, session, r.Person)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not update person"))
		clientError++
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
