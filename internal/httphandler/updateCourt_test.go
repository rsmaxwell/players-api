package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateCourt(t *testing.T) {

	teardown := model.SetupFull(t)
	defer teardown(t)

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie := testLogin(t, "007", "topsecret")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName       string
		setLogonCookie bool
		logonCookie    *http.Cookie
		id             string
		court          map[string]interface{}
		expectedStatus int
	}{
		{
			testName:       "Good request",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             goodCourtID,
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{"bob", "jill", "alice"},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad userID",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             "junk",
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{},
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			testName:       "Bad player",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             goodCourtID,
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{"junk"},
				},
			},
			expectedStatus: http.StatusNotFound,
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
			requestBody, err := json.Marshal(UpdateCourtRequest{
				Court: test.court,
			})
			require.Nil(t, err, "err should be nothing")

			r, err := http.NewRequest("PUT", contextPath+"/court/"+test.id, bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			// Serve the request
			router.ServeHTTP(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			// Check the response
			if w.Code == http.StatusOK {
				ref := common.Reference{Type: "court", ID: test.id}
				c, err := court.Load(&ref)
				require.Nil(t, err, "err should be nothing")

				if i, ok := test.court["Name"]; ok {
					value, ok := i.(string)
					if !ok {
						t.Errorf("The type of 'test.court[\"Name\"]' should be a string")
					}
					c.Container.Name = value
					assert.Equal(t, c.Container.Name, value, "The Court name was not updated correctly")
				}

				if i, ok := test.court["Players"]; ok {
					value, ok := i.([]string)
					if !ok {
						t.Errorf("The type of 'test.court[\"Players\"]' should be an array of strings")
					}
					c.Container.Players = value
					if !common.EqualArrayOfStrings(c.Container.Players, value) {
						t.Errorf("The Court name was not updated correctly:\n got %v\n want %v", c.Container.Players, value)
					}
				}
			}
		})
	}
}
