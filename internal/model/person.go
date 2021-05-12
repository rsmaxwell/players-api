package model

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgconn"
	"golang.org/x/crypto/bcrypt"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
)

// LimitedPerson type
type LimitedPerson struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

// Person type
type Person struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstName" validate:"required,min=3,max=20"`
	LastName    string `json:"lastName" validate:"required,min=3,max=20"`
	DisplayName string `json:"displayName" validate:"required,min=3,max=20"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required,min=3,max=20"`
	Hash        []byte `json:"hash"`
	Status      string `json:"status"`
}

// NullPerson type
type NullPerson struct {
	ID          int            `db:"id"`
	FirstName   sql.NullString `db:"firstname"`
	LastName    sql.NullString `db:"lastname"`
	DisplayName sql.NullString `db:"displayname"`
	Email       sql.NullString `db:"email"`
	Phone       sql.NullString `db:"phone"`
	Hash        sql.NullString `db:"hash"`
	Status      sql.NullString `db:"status"`
}

const (
	// PersonTable is the name of the person table
	PersonTable = "person"
)

var (
	functionUpdatePerson      = debug.NewFunction(pkg, "UpdatePerson")
	functionSavePerson        = debug.NewFunction(pkg, "SavePerson")
	functionFindPersonByEmail = debug.NewFunction(pkg, "FindPersonByEmail")
	functionListPeople        = debug.NewFunction(pkg, "ListPeople")
	functionPersonExists      = debug.NewFunction(pkg, "PersonExists")
	functionLoadPerson        = debug.NewFunction(pkg, "LoadPerson")
	functionDeletePersonBasic = debug.NewFunction(pkg, "DeletePersonBasic")
	functionAuthenticate      = debug.NewFunction(pkg, "Authenticate")
	functionCheckPassword     = debug.NewFunction(pkg, "CheckPassword")
)

const (
	// StatusAdmin constant
	StatusAdmin = "admin"

	// StatusNormal constant
	StatusNormal = "normal"

	// StatusSuspended constant
	StatusSuspended = "suspended"
)

var (
	// AllStates lists all the states
	AllStates []string
)

func init() {
	// AllRoles lists all the roles
	AllStates = []string{StatusAdmin, StatusNormal, StatusSuspended}
}

// NewPerson initialises a Person object
func NewPerson(firstname string, lastname string, displayName string, email string, phone string, hash []byte) *Person {
	p := new(Person)
	p.FirstName = firstname
	p.LastName = lastname
	p.DisplayName = displayName
	p.Email = email
	p.Phone = phone
	p.Hash = hash
	p.Status = StatusSuspended
	return p
}

// SavePerson writes a new Person to disk and returns the generated id
func (p *Person) SavePerson(db *sql.DB) error {
	f := functionSavePerson

	fields := "firstname, lastname, displayname, username, email, phone, hash, status"
	values := "$1, $2, $3, $4, $5, $6, $7, $8"
	sqlStatement := "INSERT INTO " + PersonTable + " (" + fields + ") VALUES (" + values + ") RETURNING id"

	err := db.QueryRow(sqlStatement, p.FirstName, p.LastName, p.DisplayName, p.Email, p.Phone, hex.EncodeToString(p.Hash), p.Status).Scan(&p.ID)
	if err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		if ok {
			if pgerr.Code == "23505" {
				return err
			}
		}

		message := "Could not insert into " + PersonTable
		f.Errorf(message)
		d := f.DumpSQLError(err, message, sqlStatement)
		p.Dump(d)
		return err
	}

	return nil
}

