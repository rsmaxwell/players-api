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

func TestGetCourt(t *testing.T) {
	
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
		testName       string
		token          string
		courtID        string
		expectedStatus int
		expectedResult string
	}{
		{
			testName:       "Good request",
			token:          goodToken,
			courtID:        goodCourtID,
			expectedStatus: http.StatusOK,
			expectedResult: "Court 1",
		},
		{
			testName:       "Bad token",
			token:          "junk",
			courtID:        goodCourtID,
			expectedStatus: http.StatusUnauthorized,
			expectedResult: "",
		},
		{
			testName:       "Bad userID",
			token:          goodToken,
			courtID:        "junk",
			expectedStatus: http.StatusNotFound,
			expectedResult: "",
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(GetCourtRequest{
				Token: test.token,
			})
			require.Nil(t, err, "err should be nothing")

			// Create a request
			req, err := http.NewRequest("GET", contextPath+"/court/"+test.courtID, bytes.NewBuffer(requestBody))
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

			var response GetCourtResponse
			err = json.Unmarshal(bytes, &response)
			require.Nil(t, err, "err should be nothing")

			actual := response.Court.Container.Name
			if actual != test.expectedResult {
				require.Fail(t, fmt.Sprintf("handler returned unexpected body: got %v want %v", actual, test.expectedResult))
			}
		})
	}
}
