package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/commands"
	"github.com/rsmaxwell/players-api/destination"
	"github.com/stretchr/testify/require"
)

func TestPostMove(t *testing.T) {

	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		source         destination.Reference
		target         destination.Reference
		players        []string
		expectedStatus int
	}{
		{
			testName: "Good request",
			token:    MyToken,
			source: destination.Reference{
				Type: "queue",
				ID:   "",
			},
			target: destination.Reference{
				Type: "court",
				ID:   "1000",
			},
			players:        []string{"007", "bob", "john"},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "Bad token",
			token:    "junk",
			source: destination.Reference{
				Type: "queue",
				ID:   "",
			},
			target: destination.Reference{
				Type: "queue",
				ID:   "1001",
			},
			players:        []string{"007", "bob", "john"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName: "Bad player",
			token:    MyToken,
			source: destination.Reference{
				Type: "queue",
				ID:   "",
			},
			target: destination.Reference{
				Type: "court",
				ID:   "1000",
			},
			players:        []string{"007", "junk", "john"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(PostMoveRequest{
				Token:   test.token,
				Source:  test.source,
				Target:  test.target,
				Players: test.players,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("POST", "/move", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			router := mux.NewRouter()
			SetupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rw := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can call its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rw, req)

			// Check the status code is what we expect.
			if rw.Code != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus)
			}

			if rw.Code == http.StatusOK {

				// Check the moved people have actually moved
				for _, personID := range test.players {

					ref, err := findPlayer(personID)
					require.Nil(t, err, fmt.Sprintf("error: %s", err))
					require.NotNil(t, ref, fmt.Sprintf("person[%s] not found", personID))

					// Check the moved person is NOT at the source
					found := commands.EqualsContainerReference(ref, &test.source)
					require.False(t, found, fmt.Sprintf("person[%s] is still at the source: %s", personID, destination.FormatReference(&test.source)))

					// Check the moved person IS at the target
					found = commands.EqualsContainerReference(ref, &test.target)
					require.True(t, found, fmt.Sprintf("person[%s] is not at the target: %s", personID, destination.FormatReference(&test.target)))
				}
			}
		})
	}
}

func findPlayer(id string) (*destination.Reference, error) {

	courts, err := destination.ListCourts()
	if err != nil {
		return nil, err
	}

	// Look for the player on one of the courts
	for _, courtID := range courts {

		c, err := destination.LoadCourt(courtID)
		if err != nil {
			return nil, err
		}

		for _, personID := range c.Container.Players {

			if id == personID {
				ref := destination.Reference{Type: "court", ID: courtID}
				return &ref, nil
			}
		}
	}

	// Look for the player on the queue
	q, err := destination.LoadQueue()
	if err != nil {
		return nil, err
	}

	for _, personID := range q.Container.Players {

		if id == personID {
			ref := destination.Reference{Type: "queue", ID: ""}
			return &ref, nil
		}
	}

	return nil, nil
}
