package httpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/person"
)

// List people Response
type listPeopleResponse struct {
	People []int `json:"people"`
}

// GetAllPeople method
func GetAllPeople(rw http.ResponseWriter, req *http.Request) {

	setHeaders(rw, req)

	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		WriteResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	listOfPeople, err := person.ListAllPeople()
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, "Error getting list of people")
		serverError++
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(listPeopleResponse{
		People: listOfPeople,
	})
}
