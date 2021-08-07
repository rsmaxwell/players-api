package model

import (
	"testing"

	_ "github.com/jackc/pgx/stdlib"
)

func TestPeople(t *testing.T) {

	teardown, db, _ := Setup(t)
	defer teardown(t)

	err := DeleteAllRecords(db)
	if err != nil {
		t.Log("Could not setup the model")
		t.FailNow()
	}

	err = Populate(db)
	if err != nil {
		t.Log("Could not populate")
		t.FailNow()
	}

	listOfCourts, err := ListCourts(db)
	if err != nil {
		t.Log("Could not list the courts")
		t.FailNow()
	}
	if len(listOfCourts) == 0 {
		t.Log("Could not find any courts")
		t.FailNow()
	}
	var c Court
	// c.ID = listOfCourts[0]
	// err = c.LoadCourt(db)
	// if err != nil {
	// 	t.Log("Could not load court")
	// 	t.FailNow()
	// }

	listOfWaiters, err := ListWaiters(db)
	if err != nil {
		t.Log("Could not get the first waiter")
		t.FailNow()
	}
	if len(listOfWaiters) == 0 {
		t.Log("Could not find any waiters")
		t.FailNow()
	}
	var p FullPerson
	p.ID = listOfWaiters[0].Person
	err = p.LoadPerson(db)
	if err != nil {
		t.Log("Could not get the first waiter")
		t.FailNow()
	}

	err = RemoveWaiter(db, p.ID)
	if err != nil {
		t.Log("Could not remove waiter")
		t.FailNow()
	}

	err = AddWaiter(db, p.ID)
	if err != nil {
		t.Log("Could not make a person into a player")
		t.FailNow()
	}

	err = MakePersonInactive(db, p.ID)
	if err != nil {
		t.Log("Could not make a person inactive")
		t.FailNow()
	}

	p.FirstName = "smersh"
	p.LastName = "Bomb"

	err = p.UpdatePerson(db)
	if err != nil {
		t.Log("Could not update person")
		t.FailNow()
	}

	var p2 FullPerson
	p2.ID = p.ID
	err = p2.LoadPerson(db)
	if err != nil {
		t.Log("Could not load person")
		t.FailNow()
	}

	p2.FirstName = "xxxxx"
	p2.Email = "fabdelkader.browx@balaways.com"
	p2.Phone = "+44 012 098765"
	err = p2.SavePerson(db)
	if err != nil {
		message := "Could not save person"
		t.Log(message)
		t.Log(err)
		t.FailNow()
	}

	c.Name = "AAAAA"

	err = c.UpdateCourt(db)
	if err != nil {
		t.Log("Could not update court")
		t.FailNow()
	}

	err = c.DeleteCourt(db)
	if err != nil {
		t.Log("Could not delete court")
		t.FailNow()
	}

	err = p.DeletePerson(db)
	if err != nil {
		t.Log("Could not delete person")
		t.FailNow()
	}
	err = p2.DeletePerson(db)
	if err != nil {
		t.Log("Could not delete person")
		t.FailNow()
	}
}
