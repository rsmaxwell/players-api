package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/stdlib"
)

func TestUpdateCourt(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie, accessToken := GetSigninToken(t, db, model.GoodEmail, model.GoodPassword)
	goodCourt := GetFirstCourt(t, db)

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName               string
		setLogonCookie         bool
		logonCookie            *http.Cookie
		setAuthorizationHeader bool
		accessToken            string
		id                     int
		court                  map[string]interface{}
		expectedStatus         int
	}{
		{
			testName:               "Good request",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			id:                     goodCourt.ID,
			court: map[string]interface{}{
				"name": "COURT 101",
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName:               "Bad userID",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			id:                     999999999,
			court: map[string]interface{}{
				"name": "COURT 101",
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

			command := fmt.Sprintf("/courts/%d", test.id)
			r, err := http.NewRequest("PUT", contextPath+command, bytes.NewBuffer(requestBody))
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

			// Check the response
			if w.Code == http.StatusOK {
				var c model.Court
				c.ID = test.id
				err := c.LoadCourt(db)
				require.Nil(t, err, "err should be nothing")

				if value, ok := test.court["name"]; ok {
					if name, ok := value.(string); ok {
						if c.Name != name {
							t.Errorf("Expected %s, actual: %s", name, c.Name)
						}
					}
				}
			}

			if w.Code != test.expectedStatus {
				t.Errorf("Unexpectred status. Expected %d, actual: %d", test.expectedStatus, w.Code)
				t.FailNow()
			}
		})
	}
}
