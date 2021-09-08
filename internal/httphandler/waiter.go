package httphandler

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionListWaiters = debug.NewFunction(pkg, "ListWaiters")
)

// DisplayWaiter type
type DisplayWaiter struct {
	PersonID int    `json:"personID"`
	Knownas  string `json:"knownas"`
	Start    int64  `json:"start"`
}

// ListWaiters method
func ListWaiters(writer http.ResponseWriter, request *http.Request) {
	f := functionListWaiters

	_, err := checkAuthenticated(request)
	if err != nil {
		writeResponseError(writer, request, err)
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

	waiters, err := model.ListWaiters(context.Background(), db)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	var list []DisplayWaiter
	for _, waiter := range waiters {

		p := model.FullPerson{ID: waiter.Person}
		err := p.LoadPerson(context.Background(), db)
		if err != nil {
			writeResponseError(writer, request, err)
			return
		}

		w := DisplayWaiter{}
		w.PersonID = waiter.Person
		w.Knownas = p.Knownas
		w.Start = waiter.Start.Unix()

		list = append(list, w)
	}

	writeResponseObject(writer, request, http.StatusOK, list)
}
