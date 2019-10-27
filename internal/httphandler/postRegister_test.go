package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/model"
)

func TestRegister(t *testing.T) {

	teardown := model.SetupEmpty(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		userID         string
		password       string
		firstName      string
		lastName       string
		email          string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			userID:         "007",
			password:       "topsecret",
			firstName:      "James",
			lastName:       "Bond",
			email:          "james@mi6.co.uk",
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Space in userID",
			userID:         "0 7",
			password:       "topsecret",
			firstName:      "James",
			lastName:       "Bond",
			email:          "james@mi6.co.uk",
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName:       "Path in userID",
			userID:         "../007",
			password:       "topsecret",
			firstName:      "James",
			lastName:       "Bond",
			email:          "james@mi6.co.uk",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			initialNumberOfPeople, err := person.Size()
			require.Nil(t, err)

			// Set up the handlers on the router
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()

			// Create a request
			requestBody, err := json.Marshal(RegisterRequest{
				UserID:    test.userID,
				Password:  test.password,
				FirstName: test.firstName,
				LastName:  test.lastName,
				Email:     test.email,
			})
			if err != nil {
				log.Fatalln(err)
			}

			req, err := http.NewRequest("POST", contextPath+"/register", bytes.NewBuffer(requestBody))
			require.Nil(t, err)

			// Serve the request
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			// Check the response
			if rw.Code == http.StatusOK {
				finalNumberOfPeople, err := person.Size()
				require.Nil(t, err)
				require.Equal(t, initialNumberOfPeople+1, finalNumberOfPeople, "Person was not registered")

				// Check the status of the new person
				p, err := person.Load(test.userID)
				require.Nil(t, err)
				if initialNumberOfPeople == 0 {
					require.Equal(t, person.RoleAdmin, p.Role, "Unexpected role")
				} else {
					require.Equal(t, person.RoleSuspended, p.Role, "Unexpected role")
				}
			}
		})
	}
}
