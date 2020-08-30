package logon

import (
	"database/sql"
	"errors"

	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// Logon type
type Logon struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,len=30"`
}

var (
	pkg = debug.NewPackage("logon")

	functionToPerson = debug.NewFunction(pkg, "ToPerson")
)

// New initialises a Logon object
func New(email string, password string) *Logon {
	l := new(Logon)
	l.Email = email
	l.Password = password
	return l
}

// ToPerson converts a Registration into a person
func (l *Logon) ToPerson(db *sql.DB) (*model.Person, error) {
	f := functionToPerson

	// Find the person with a matching email
	sqlStatement := "SELECT * FROM " + model.PersonTable + " WHERE email=" + basic.Quote(l.Email)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not select person"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return nil, err
	}
	defer rows.Close()

	var p model.Person
	count := 0
	for rows.Next() {
		count++

		var np model.NullPerson
		err := rows.Scan(&np.ID)
		if err != nil {
			message := "Could not scan the person"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}

		p.ID = np.ID
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list the people"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}

	if count == 0 {
		message := "Person email " + l.Email + " not found"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, errors.New(message)
	} else if count > 1 {
		message := "Found " + string(count) + " courts with email " + l.Email
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, errors.New(message)
	}

	err = p.LoadPerson(db)
	if err != nil {
		message := "Could not load the person"
		f.Errorf(message)
		return nil, err
	}

	// Compare the matching person's hash with the given password
	err = bcrypt.CompareHashAndPassword(p.Hash, []byte(l.Password))
	if err != nil {
		message := "The email/password is not valid"
		f.Errorf(message)
		return nil, err
	}
	return &p, nil
}
