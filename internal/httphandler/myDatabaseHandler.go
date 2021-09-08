package httphandler

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type MyDatabase struct {
	handler http.Handler
	db      *sql.DB
}

func AddDatabaseContext(handlerToWrap http.Handler, db *sql.DB) *MyDatabase {
	return &MyDatabase{handler: handlerToWrap, db: db}
}

func (h *MyDatabase) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx1, cancel := context.WithTimeout(r.Context(), time.Duration(60*time.Second))
	defer cancel()

	ctx2 := context.WithValue(ctx1, ContextDatabaseKey, h.db)
	r3 := r.WithContext(ctx2)

	h.handler.ServeHTTP(w, r3)
}
