package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
)

func TestListCourts(t *testing.T) {

	teardown := model.SetupFull(t)
	defer teardown(t)

	// ***************************************************************
	// * Login to get a valid token
	// ***************************************************************
	goodToken, err := getLoginToken(t, goodUserID, goodPassword)
	require.Nil(t, err, "err should be nothing")

	allCourtIDs, err := court.List()
	require.Nil(t, err, "err should be nothing")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName       string
		token          string
		expectedStatus int
		expectedResult []string
	}{
		{
			testName:       "Good request",
			token:          goodToken,
			expectedStatus: http.StatusOK,
			expectedResult: allCourtIDs,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: []string{},
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(ListCourtsRequest{
				Token: test.token,
			})
			require.Nil(t, err, "err should be nothing")

			// Create a  request to pass to our handler.
			req, err := http.NewRequest("GET", contextPath+"/court", bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			// Pass the request to our handler
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			// Check the response
			if rw.Code == http.StatusOK {
				bytes, err := ioutil.ReadAll(rw.Body)
				require.Nil(t, err, "err should be nothing")

				var response ListCourtsResponse
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")

				// Check the response body is what we expect.
				if !common.EqualArrayOfStrings(response.Courts, test.expectedResult) {
					t.Logf("actual:   %v", response.Courts)
					t.Logf("expected: %v", test.expectedResult)
					t.Errorf("Unexpected list of courts")
				}
			}
		})
	}
}
