package httpHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/logger"
)

// UpdateCourt method
func UpdateCourt(rw http.ResponseWriter, req *http.Request, idString string) {

	logger.Logger.Printf("writeUpdateCourtResponse")

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

	var c court.JSONCourt

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse court data for id:%s", idString))
		clientError++
		return
	}

	_, err = court.UpdateCourt(id, c)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not update court:%s", idString))
		clientError++
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}