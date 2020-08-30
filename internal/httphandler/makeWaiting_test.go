package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/stdlib"
)

func TestMakeWaiting(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie := GetLoginToken(t, db, model.GoodUserName, model.GoodPassword)
	anotherPerson := FindPersonByUserName(t, db, model.AnotherUserName)

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName       string
		setLogonCookie bool
		logonCookie    *http.Cookie
		id             int
		expectedStatus int
	}{
		{
			testName:       "Good request",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             anotherPerson.ID,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad userID",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             999999999,
			expectedStatus: http.StatusBadRequest,
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Set up the handlers on the router
			router := mux.NewRouter()
			SetupHandlers(router)
			w := httptest.NewRecorder()

			// Create a request
			requestBody, err := json.Marshal(MakeWaitingRequest{
				ID: test.id,
			})
			require.Nil(t, err, "err should be nothing")

			command := fmt.Sprintf("/users/towaiting/%d", test.id)
			r, err := http.NewRequest("PUT", contextPath+command, bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			// ---------------------------------------

			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60*time.Second))
			defer cancel()
			r2 := r.WithContext(ctx)

			ctx = context.WithValue(r2.Context(), ContextDatabaseKey, db)
			r3 := r.WithContext(ctx)

			// ---------------------------------------

			// Serve the request
			router.ServeHTTP(w, r3)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			if w.Code == http.StatusOK {
				// Check the person is inactive
				listOfPlayers, err := model.ListPlayersForPerson(db, test.id)
				require.Nil(t, err, "err should be nothing")
				require.Equal(t, 0, len(listOfPlayers), "person is still playing")

				listOfWaiters, err := model.ListWaitersForPerson(db, test.id)
				require.Nil(t, err, "err should be nothing")
				require.Equal(t, 1, len(listOfWaiters), "Unexpected number of waiters for person: %d", len(listOfWaiters))
			}

			if w.Code != test.expectedStatus {
				require.FailNow(t, "Unexpected status: expected: %d, actual: %d", test.expectedStatus, w.Code)
			}
		})
	}
}
