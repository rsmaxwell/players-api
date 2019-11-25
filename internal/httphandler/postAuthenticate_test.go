package httphandler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"

	"github.com/gorilla/mux"
)

func TestAuthenticate(t *testing.T) {

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
			role:           person.RoleAdmin,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Space in userID",
			userID:         "0 7",
			password:       "topsecret",
			role:           person.RoleNormal,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Path in userID",
			userID:         "../007",
			password:       "topsecret",
			role:           person.RoleSuspended,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Create a request to pass to our handler.
			req, err := http.NewRequest("POST", contextPath+"/users/authenticate", nil)
			require.Nil(t, err, "err should be nothing")

			req.Header.Set("Authorization", model.BasicAuth(test.userID, test.password))

			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			if rw.Code == http.StatusOK {
				_, err := ioutil.ReadAll(rw.Body)
				if err != nil {
					log.Fatalln(err)
				}
			}
		})
	}
}
