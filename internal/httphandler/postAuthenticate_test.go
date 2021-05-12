package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"

	"github.com/gorilla/mux"

	_ "github.com/jackc/pgx/stdlib"
)

func TestAuthenticate(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		email          string
		password       string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			email:          model.GoodUserName,
			password:       model.GoodPassword,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad userid",
			email:          "junk",
			password:       "junk",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Create a request to pass to our handler.
			command := fmt.Sprintf("/users/authenticate")
			r, err := http.NewRequest("POST", contextPath+command, nil)
			require.Nil(t, err, "err should be nothing")

			r.Header.Set("Authorization", BasicAuth(test.email, test.password))

			router := mux.NewRouter()
			SetupHandlers(router)
			w := httptest.NewRecorder()

			// ---------------------------------------

			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60*time.Second))
			defer cancel()
			r2 := r.WithContext(ctx)

			ctx = context.WithValue(r2.Context(), ContextDatabaseKey, db)
			r3 := r.WithContext(ctx)

			// ---------------------------------------

			router.ServeHTTP(w, r3)

			if w.Code == http.StatusOK {
				bytes, err := ioutil.ReadAll(w.Body)
				if err != nil {
					log.Fatalln(err)
				}

				var resp PostAuthenticateResponse
				err = json.Unmarshal(bytes, &resp)
				require.Nil(t, err, "err should be nothing")
				require.Equal(t, resp.Person.Email, test.email)
			}

			if w.Code != test.expectedStatus {
				t.Logf("Unexpected status: expected:%v actual:%v", test.expectedStatus, w.Code)
				t.FailNow()
			}
		})
	}
}
