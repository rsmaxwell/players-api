package queue

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
)

func TestNewInfoJunkQueue(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	queuefile, err := makeQueueFilename()
	require.Nil(t, err, "err should be nothing")

	err = ioutil.WriteFile(queuefile, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	ref := &common.Reference{Type: "queue", ID: ""}
	_, err = Load(ref)
	if err != nil {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			if cerr.Code() == http.StatusInternalServerError {
				// ok
			} else {
				require.Fail(t, fmt.Sprintf("Unexpected error: %d", cerr.Code()))
			}
		} else {
			require.Fail(t, fmt.Sprintf("Unexpected error: %s", err.Error()))
		}
	} else {
		require.Fail(t, "Unexpected success")
	}
}

// SetupTestcase function
func setupTestcase(t *testing.T) func(t *testing.T) {

	queuefile, err := makeQueueFilename()
	require.Nil(t, err, "err should be nothing")

	_, err = os.Stat(queuefile)
	if err == nil {
		err = os.Remove(queuefile)
		require.Nil(t, err, "err should be nothing")
	}

	return func(t *testing.T) {
	}
}
