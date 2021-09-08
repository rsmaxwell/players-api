package httphandler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/rsmaxwell/players-api/internal/model"

	_ "github.com/jackc/pgx/stdlib"
)

func TestMakePlayerPlay(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	ctx := context.Background()

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie, accessToken := GetSigninToken(t, db, model.GoodEmail, model.GoodPassword)
	anotherPerson, _ := model.FindPersonByEmail(ctx, db, model.AnotherEmail)
	goodCourt := GetFirstCourt(t, db)

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName               string
		setLogonCookie         bool
		logonCookie            *http.Cookie
		setAuthorizationHeader bool
		accessToken            string
		personID               int
		courtID                int
		position               int
		expectedStatus         int
	}{
		{
			testName:               "Good request",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			personID:               anotherPerson.ID,
			courtID:                goodCourt.ID,
			position:               1,
			expectedStatus:         http.StatusOK,
		},
		{
			testName:               "Bad personID",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			personID:               999999999,
			courtID:                goodCourt.ID,
			position:               1,
			expectedStatus:         http.StatusNotFound,
		},
		{
			testName:               "Bad courtID",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			personID:               anotherPerson.ID,
			courtID:                999999999,
			position:               1,
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

			command := fmt.Sprintf("/people/toplaying/%d/%d/%d", test.personID, test.courtID, test.position)
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
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			if w.Code == http.StatusOK {
				// Check the person is playing
				listOfPlayers, err := model.ListPlayersForPerson(ctx, db, test.personID)
				require.Nil(t, err, "err should be nothing")
				require.Equal(t, 1, len(listOfPlayers), "Unexpected number of players for person: %d", len(listOfPlayers))

				listOfWaiters, err := model.ListWaitersForPerson(ctx, db, test.personID)
				require.Nil(t, err, "err should be nothing")
				require.Equal(t, 0, len(listOfWaiters), "person is still waiting")
			}

			if w.Code != test.expectedStatus {
				require.FailNow(t, "Unexpected status: expected: %d, actual: %d", test.expectedStatus, w.Code)
			}
		})
	}
}
