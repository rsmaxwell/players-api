package httphandler

import (
	"database/sql"
	"testing"

	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/stdlib"
)

func TestGetLoginToken(t *testing.T) {

	teardown, db, _ := model.Setup(t)
	defer teardown(t)

	cookie, accessToken := GetSigninToken(t, db, model.GoodEmail, model.GoodPassword)

	t.Logf("cookie: %s", cookie)
	t.Logf("token: %s", accessToken)
}

func GetFirstCourt(t *testing.T, db *sql.DB) *model.Court {

	listOfCourts, err := model.ListCourts(db)
	require.Nil(t, err, "err should be nothing")

	numberOfCourts := len(listOfCourts)
	require.True(t, numberOfCourts > 0, "There are no courts")

	var c = listOfCourts[0]
	return &c
}
