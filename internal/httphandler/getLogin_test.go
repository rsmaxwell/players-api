package httphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/session"
	"github.com/stretchr/testify/require"

	"github.com/gorilla/mux"
)

func TestLogin(t *testing.T) {

	teardown := model.SetupOne(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		userID         string
		password       string
		role           string
		expectedStatus int
	}{
		{
			testName:       "Good request from admin",
			userID:         "007",
			password:       "topsecret",
			role:           model.RoleAdmin,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Space in userID",
			userID:         "0 7",
			password:       "topsecret",
			role:           model.RoleNormal,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Path in userID",
			userID:         "../007",
			password:       "topsecret",
			role:           model.RoleSuspended,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", "/login", nil)
			require.Nil(t, err, "err should be nothing")

			req.Header.Set("Authorization", model.BasicAuth(test.userID, test.password))

			router := mux.NewRouter()
			SetupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rw := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can call its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rw, req)

			// Check the status code is what we expect.
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			if rw.Code == http.StatusOK {
				// Check the response is what we expect.
				bytes, err := ioutil.ReadAll(rw.Body)
				if err != nil {
					log.Fatalln(err)
				}

				var response LogonResponse

				err = json.Unmarshal(bytes, &response)
				if err != nil {
					log.Fatalln(err)
				}

				// Check the response contains a valid token.
				session := session.LookupToken(response.Token)
				if session == nil {
					t.Errorf("Invalid token returned: token:%s", response.Token)
				}
			}
		})
	}
}
