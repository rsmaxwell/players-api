package httphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/destination"
	"github.com/rsmaxwell/players-api/person"
	"github.com/rsmaxwell/players-api/session"
)

// CreateCourtRequest structure
type CreateCourtRequest struct {
	Token string            `json:"token"`
	Court destination.Court `json:"court"`
}

// CreateCourtResponse structure
type CreateCourtResponse struct {
	ID string `json:"id"`
}

// CreateCourt method
func CreateCourt(rw http.ResponseWriter, req *http.Request) {

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	var r CreateCourtRequest
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

	// Check there are not too many players on the court
	info, err := destination.GetCourtInfo()
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	if len(r.Court.Container.Players) > info.PlayersPerCourt {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too many players on court"))
		serverError++
		return
	}

	// Check the people on the court are valid
	for _, id := range r.Court.Container.Players {
		p, err := person.Load(id)
		if err != nil {
			errorHandler(rw, req, err)
			return
		}

		if p == nil {
			WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("person [%s] not found", id))
			clientError++
			return
		}

		if p.Player == false {
			WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("person [%s] is not a player", id))
			clientError++
			return
		}
	}

	id, err := r.Court.Insert()
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(CreateCourtResponse{
		ID: id,
	})
}
