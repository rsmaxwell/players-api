package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/basic/destination"
	"github.com/rsmaxwell/players-api/internal/basic/queue"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
)

func TestPostMove(t *testing.T) {

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
		testName       string
		setLogonCookie bool
		logonCookie    *http.Cookie
		source         common.Reference
		target         common.Reference
		players        []string
		expectedStatus int
	}{
		{
			testName: "Good request",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			source: common.Reference{
				Type: "queue",
				ID:   "",
			},
			target: common.Reference{
				Type: "court",
				ID:   "1000",
			},
			players:        []string{"007", "bob", "john"},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "Bad player",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			source: common.Reference{
				Type: "queue",
				ID:   "",
			},
			target: common.Reference{
				Type: "court",
				ID:   "1000",
			},
			players:        []string{"007", "junk", "john"},
			expectedStatus: http.StatusBadRequest,
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
			requestBody, err := json.Marshal(PostMoveRequest{
				Source:  test.source,
				Target:  test.target,
				Players: test.players,
			})
			require.Nil(t, err, "err should be nothing")

			r, err := http.NewRequest("POST", contextPath+"/users/move", bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			// Serve the request
			router.ServeHTTP(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			// Check the response
			if w.Code == http.StatusOK {

				// Check the moved people have actually moved
				for _, personID := range test.players {

					ref, err := findPlayer(t, personID)
					require.Nil(t, err, fmt.Sprintf("error: %s", err))
					require.NotNil(t, ref, fmt.Sprintf("person[%s] not found", personID))

					// Check the moved person is NOT at the source
					found := model.EqualsContainerReference(ref, &test.source)
					require.False(t, found, fmt.Sprintf("person[%s] is still at the source: %s", personID, destination.FormatReference(&test.source)))

					// Check the moved person IS at the target
					found = model.EqualsContainerReference(ref, &test.target)
					require.True(t, found, fmt.Sprintf("person[%s] is not at the target: %s", personID, destination.FormatReference(&test.target)))
				}
			}
		})
	}
}

func findPlayer(t *testing.T, id string) (*common.Reference, error) {

	courts, err := court.List()
	if err != nil {
		return nil, err
	}

	// Look for the player on one of the courts
	for _, courtID := range courts {

		ref := common.Reference{Type: "court", ID: courtID}
		c, err := court.Load(&ref)
		require.Nil(t, err, "err should be nothing")

		for _, personID := range c.Container.Players {

			if id == personID {
				ref := common.Reference{Type: "court", ID: courtID}
				return &ref, nil
			}
		}
	}

	// Look for the player on the queue
	ref := common.Reference{Type: "queue", ID: ""}
	q, err := queue.Load(&ref)
	require.Nil(t, err, "err should be nothing")

	for _, personID := range q.Container.Players {

		if id == personID {
			ref := common.Reference{Type: "queue", ID: ""}
			return &ref, nil
		}
	}

	return nil, nil
}
