package httphandler

import (
	"context"
	"database/sql"
	"fmt"
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
	filters[""] = ``
	filters["all"] = ``
	filters["players"] = `WHERE status = 'player'`
	filters["inactive"] = `WHERE status = 'inactive'`
	filters["suspended"] = `WHERE status = 'suspended'`
}

// ListPeople method
func ListPeople(writer http.ResponseWriter, request *http.Request) {
	f := functionListPeople

	_, err := checkAuthenticated(request)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	filter := request.URL.Query().Get("filter")

	var whereClause string
	var ok bool
	if whereClause, ok = filters[filter]; !ok {
		message := fmt.Sprintf("unexpected filter name: '%s'", filter)
		writeResponseMessage(writer, request, http.StatusBadRequest, message)
		return
	}

	object := request.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}

	list, err := model.ListPeople(context.Background(), db, whereClause)
	if err != nil {
		message := "problem listing people"
		DumpError(f, request, err, message)
		writeResponseError(writer, request, err)
		return
	}

	var listOfPeople []model.Person
	for _, person := range list {
		listOfPeople = append(listOfPeople, *person.ToLimited())
	}

	writeResponseObject(writer, request, http.StatusOK, listOfPeople)
}
