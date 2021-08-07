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
	"golang.org/x/crypto/bcrypt"

	_ "github.com/jackc/pgx/stdlib"
)

func TestUpdatePerson(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	// ***************************************************************
	// * Login
	// ***************************************************************
	logonCookie, accessToken := GetSigninToken(t, db, model.GoodEmail, model.GoodPassword)
	goodPerson, _ := model.FindPersonByEmail(db, model.GoodEmail)

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
		person                 map[string]interface{}
		expectedStatus         int
	}{
		{
			testName:               "Good request",
			setLogonCookie:         true,
			logonCookie:            logonCookie,
			setAuthorizationHeader: true,
			accessToken:            accessToken,
			id:                     goodPerson.ID,
			person: map[string]interface{}{
				"firstname": goodPerson.FirstName,
				"lastname":  goodPerson.LastName,
				"email":     goodPerson.Email,
				"password":  model.GoodPassword,
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
			expectedStatus:         http.StatusNotFound,
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
			requestBody, err := json.Marshal(UpdatePersonRequest{
				Person: test.person,
			})
			require.Nil(t, err, "err should be nothing")

			command := fmt.Sprintf("/people/%d", test.id)
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

			// Check the person was actually updated
			if w.Code == http.StatusOK {
				var p model.FullPerson
				p.ID = test.id
				err := p.LoadPerson(db)
				require.Nil(t, err, "err should be nothing")

				if value, ok := test.person["firstname"]; ok {
					if firstName, ok := value.(string); ok {
						if p.FirstName != firstName {
							t.Errorf("Expected %s, actual: %s", firstName, p.FirstName)
						}
					}
				}

				if value, ok := test.person["lastname"]; ok {
					if lastName, ok := value.(string); ok {
						if p.LastName != lastName {
							t.Errorf("Expected %s, actual: %s", lastName, p.LastName)
						}
					}
				}

				if value, ok := test.person["email"]; ok {
					if email, ok := value.(string); ok {
						if p.Email != email {
							t.Errorf("Expected %s, actual: %s", email, p.Email)
						}
					}
				}

				if value, ok := test.person["password"]; ok {
					if password, ok := value.(string); ok {
						err = bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(password))
						if err != nil {
							t.Errorf("The password was not updated correctly")
						}
					}
				}
			}

			if w.Code != test.expectedStatus {
				t.Logf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus)
				t.FailNow()
			}
		})
	}
}
