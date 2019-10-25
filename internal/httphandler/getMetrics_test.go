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

	"github.com/rsmaxwell/players-api/internal/model"
)

func TestGetMetrics(t *testing.T) {

	teardown := model.SetupFull(t)
	defer teardown(t)

	// ***************************************************************
	// * Login to get a valid token
	// ***************************************************************
	goodToken, err := getLoginToken(t, goodUserID, goodPassword)
	require.Nil(t, err, "err should be nothing")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName              string
		token                 string
		expectedStatus        int
		expectedClientSuccess int
	}{
		{
			testName:              "Good request",
			token:                 goodToken,
			expectedStatus:        http.StatusOK,
			expectedClientSuccess: 0,
		},
		{
			testName:              "Bad token",
			token:                 "junk",
			expectedStatus:        http.StatusUnauthorized,
			expectedClientSuccess: 0,
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(GetMetricsRequest{
				Token: test.token,
			})
			require.Nil(t, err, "err should be nothing")

			// Create a request
			req, err := http.NewRequest("GET", contextPath+"/metrics", bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			// Pass the request to our handler
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			// Check the response
			bytes, err := ioutil.ReadAll(rw.Body)
			require.Nil(t, err, "err should be nothing")

			if rw.Code == http.StatusOK {
				var response GetMetricsResponse
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")
				require.Equal(t, test.expectedClientSuccess, response.Data.ClientSuccess, fmt.Sprintf("Unexpected metrics: expected: %v, actual:   %v", test.expectedClientSuccess, response.Data.ClientSuccess))
			}
		})
	}
}
