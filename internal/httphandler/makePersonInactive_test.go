package httphandler

import (
	"context"
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

func TestMakeInactive(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie, accessToken := GetSigninToken(t, db, model.GoodEmail, model.GoodPassword)
	anotherPerson, _ := model.FindPersonByEmail(db, model.AnotherEmail)

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName               string
		setLogonCookie         bool
		logonCookie            *http.Cookie
		setAuthorizationHeader bool
		accessToken            string
		id                     int
		expectedStatus         int
	}{
		{
			testName:               "Good request",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			id:                     anotherPerson.ID,
			expectedStatus:         http.StatusOK,
		},
		{
			testName:               "Bad userID",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			id:                     999999999,
			expectedStatus:         http.StatusNotFound,
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

			command := fmt.Sprintf("/people/toinactive/%d", test.id)
			r, err := http.NewRequest("PUT", contextPath+command, nil)
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			if test.setAuthorizationHeader {
				r.Header.Set("Authorization", "Bearer "+test.accessToken)
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

			if w.Code == http.StatusOK {
				listOfPlayers, err := model.ListPlayersForPerson(db, test.id)
				require.Nil(t, err)
				require.Zero(t, len(listOfPlayers))

				listOfWaiters, err := model.ListWaitersForPerson(db, test.id)
				require.Nil(t, err)
				require.Zero(t, len(listOfWaiters))
			}

			if w.Code != test.expectedStatus {
				require.FailNow(t, "Unexpected status: expected: %d, actual: %d", test.expectedStatus, w.Code)
			}
		})
	}
}
