package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/rsmaxwell/players-api/internal/model"

	_ "github.com/jackc/pgx/stdlib"
)

func TestRegister(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	ctx := context.Background()

	tests := []struct {
		testName       string
		registration   model.Registration
		expectedStatus int
	}{
		{
			testName: "Good request",
			registration: model.Registration{
				FirstName: "James",
				LastName:  "Bond",
				Knownas:   "aaa",
				Email:     "007@mi6.co.uk",
				Phone:     "012345 123456",
				Password:  "topsecret",
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "Space in email",
			registration: model.Registration{
				FirstName: "James",
				LastName:  "Bond",
				Knownas:   "aaa",
				Email:     "007 @mi6.co.uk",
				Phone:     "012345 123456",
				Password:  "topsecret",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName: "Bad email",
			registration: model.Registration{
				FirstName: "James",
				LastName:  "Bond",
				Knownas:   "aaa",
				Email:     "007mi6.co.uk",
				Phone:     "012345 123999",
				Password:  "topsecret",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			listOfPeople, err := model.ListPeople(ctx, db, "")
			require.Nil(t, err)
			initialNumberOfPeople := len(listOfPeople)

			// Set up the handlers on the router
			router := mux.NewRouter()
			SetupHandlers(router)
			w := httptest.NewRecorder()

			// Create a request
			requestBody, err := json.Marshal(test.registration)
			if err != nil {
				log.Fatalln(err)
			}

			r, err := http.NewRequest("POST", contextPath+"/register", bytes.NewBuffer(requestBody))
			require.Nil(t, err)

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

				listOfPeople, err = model.ListPeople(ctx, db, "")
				require.Nil(t, err)
				finalNumberOfPeople := len(listOfPeople)

				require.Equal(t, initialNumberOfPeople+1, finalNumberOfPeople, "Person was not registered")

				// Check the status of the new person
				p, err := model.FindPersonByEmail(ctx, db, test.registration.Email)
				require.Nil(t, err)
				if initialNumberOfPeople == 0 {
					require.Equal(t, model.StatusAdmin, p.Status, "Unexpected role")
				} else {
					require.Equal(t, model.StatusSuspended, p.Status, "Unexpected role")
				}
			}

			if w.Code != test.expectedStatus {
				t.Logf("handler returned wrong status code: expected:%v actual:%v", test.expectedStatus, w.Code)
				t.FailNow()
			}
		})
	}
}
