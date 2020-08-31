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

	cookie := GetLoginToken(t, db, model.GoodUserName, model.GoodPassword)

	t.Logf("cookie: %s", cookie)
}

// FindPersonByUsername function
func FindPersonByUserName(t *testing.T, db *sql.DB, userName string) *model.Person {

	q := model.Query{}
	q.Conditions = make(map[string]model.Condition)
	q.Conditions["username"] = model.Condition{"=", userName}

	arrayOfPeopleIDs, err := model.ListPeople(db, &q)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(arrayOfPeopleIDs) < 1 {
		t.Errorf(err.Error())
	}

	id := arrayOfPeopleIDs[0]
	p := model.Person{ID: id}
	err = p.LoadPerson(db)
	if err != nil {
		t.Errorf(err.Error())
	}

	return &p
}

func GetFirstCourt(t *testing.T, db *sql.DB) *model.Court {

	listOfCourts, err := model.ListCourts(db)
	require.Nil(t, err, "err should be nothing")

	numberOfCourts := len(listOfCourts)
	require.True(t, numberOfCourts > 0, "There are no courts")

	var c model.Court
	c.ID = listOfCourts[0]
	err = c.LoadCourt(db)
	require.Nil(t, err, "err should be nothing")

	return &c
}
