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

func TestListCourts(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie, accessToken := GetSigninToken(t, db, model.GoodEmail, model.GoodPassword)

	// ***************************************************************
	// * Get the list of all courts
	// ***************************************************************
	// allCourtIDs, err := model.ListCourts(db)
	// require.Nil(t, err, "err should be nothing")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName               string
		setLogonCookie         bool
		logonCookie            *http.Cookie
		setAuthorizationHeader bool
		accessToken            string
		expectedStatus         int
		expectedResult         []int
	}{
		{
			testName:               "Good request",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			expectedStatus:         http.StatusOK,
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
			command := "/courts"
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
			// if w.Code == http.StatusOK {
			// bytes, err := ioutil.ReadAll(w.Body)
			// require.Nil(t, err, "err should be nothing")

			// var response ListCourtsResponse
			// err = json.Unmarshal(bytes, &response)
			// require.Nil(t, err, "err should be nothing")

			// Check the response body is what we expect.
			// if !model.EqualIntArray(response.Courts, test.expectedResult) {
			// 	t.Logf("actual:   %v", response.Courts)
			// 	t.Logf("expected: %v", test.expectedResult)
			// 	t.Errorf("Unexpected list of courts")
			// }
			// }
		})
	}
}