// UpdatePerson method
func (p *Person) UpdatePerson(db *sql.DB) error {
	f := functionUpdatePerson

	items := "firstname=$1, lastname=$2, displayname=$3, username=$4, email=$5, phone=$6, hash=$7, status=$8"
	sqlStatement := "UPDATE " + PersonTable + " SET " + items + " WHERE id=" + strconv.Itoa(p.ID)
	_, err := db.Exec(sqlStatement, p.FirstName, p.LastName, p.DisplayName, p.Email, p.Phone, hex.EncodeToString(p.Hash), p.Status)
	if err != nil {
		message := "Could not update person"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	return err
}

// LoadPerson returns the Person with the given ID
func (p *Person) LoadPerson(db *sql.DB) error {
	f := functionLoadPerson

	// Query the person
	fields := "firstname, lastname, displayname, username, email, phone, hash, status"
	sqlStatement := "SELECT " + fields + " FROM " + PersonTable + " WHERE ID=" + strconv.Itoa(p.ID)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not select all people"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++

		var np NullPerson
		err := rows.Scan(&np.FirstName, &np.LastName, &np.DisplayName, &np.Email, &np.Phone, &np.Hash, &np.Status)
		if err != nil {
			message := "Could not scan the person"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}

		if np.FirstName.Valid {
			p.FirstName = np.FirstName.String
		}

		if np.LastName.Valid {
			p.LastName = np.LastName.String
		}

		if np.DisplayName.Valid {
			p.DisplayName = np.DisplayName.String
		}

		if np.Email.Valid {
			p.Email = np.Email.String
		}

		if np.Phone.Valid {
			p.Phone = np.Phone.String
		}

		if np.Hash.Valid {
			p.Hash, err = hex.DecodeString(np.Hash.String)
			if err != nil {
				message := "Could not scan the Hash HexString"
				f.Errorf(message)
				f.DumpError(err, message)
				return err
			}
		}

		if np.Status.Valid {
			p.Status = np.Status.String
		}
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list the courts"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	if count == 0 {
		return codeerror.NewNotFound(fmt.Sprintf("People id %d not found", p.ID))
	} else if count > 1 {
		message := fmt.Sprintf("Found %d people with id %d", count, p.ID)
		err := codeerror.NewInternalServerError(message)
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	return nil
}

// DeletePersonBasic the person with the given ID
func (p *Person) DeletePersonBasic(db *sql.DB) error {
	f := functionDeletePersonBasic

	sqlStatement := "DELETE FROM " + PersonTable + " WHERE id=" + strconv.Itoa(p.ID) + " AND status != '" + StatusAdmin + "'"
	_, err := db.Exec(sqlStatement)
	if err != nil {
		message := "Could not delete person"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	return nil
}

// FindPersonByEmail function
func FindPersonByEmail(db *sql.DB, email string) (*Person, error) {
	f := functionFindPersonByEmail

	q := make(Query)
	q["username"] = Condition{Operation: "=", Value: email}

	arrayOfPeopleIDs, err := ListPeople(db, &q)
	if err != nil {
		return nil, err
	}

	if len(arrayOfPeopleIDs) <= 0 {
		err := codeerror.NewNotFound(fmt.Sprintf("Person not found: email:%s", email))
		return nil, err
	}

	if len(arrayOfPeopleIDs) > 1 {
		message := fmt.Sprintf("Too many matches. email:%s, count:%d", email, len(arrayOfPeopleIDs))
		err := codeerror.NewNotFound(message)
		d := f.DumpError(err, message)
		d.AddIntArray("peopleIDs.txt", arrayOfPeopleIDs)
		return nil, err
	}

	id := arrayOfPeopleIDs[0]
	p := Person{ID: id}
	err = p.LoadPerson(db)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// ListPeople returns a list of the people IDs
func ListPeople(db *sql.DB, query *Query) ([]int, error) {
	f := functionListPeople

	// Query the people
	allFields := []string{"id", "firstname", "lastname", "displayname", "username", "email", "phone", "status"}
	returnedFields := []string{"id"}
	sqlStatement, values, err := BuildQuery(PersonTable, allFields, returnedFields, query)
	rows, err := db.Query(sqlStatement, values...)
	if err != nil {
		message := "Could not select all from " + PersonTable
		f.Errorf(message)
		d := f.DumpSQLError(err, message, sqlStatement)
		d.AddArray("values.txt", values)
		return nil, err
	}
	defer rows.Close()

	var list []int
	for rows.Next() {

		var id int
		err := rows.Scan(&id)
		if err != nil {
			message := "Could not scan the person"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}
		list = append(list, id)
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list all from " + PersonTable
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}

	return list, nil
}

// PersonExists returns 'true' if the person exists
func (p *Person) PersonExists(db *sql.DB) (bool, error) {
	f := functionPersonExists

	// Query the person
	sqlStatement := "SELECT * FROM " + PersonTable + " WHERE id=$1"
	rows, err := db.Query(sqlStatement, p.ID)
	if err != nil {
		message := "Could not select people"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return false, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list the people"
		f.Errorf(message)
		f.DumpError(err, message)
		return false, err
	}

	if count == 0 {
		return false, nil
	} else if count > 1 {
		message := "Found " + string(count) + " people with id " + string(p.ID)
		f.Errorf(message)
		f.DumpError(err, message)
		return true, errors.New(message)
	}

	return true, nil
}

// Authenticate method
func (p *Person) Authenticate(db *sql.DB, password string) error {
	f := functionAuthenticate
	f.DebugVerbose("id: %d, password:%s", p.ID, "********")

	err := p.CheckPassword(password)
	if err != nil {
		f.DebugVerbose("password check failed for person [%d]", p.ID)
		return codeerror.NewUnauthorized("Not Authorized")
	}

	err = p.CanLogin()
	if err != nil {
		f.DebugVerbose("person [%d] not authorized to login", p.ID)
		return codeerror.NewForbidden("Forbidden")
	}

	return nil
}

// CheckPassword checks the validity of the password
func (p *Person) CheckPassword(password string) error {
	f := functionCheckPassword

	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(password))
	if err != nil {
		message := "The password was invalid for this user"
		f.Errorf(message)
		d := f.DumpError(err, message)
		d.AddString("hash.txt", hex.EncodeToString(p.Hash))
		return err
	}
	return nil
}

// CanLogin checks the user is allowed to login
func (p *Person) CanLogin() error {

	if p.Status == StatusAdmin {
		return nil
	}
	if p.Status == StatusNormal {
		return nil
	}

	return fmt.Errorf("Not Authorized")
}

// CanEditCourt checks the user is allowed update a court
func (p *Person) CanEditCourt() error {

	if p.Status == StatusAdmin {
		return nil
	}
	if p.Status == StatusNormal {
		return nil
	}

	return fmt.Errorf("Not Authorized")
}

// CanGetMetrics checks the user is allowed get the metrics
func (p *Person) CanGetMetrics() error {

	if p.Status == StatusAdmin {
		return nil
	}
	if p.Status == StatusNormal {
		return nil
	}

	return fmt.Errorf("Not Authorized")
}

// CanEditOtherPeople checks the user is allowed update a court
func (p *Person) CanEditOtherPeople() error {

	if p.Status == StatusAdmin {
		return nil
	}
	if p.Status == StatusNormal {
		return nil
	}

	return fmt.Errorf("Not Authorized")
}

// CanEditSelf checks the user is allowed update a court
func (p *Person) CanEditSelf() error {

	if p.Status == StatusAdmin {
		return nil
	}
	if p.Status == StatusNormal {
		return nil
	}

	return fmt.Errorf("Not Authorized")
}

// ToLimited converts a person to a Limited person
func (p *Person) ToLimited() *LimitedPerson {
	lp := &LimitedPerson{
		ID:          p.ID,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		DisplayName: p.DisplayName,
		Email:       p.Email,
		Phone:       p.Phone,
	}
	return lp
}

// ToLimitedPerson converts a NullPerson to a Limited person
func (np *NullPerson) ToLimitedPerson() *LimitedPerson {

	lp := LimitedPerson{ID: np.ID}

	if np.FirstName.Valid {
		lp.FirstName = np.FirstName.String
	}

	if np.LastName.Valid {
		lp.LastName = np.LastName.String
	}

	if np.DisplayName.Valid {
		lp.DisplayName = np.DisplayName.String
	}

	if np.Email.Valid {
		lp.Email = np.Email.String
	}

	if np.Phone.Valid {
		lp.Phone = np.Phone.String
	}

	return &lp
}

// Dump writes the person to a dump file
func (p *Person) Dump(d *debug.Dump) {

	bytearray, err := json.Marshal(p)
	if err != nil {
		return
	}

	title := fmt.Sprintf("person.%d.json", p.ID)
	d.AddByteArray(title, bytearray)
}
