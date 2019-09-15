package httpHandler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/logger"
)

// DeleteCourt method
func DeleteCourt(rw http.ResponseWriter, req *http.Request, idString string) {

	logger.Logger.Printf("writeDeleteCourtResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		WriteResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	// Convert the ID into a number
	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Logger.Printf(err.Error())
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("The ID:%s is not a number", idString))
		clientError++
		return
	}

	err = court.DeleteCourt(id)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not delete court:%s", idString))
		clientError++
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
