package httphandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionMakePlaying = debug.NewFunction(pkg, "MakePlaying")
)

// MakePlayingRequest structure
type MakePlayingRequest struct {
	ID int `json:"id"`
}

// MakePlaying method
func MakePlaying(w http.ResponseWriter, r *http.Request) {
	f := functionMakePlaying

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	str := mux.Vars(r)["id1"]
	personID, err := strconv.Atoi(str)
	if err != nil {
		f.DebugVerbose("Count not convert '" + str + "' into an int")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	f.DebugVerbose("PersonID: %d", personID)

	str = mux.Vars(r)["id2"]
	courtID, err := strconv.Atoi(str)
	if err != nil {
		f.DebugVerbose("Count not convert '" + str + "' into an int")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	f.DebugVerbose("courtID: %d", courtID)

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		err = fmt.Errorf("unexpected context type")
		writeResponseError(w, r, err)
		return
	}

	p := model.Person{ID: personID}
	exists, err := p.PersonExists(db)
	if err != nil {
		message := fmt.Sprintf("Unexpected error checking person [%d] exists", personID)
		f.DumpError(err, message)
		writeResponseError(w, r, codeerror.NewInternalServerError(message))
		return
	}
	if !exists {
		message := fmt.Sprintf("person [%d] not found", personID)
		writeResponseError(w, r, codeerror.NewBadRequest(message))
		return
	}

	c := model.Court{ID: courtID}
	exists, err = c.CourtExists(db)
	if err != nil {
		message := fmt.Sprintf("Unexpected error checking court [%d] exists", courtID)
		f.DumpError(err, message)
		writeResponseError(w, r, codeerror.NewInternalServerError(message))
		return
	}
	if !exists {
		message := fmt.Sprintf("court [%d] not found", courtID)
		writeResponseError(w, r, codeerror.NewBadRequest(message))
		return
	}

	err = model.MakePlaying(db, personID, courtID)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
