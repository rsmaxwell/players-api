package httpHandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/logger"
)

// Court details Response
type courtDetailsResponse struct {
	Court court.Court `json:"court"`
}

// GetCourt method
func GetCourt(rw http.ResponseWriter, req *http.Request, idString string) {

	logger.Logger.Printf("writeGetCourtByIDResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		WriteResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		WriteResponse(rw, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	court, err := court.GetCourtDetails(id)

	if err != nil {
		WriteResponse(rw, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	setHeaders(rw, req)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(courtDetailsResponse{
		Court: *court,
	})

	rw.WriteHeader(http.StatusOK)
}
