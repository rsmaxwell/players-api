package httphandler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/stdlib"
)

func TestDeletePerson(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie, accessToken := GetSigninToken(t, db, model.GoodEmail, model.GoodPassword)
	goodPerson, _ := model.FindPersonByEmail(db, model.GoodEmail)
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
		userID                 int
		expectedStatus         int
	}{
		{
			testName:               "Good request",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			userID:                 anotherPerson.ID,
			expectedStatus:         http.StatusOK,
		},
		{
			testName:               "Bad userID",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			userID:                 999999999,
			expectedStatus:         http.StatusOK,
		},
		{
			testName:               "delete myself",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			userID:                 goodPerson.ID,
			expectedStatus:         http.StatusUnauthorized,
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
			command := fmt.Sprintf("/people/%d", test.userID)
			r, err := http.NewRequest("DELETE", contextPath+command, nil)
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
				p := model.FullPerson{ID: test.userID}
				err := p.LoadPerson(db)
				require.NotNil(t, err, "Person was not deleted")
			}
		})
	}
}
