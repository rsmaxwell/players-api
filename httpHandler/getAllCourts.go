package httpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/logger"
)

// List courts Response
type listCourtsResponse struct {
	Courts []int `json:"courts"`
}

// GetAllCourts method
func GetAllCourts(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeGetListOfCourtsResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		WriteResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	listOfCourts, err := court.ListAllCourts()
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, "Error getting list of courts")
		serverError++
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(listCourtsResponse{
		Courts: listOfCourts,
	})
}
