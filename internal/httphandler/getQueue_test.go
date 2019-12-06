package httphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"
)

func TestGetQueue(t *testing.T) {

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
		testName              string
		setLogonCookie        bool
		logonCookie           *http.Cookie
		userID                string
		expectedStatus        int
		expectedResultName    string
		expectedResultPlayers []string
	}{
		{
			testName:              "Good request",
			setLogonCookie:        true,
			logonCookie:           logonCookie,
			expectedStatus:        http.StatusOK,
			expectedResultName:    "Queue",
			expectedResultPlayers: []string{"one", "two"},
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
			r, err := http.NewRequest("GET", contextPath+"/queue", nil)
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			// Serve the request
			router.ServeHTTP(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			// Check the response
			bytes, err := ioutil.ReadAll(w.Body)
			require.Nil(t, err, "err should be nothing")

			if w.Code == http.StatusOK {
				var response GetQueueResponse
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")

				actualName := response.Queue.Container.Name
				require.Equal(t, test.expectedResultName, actualName, fmt.Sprintf("handler returned unexpected body: want %v, got %v", test.expectedResultName, actualName))

				actualPlayers := response.Queue.Container.Players
				if common.EqualArrayOfStrings(actualPlayers, test.expectedResultPlayers) {
					require.Fail(t, fmt.Sprintf("handler returned unexpected body: want %v, got %v", test.expectedResultPlayers, actualPlayers))
				}
			}
		})
	}
}
