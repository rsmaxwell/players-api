package httpHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/logger"
	"github.com/rsmaxwell/players-api/person"
	"github.com/rsmaxwell/players-api/session"
)

// CreateCourt method
func CreateCourt(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("CreateCourt")

	var r court.CreateCourtRequest

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted in CreateCourt Request"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse data for CreateCourt Request"))
		clientError++
		return
	}

	session := session.LookupToken(r.Token)
	if session == nil {
		WriteResponse(rw, http.StatusBadRequest, "Not Authorized")
		return
	}

	// Check there are not too many players on the court
	info, err := court.GetCourtInfo()
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, "CourtInfo error")
		serverError++
		return
	}

	if len(r.Court.Players) > info.PlayersPerCourt {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too many players on court"))
		serverError++
		return
	}

	// Check the people on the court are valid
	for _, id := range r.Court.Players {
		p, err := person.GetPersonDetails(id)
		if err != nil {
			WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not create court. err:%s", err))
			serverError++
			return
		}

		if p == nil {
			WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not create court as person [%s] does not exist", id))
			clientError++
			return
		}

		if p.Player == false {
			WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not create court as person [%s] is not a player", id))
			clientError++
			return
		}
	}

	err = court.AddCourt(r.Court)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not create a new court. err:%s", err))
		serverError++
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}
