package httphandler

import (
	"bytes"
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

func TestListPeople(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie := GetLoginToken(t, db, model.GoodUserName, model.GoodPassword)

	// ***************************************************************
	// * Get a list of all the people
	// ***************************************************************
	allPeopleIDs, err := model.ListPeople(db, nil)
	require.Nil(t, err, "err should be nothing")

	// ***************************************************************
	// * Testcases
	// ***************************************************************

	tests := []struct {
		testName       string
		setLogonCookie bool
		logonCookie    *http.Cookie
		query          model.Query
		expectedStatus int
		expectedResult []int
	}{
		{
			testName:       "Good request",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			query:          map[string]model.Condition{"status": {Operation: "<>", Value: "suspended"}},
			expectedStatus: http.StatusOK,
			expectedResult: allPeopleIDs,
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
			requestBody, err := json.Marshal(ListPeopleRequest{Query: test.query})
			require.Nil(t, err, "err should be nothing")

			requestString := string(requestBody)
			t.Errorf("requestString: %s", requestString)

			command := fmt.Sprintf("/users")
			r, err := http.NewRequest("GET", contextPath+command, bytes.NewBuffer(requestBody))
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

			// Check the response
			if w.Code == http.StatusOK {
				bytes, err := ioutil.ReadAll(w.Body)
				require.Nil(t, err, "err should be nothing")

				var response ListPeopleResponse
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")

				// Check the response body is what we expect.
				if !model.EqualIntArray(response.People, test.expectedResult) {
					t.Logf("actual:   %v", response.People)
					t.Logf("expected: %v", test.expectedResult)
					t.Errorf("Unexpected list of people")
				}
			}
		})
	}
}
