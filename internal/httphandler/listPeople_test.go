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
	signonCookie, accessToken := GetSigninToken(t, db, model.GoodEmail, model.GoodPassword)

	// ***************************************************************
	// * Get a list of all the people
	// ***************************************************************
	//allPeople, err := model.FindPeople(db, "all")
	//require.Nil(t, err, "err should be nothing")

	// ***************************************************************
	// * Testcases
	// ***************************************************************

	tests := []struct {
		testName               string
		setSignonCookie        bool
		signonCookie           *http.Cookie
		setAuthorizationHeader bool
		accessToken            string
		filter                 string
		expectedStatus         int
		expectedResult         []model.Person
	}{
		// {
		// 	testName:               "Good request",
		// 	setSignonCookie:        true,
		// 	signonCookie:           signonCookie,
		// 	setAuthorizationHeader: true,
		// 	accessToken:            accessToken,
		// 	filter:                 "players",
		// 	expectedStatus:         http.StatusOK,
		// 	expectedResult:         allPeople,
		// },
		{
			testName:               "Bad request",
			setSignonCookie:        true,
			signonCookie:           signonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			filter:                 "junk",
			expectedStatus:         http.StatusBadRequest,
			expectedResult:         nil,
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
			requestBody, err := json.Marshal(ListPeopleRequest{Filter: test.filter})
			require.Nil(t, err, "err should be nothing")

			command := "/people"
			r, err := http.NewRequest("POST", contextPath+command, bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			if test.setSignonCookie {
				r.AddCookie(test.signonCookie)
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
			if w.Code == http.StatusOK {
				bytes, err := ioutil.ReadAll(w.Body)
				require.Nil(t, err, "err should be nothing")

				var response []model.Person
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")

				// Check the response body is what we expect.
				// if !model.EqualIntArray(response.People, test.expectedResult) {
				// 	t.Logf("actual:   %v", response.People)
				// 	t.Logf("expected: %v", test.expectedResult)
				// 	t.Errorf("Unexpected list of people")
				// }
			} else if w.Code == http.StatusBadRequest {
				bytes, err := ioutil.ReadAll(w.Body)
				require.Nil(t, err, "err should be nothing")

				var response MessageResponse
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")
				fmt.Printf(response.Message)
			}
		})
	}
}
