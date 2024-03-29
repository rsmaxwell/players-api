package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/stdlib"
)

func TestGetPerson(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	ctx := context.Background()

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie, accessToken := GetSigninToken(t, db, model.GoodEmail, model.GoodPassword)
	goodPerson, _ := model.FindPersonByEmail(ctx, db, model.GoodEmail)

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
		expectedResultName     string
	}{
		{
			testName:               "Good request",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			userID:                 goodPerson.ID,
			expectedStatus:         http.StatusOK,
			expectedResultName:     model.GoodFirstName,
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
			r, err := http.NewRequest("GET", contextPath+command, nil)
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

			// Check the response
			bytes, err := ioutil.ReadAll(w.Body)
			require.Nil(t, err, "err should be nothing")

			if w.Code == http.StatusOK {
				var person model.Person
				err = json.Unmarshal(bytes, &person)
				require.Nil(t, err, "err should be nothing")
				require.Equal(t, test.expectedResultName, person.FirstName, fmt.Sprintf("handler returned unexpected body: want %v, got %v", test.expectedResultName, person.FirstName))
			}
		})
	}
}
