package httphandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// ListPeopleRequest structure
type ListPeopleRequest struct {
	Filter string `json:"filter"`
}

// ListPeopleResponse structure
// type ListPeopleResponse struct {
// 	Message string                `json:"message"`
// 	People  []model.LimitedPerson `json:"people"`
// }

var (
	functionListPeople = debug.NewFunction(pkg, "ListPeople")

	filters = make(map[string]string)
)

func init() {
	filters["all"] = ``
	filters["players"] = `WHERE status = 'player'`
	filters["inactive"] = `WHERE status = 'inactive'`
	filters["suspended"] = `WHERE status = 'suspended'`
}

// ListPeople method
func ListPeople(w http.ResponseWriter, r *http.Request) {
	f := functionListPeople
	f.DebugAPI("")

	if r.Method == http.MethodOptions {
		writeResponseMessage(w, r, http.StatusOK, "ok")
		return
	}

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request ListPeopleRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var whereClause string
	var ok bool
	if whereClause, ok = filters[request.Filter]; !ok {
		message := fmt.Sprintf("unexpected filter name: '%s'", request.Filter)
		writeResponseMessage(w, r, http.StatusBadRequest, message)
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

	list, err := model.ListPeople(db, whereClause)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	var listOfPeople []model.Person
	for _, person := range list {
		listOfPeople = append(listOfPeople, *person.ToLimited())
	}

	writeResponseObject(w, r, http.StatusOK, listOfPeople)
}
