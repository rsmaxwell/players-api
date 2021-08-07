package httphandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// UpdateCourtRequest structure
type UpdateCourtRequest struct {
	Court map[string]interface{} `json:"court"`
}

var (
	functionUpdateCourt = debug.NewFunction(pkg, "UpdateCourt")
)

// UpdateCourt method
func UpdateCourt(w http.ResponseWriter, r *http.Request) {
	f := functionUpdateCourt
	f.DebugAPI("")

	userID, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	user := model.FullPerson{ID: userID}
	err = user.LoadPerson(db)
	if err != nil {
		message := fmt.Sprintf("Could not load person [%d]", userID)
		f.DumpError(err, message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}
	err = user.CanEditCourt()
	if err != nil {
		f.DebugVerbose(fmt.Sprintf("Person [%d] is not allowed to edit court", userID))
		writeResponseMessage(w, r, http.StatusForbidden, "Forbidden")
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request UpdateCourtRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	str := mux.Vars(r)["id"]
	id, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}

	c := model.Court{ID: id}
	err = c.LoadCourt(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	if val, ok := request.Court["name"]; ok {
		c.Name, ok = val.(string)
		if !ok {
			err = fmt.Errorf("unexpected type for 'name'")
			writeResponseError(w, r, err)
			return
		}
	}

	err = c.UpdateCourt(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "ok")
}
