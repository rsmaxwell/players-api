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

	sess, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	userID, ok := sess.Values["userID"].(int)
	if !ok {
		f.DebugVerbose("could not get 'userID' from the session")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		err = fmt.Errorf("Unexpected context type")
		writeResponseError(w, r, err)
		return
	}

	user := model.Person{ID: userID}
	err = user.LoadPerson(db)
	if err != nil {
		f.DebugVerbose("Could not load person")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request UpdatePersonRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	str := mux.Vars(r)["id"]
	id, err := strconv.Atoi(str)
	if err != nil {
		f.DebugVerbose("Count not convert '" + str + "' into an int")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	f.DebugVerbose("ID: %d", id)
	f.DebugVerbose("Person: %v", request.Person)

	if userID == id {
		err = user.CanEditSelf()
		if err != nil {
			f.DebugVerbose("Not allowed to edit self")
			writeResponseMessage(w, r, http.StatusForbidden, "", "Forbidden")
			return
		}
	} else {
		err = user.CanEditOtherPeople()
		if err != nil {
			f.DebugVerbose("Not allowed to edit other people")
			writeResponseMessage(w, r, http.StatusForbidden, "", "Forbidden")
			return
		}
	}

	var p model.Person
	p.ID = id
	err = p.LoadPerson(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	if val, ok := request.Person["firstname"]; ok {
		p.FirstName, ok = val.(string)
		if !ok {
			err = fmt.Errorf("Unexpected type for 'firstName'")
			writeResponseError(w, r, err)
			return
		}
	}

	if val, ok := request.Person["lastname"]; ok {
		p.LastName, ok = val.(string)
		if !ok {
			err = fmt.Errorf("Unexpected type for 'lastName'")
			writeResponseError(w, r, err)
			return
		}
	}

	if val, ok := request.Person["displayname"]; ok {
		p.DisplayName, ok = val.(string)
		if !ok {
			err = fmt.Errorf("Unexpected type for 'displayname'")
			writeResponseError(w, r, err)
			return
		}
	}

	if val, ok := request.Person["username"]; ok {
		p.UserName, ok = val.(string)
		if !ok {
			err = fmt.Errorf("Unexpected type for 'username'")
			writeResponseError(w, r, err)
			return
		}
	}

	if val, ok := request.Person["email"]; ok {
		p.Email, ok = val.(string)
		if !ok {
			err = fmt.Errorf("Unexpected type for 'email'")
			writeResponseError(w, r, err)
			return
		}
	}

	if val, ok := request.Person["phone"]; ok {
		p.Phone, ok = val.(string)
		if !ok {
			err = fmt.Errorf("Unexpected type for 'phone'")
			writeResponseError(w, r, err)
			return
		}
	}

	if val, ok := request.Person["password"]; ok {
		password, ok := val.(string)
		if !ok {
			err = fmt.Errorf("Unexpected type for 'email'")
			writeResponseError(w, r, err)
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
			err = fmt.Errorf("Unexpected type for 'status'")
			writeResponseError(w, r, err)
			return
		}
	}

	err = p.UpdatePerson(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
