package httpHandler

import (
	"fmt"
	"net/http"

	"github.com/rsmaxwell/players-api/person"
)

// DeletePerson method
func DeletePerson(rw http.ResponseWriter, req *http.Request, id string) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		WriteResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	err := person.DeletePerson(id)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not delete person:%s", id))
		clientError++
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
