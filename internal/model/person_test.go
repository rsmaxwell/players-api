package model

import (
	"database/sql"
	"testing"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/jackc/pgx/stdlib"
)

func TestPeopleBasic(t *testing.T) {
	teardown, db, _ := Setup(t)
	defer teardown(t)

	r := Registration{
		FirstName: "James2", LastName: "Bond2", DisplayName: "038", UserName: "018", Email: "018@mi6.gov.uk", Phone: "+44 1234 222222", Password: "TopSecret",
	}

	p, err := r.ToPerson()
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	err = p.SavePerson(db)
	if err != nil {
		t.Log("Could not create new person")
		t.Logf("%T   %s", err, err.Error())
		t.FailNow()
	}

	p.CheckPerson(t, db, r.FirstName, r.LastName, r.DisplayName, r.UserName, r.Email, r.Phone, r.Password, StatusSuspended)

	FirstName2 := "Smersh11"
	LastName2 := "Bomb11"
	DisplayName2 := "00711"
	UserName2 := "00711"
	Email2 := "00811@mi6.gov.uk"
	Phone2 := "+44 1234 222222"
	Password2 := "qwerty"

	hash, err := bcrypt.GenerateFromPassword([]byte(Password2), bcrypt.DefaultCost)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	p.FirstName = FirstName2
	p.LastName = LastName2
	p.DisplayName = DisplayName2
	p.UserName = UserName2
	p.Email = Email2
	p.Phone = Phone2
	p.Hash = hash
	p.Status = StatusNormal

	err = p.UpdatePerson(db)
	if err != nil {
		t.Log("Could not update person")
		t.FailNow()
	}

	p.CheckPerson(t, db, FirstName2, LastName2, DisplayName2, UserName2, Email2, Phone2, Password2, StatusNormal)

	var p2 Person
	p2.ID = p.ID
	err = p2.LoadPerson(db)
	if err != nil {
		t.Log("Could not load person")
		t.FailNow()
	}

	FirstName3 := "xxxxx"
	LastName3 := "yyyyy"
	DisplayName3 := "008"
	UserName3 := "009"
	Email3 := "010@mi6.gov.uk"
	Phone3 := "+44 1234 333333"
	Password3 := "topcat"

	hash, err = bcrypt.GenerateFromPassword([]byte(Password3), bcrypt.DefaultCost)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	p.FirstName = FirstName3
	p.LastName = LastName3
	p.DisplayName = DisplayName3
	p.UserName = UserName3
	p.Email = Email3
	p.Phone = Phone3
	p.Hash = hash
	p.Status = StatusNormal

	err = p.SavePerson(db)
	if err != nil {
		t.Log("Could not save person")
		t.FailNow()
	}

	p.CheckPerson(t, db, FirstName3, LastName3, DisplayName3, UserName3, Email3, Phone3, Password3, StatusNormal)

	err = p.DeletePersonBasic(db)
	if err != nil {
		t.Log("Could not delete person")
		t.FailNow()
	}
	err = p2.DeletePersonBasic(db)
	if err != nil {
		t.Log("Could not delete person")
		t.FailNow()
	}
}

func (p *Person) CheckPerson(t *testing.T, db *sql.DB, firstname string, lastname string, displayname string, username string, email string, phone string, password string, status string) {

	err := p.LoadPerson(db)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if p.FirstName != firstname {
		t.Logf("Unexpected FirstName. expected: '%s' actual: '%s'", firstname, p.FirstName)
		t.FailNow()
	}

	if p.LastName != lastname {
		t.Logf("Unexpected LastName. expected: '%s' actual: '%s'", lastname, p.LastName)
		t.FailNow()
	}

	if p.DisplayName != displayname {
		t.Logf("Unexpected displayName. expected: '%s' actual: '%s'", displayname, p.DisplayName)
		t.FailNow()
	}

	if p.UserName != username {
		t.Logf("Unexpected userName. expected: '%s' actual: '%s'", username, p.UserName)
		t.FailNow()
	}

	if p.Email != email {
		t.Logf("Unexpected Email. expected: '%s' actual: '%s'", email, p.Email)
		t.FailNow()
	}

	if p.Phone != phone {
		t.Logf("Unexpected Phone. expected: '%s' actual: '%s'", phone, p.Phone)
		t.FailNow()
	}

	if p.Status != status {
		t.Logf("Unexpected Status. expected: '%s' actual: '%s'", status, p.Status)
		t.FailNow()
	}

	err = bcrypt.CompareHashAndPassword(p.Hash, []byte(password))
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
}
