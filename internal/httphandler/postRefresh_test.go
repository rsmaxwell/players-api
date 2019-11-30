package httphandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/rsmaxwell/players-api/internal/model"
)

func TestRefresh(t *testing.T) {

	teardown := model.SetupFull(t)
	defer teardown(t)

	// ***************************************************************
	// * Login to get tokens
	// ***************************************************************
	accessTokenString, refreshTokenCookie := testLogin(t, "007", "topsecret")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName            string
		setAccessToken      bool
		accessToken         string
		useGoodRefreshToken bool
		setRefreshToken     bool
		refreshToken        string
		expectedStatus      int
	}{
		{
			testName:            "Good request",
			setAccessToken:      true,
			accessToken:         "Bearer " + accessTokenString,
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			expectedStatus:      http.StatusOK,
		},
		{
			testName:            "junk accessToken",
			setAccessToken:      true,
			accessToken:         "Bearer " + "junk",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			expectedStatus:      http.StatusBadRequest,
		},
		{
			testName:            "no 'Bearer' prefix before accessToken",
			setAccessToken:      true,
			accessToken:         "junk",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			expectedStatus:      http.StatusBadRequest,
		},
		{
			testName:            "missing accessToken",
			setAccessToken:      false,
			accessToken:         "",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			expectedStatus:      http.StatusUnauthorized,
		},
		{
			testName:            "junk refreshToken",
			setAccessToken:      true,
			accessToken:         "Bearer " + "junk",
			useGoodRefreshToken: false,
			setRefreshToken:     true,
			refreshToken:        "junk",
			expectedStatus:      http.StatusBadRequest,
		},
		{
			testName:            "missing refreshToken",
			setAccessToken:      true,
			accessToken:         "Bearer " + accessTokenString,
			useGoodRefreshToken: false,
			setRefreshToken:     false,
			refreshToken:        "",
			expectedStatus:      http.StatusUnauthorized,
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

			rw := httptest.NewRecorder()

			req, err := http.NewRequest("POST", contextPath+"/users/refresh", nil)
			require.Nil(t, err)

			setAccessToken(req, test.setAccessToken, test.accessToken)
			setRefreshToken(req, test.useGoodRefreshToken, test.setRefreshToken, refreshTokenCookie, test.refreshToken)

			// Serve the request
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			if rw.Code == http.StatusOK {

				response := rw.Result()

				accessTokenArray2 := response.Header["Access-Token"]
				require.Equal(t, len(accessTokenArray2), 1, "accessToken array should have 1 entry: found "+strconv.Itoa(len(accessTokenArray2)))

				accessTokenString2 := accessTokenArray2[0]
				require.NotEmpty(t, accessTokenString2, "accessToken should not be empty")

				cookies := map[string]*http.Cookie{}
				for _, cookie := range response.Cookies() {
					cookies[cookie.Name] = cookie
				}

				refreshTokenCookie2 := cookies["players-api"]
				require.NotEmpty(t, refreshTokenCookie2, "expecting an refresh token")

				refreshTokenString2 := refreshTokenCookie2.Value
				require.NotEmpty(t, refreshTokenString2, "refreshToken should not be empty")
			}
		})
	}
}
