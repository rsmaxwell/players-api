package httpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/person"
)

// People details Response
type personDetailsResponse struct {
	Person person.Person `json:"person"`
}

// GetPerson method
func GetPerson(rw http.ResponseWriter, req *http.Request, id string) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		WriteResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	person, err := person.GetPersonDetails(id)
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
