package httphandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// UpdatePersonRequest structure
type UpdatePersonRequest struct {
	Person map[string]interface{} `json:"person"`
}

var (
	functionUpdatePerson = debug.NewFunction(pkg, "UpdatePerson")
)

// UpdatePerson method
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	f := functionUpdatePerson
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
		f.DebugVerbose(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request UpdatePersonRequest
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

	f.DebugVerbose("ID: %d", id)
	f.DebugVerbose("Person: %v", request.Person)

	if userID == id {
		err = user.CanEditSelf()
		if err != nil {
			f.DebugVerbose("Not allowed to edit self")
			writeResponseMessage(w, r, http.StatusForbidden, "Forbidden")
			return
		}
	} else {
		err = user.CanEditOtherPeople()
		if err != nil {
			f.DebugVerbose("Not allowed to edit other people")
			writeResponseMessage(w, r, http.StatusForbidden, "Forbidden")
			return
		}
	}

	var p model.FullPerson
	p.ID = id
	err = p.LoadPerson(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	if val, ok := request.Person["firstname"]; ok {
		p.FirstName, ok = val.(string)
		if !ok {
			writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("unexpected type for [%s]: %v", "firstName", val))
			return
		}
	}

	if val, ok := request.Person["lastname"]; ok {
		p.LastName, ok = val.(string)
		if !ok {
			writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("unexpected type for [%s]: %v", "lastName", val))
			return
		}
	}

	if val, ok := request.Person["knownas"]; ok {
		p.Knownas, ok = val.(string)
		if !ok {
			writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("unexpected type for [%s]: %v", "knownas", val))
			return
		}
	}

	if val, ok := request.Person["email"]; ok {
		p.Email, ok = val.(string)
		if !ok {
			writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("unexpected type for [%s]: %v", "email", val))
			return
		}
	}

	if val, ok := request.Person["phone"]; ok {
		p.Phone, ok = val.(string)
		if !ok {
			writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("unexpected type for [%s]: %v", "phone", val))
			return
		}
	}

	if val, ok := request.Person["password"]; ok {
		password, ok := val.(string)
		if !ok {
			writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("unexpected type for [%s]: %v", "password", val))
			return
		}

		p.Hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			writeResponseError(w, r, err)
			return
		}
	}

	if val, ok := request.Person["status"]; ok {
		p.Status, ok = val.(string)
		if !ok {
			writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("unexpected type for [%s]: %v", "status", val))
			return
		}
	}

	err = p.UpdatePerson(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "ok")
}
