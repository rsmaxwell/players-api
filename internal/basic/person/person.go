package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
)

// Person type
type Person struct {
	FirstName      string `json:"firstname" validate:"required,min=3,max=20"`
	LastName       string `json:"lastname" validate:"required,min=3,max=20"`
	Email          string `json:"email" validate:"required,email"`
	HashedPassword []byte `json:"hashedpassword" validate:"required,len=60"`
	Role           string `json:"role" validate:"required,oneof=admin normal suspended"`
	Player         bool   `json:"player"`
}

var (
	personBaseDir string
	personListDir string

	// RoleAdmin is allowed to do anything!
	RoleAdmin string

	// RoleNormal can change themselves
	RoleNormal string

	// RoleSuspended can do nothing. Only the 'admin' can change a suspended person
	RoleSuspended string

	// AllRoles is the 'all' filter which returns every person
	AllRoles []string

	validate *validator.Validate
	pkg      *debug.Package
)

func init() {
	pkg = debug.NewPackage("person")
}

func init() {
	personBaseDir = common.RootDir + "/people"
	personListDir = personBaseDir + "/list"

	RoleAdmin = "admin"
	RoleNormal = "normal"
	RoleSuspended = "suspended"

	AllRoles = []string{RoleAdmin, RoleNormal, RoleSuspended}

	validate = validator.New()
}

// makeFilename function
func makeFilename(id string) (string, error) {

	err := common.CheckCharactersInID(id)
	if err != nil {
		return "", err
	}

	err = createFileStructure()
	if err != nil {
		return "", err
	}

	filename := personListDir + "/" + id + ".json"
	return filename, nil
}

// createFileStructure  creates the people directory
func createFileStructure() error {

	_, err := os.Stat(personListDir)
	if err != nil {
		err := os.MkdirAll(personListDir, 0755)
		if err != nil {
			return codeerror.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// CheckPassword - Basic check on the user calling the service
func (p *Person) CheckPassword(password string) bool {

	err := bcrypt.CompareHashAndPassword(p.HashedPassword, []byte(password))
	if err != nil {
		return false
	}

	return true
}

// New initialises a Person
func New(firstname string, lastname string, email string, hashedPassword []byte, player bool) *Person {
	person := new(Person)
	person.FirstName = firstname
	person.LastName = lastname
	person.Email = email
	person.HashedPassword = hashedPassword
	person.Player = player
	person.Role = RoleSuspended
	return person
}

// Update method
func Update(id string, fields map[string]interface{}) error {
	f := debug.NewFunction(pkg, "Update")
	f.DebugVerbose("fields: %v:", fields)

	p, err := Load(id)
	if err != nil {
		return err
	}

	err = p.updateFields(fields)
	if err != nil {
		return err
	}

	err = p.Save(id)
	if err != nil {
		return err
	}

	return nil
}

// updateFields method
func (p *Person) updateFields(fields map[string]interface{}) error {
	f := debug.NewFunction(pkg, "updateFields")
	f.DebugVerbose("fields: %v:", fields)

	if v, ok := fields["FirstName"]; ok {
		value, ok := v.(string)
		if !ok {
			return codeerror.NewBadRequest("The type of 'Person.FirstName' should be a string")
		}
		p.FirstName = value
	}

	if v, ok := fields["LastName"]; ok {
		value, ok := v.(string)
		if !ok {
			return codeerror.NewBadRequest("The type of 'Person.LastName' should be a string")
		}
		p.LastName = value
	}

	if v, ok := fields["Email"]; ok {
		value, ok := v.(string)
		if !ok {
			return codeerror.NewBadRequest("The type of 'Person.Email' should be a string")
		}
		p.Email = value
	}

	if v, ok := fields["Password"]; ok {
		value, ok := v.(string)
		if !ok {
			return codeerror.NewBadRequest("The type of 'Person.HashedPassword' should be a string")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		p.HashedPassword = hashedPassword
	}

	return nil
}

// UpdatePlayer method
func UpdatePlayer(id string, value bool) error {
	f := debug.NewFunction(pkg, "UpdatePlayer")
	f.DebugVerbose("id: %s, player: %t", id, value)

	p, err := Load(id)
	if err != nil {
		return err
	}

	err = p.updateFieldsPlayer(value)
	if err != nil {
		return err
	}

	err = p.Save(id)
	if err != nil {
		return err
	}

	return nil
}

// updateFieldsPlayer method
func (p *Person) updateFieldsPlayer(value bool) error {
	p.Player = value
	return nil
}

// UpdateRole method
func UpdateRole(id string, value string) error {
	f := debug.NewFunction(pkg, "UpdateRole")
	f.DebugVerbose("id: %s, role: %s", id, value)

	p, err := Load(id)
	if err != nil {
		return err
	}

	err = p.updateFieldsRole(value)
	if err != nil {
		return err
	}

	err = p.Save(id)
	if err != nil {
		return err
	}

	return nil
}

// updateFieldsRole method
func (p *Person) updateFieldsRole(value string) error {
	p.Role = value
	return nil
}

// List returns a list of the person IDs with one of the allowed role values
func List(filter []string) ([]string, error) {
	f := debug.NewFunction(pkg, "List")
	f.DebugVerbose("filter: %s", filter)

	err := createFileStructure()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	list := []string{}
	for _, filenameInfo := range files {
		filename := filenameInfo.Name()
		id := strings.TrimSuffix(filename, path.Ext(filename))

		f.DebugVerbose("loading: %s", id)
		p, err := Load(id)
		if err != nil {
			return nil, err
		}
		if !common.Contains(filter, p.Role) {
			f.DebugVerbose("skipping: %s", id)
			continue
		}

		f.DebugVerbose("adding: %s", id)
		list = append(list, id)
	}

	return list, nil
}

// IsPlayer returns 'true' if the person exists and is a player
func IsPlayer(id string) bool {

	person, err := Load(id)
	if err != nil {
		return false
	}
	if person == nil {
		return false
	}

	return person.Player
}

// CanLogin function
func (p *Person) CanLogin() bool {

	switch p.Role {
	case RoleAdmin:
		return true
	case RoleNormal:
		return true
	}

	return false
}

// CanUpdateCourt function
func (p *Person) CanUpdateCourt() bool {

	switch p.Role {
	case RoleAdmin:
		return true
	case RoleNormal:
		return true
	}

	return false
}

// CanUpdatePerson function
func CanUpdatePerson(sessionID, userID string) bool {

	p, err := Load(sessionID)
	if err != nil {
		return false
	}

	switch p.Role {
	case RoleAdmin:
		return true
	case RoleNormal:
		if sessionID == userID {
			return true
		}
	}

	return false
}

// CanUpdatePersonRole function
func CanUpdatePersonRole(sessionID, userID string) bool {

	p, err := Load(sessionID)
	if err != nil {
		return false
	}

	switch p.Role {
	case RoleAdmin:
		if sessionID != userID {
			return true
		}
	}

	return false
}

// CanUpdatePersonPlayer function
func CanUpdatePersonPlayer(sessionID, userID string) bool {

	p, err := Load(sessionID)
	if err != nil {
		return false
	}

	switch p.Role {
	case RoleAdmin:
		return true
	case RoleNormal:
		if sessionID == userID {
			return true
		}
	}

	return false
}

// CanGetMetrics function
func (p *Person) CanGetMetrics() bool {

	switch p.Role {
	case RoleAdmin:
		return true
	}

	return false
}

// ***************************************************************************
//
// ***************************************************************************

// Save writes a Person to disk
func (p *Person) Save(id string) error {
	f := debug.NewFunction(pkg, "Save")
	f.DebugVerbose("id: %s", id)

	// The first user must be made an 'admin' user
	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}
	if len(files) == 0 {
		p.Role = RoleAdmin
	}

	err = validate.Struct(p)
	if err != nil {
		return codeerror.NewBadRequest(err.Error())
	}

	personJSON, err := json.Marshal(p)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	filename, err := makeFilename(id)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, personJSON, 0644)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	return nil
}

// Load returns the Person with the given ID
func Load(id string) (*Person, error) {
	f := debug.NewFunction(pkg, "Load")
	f.DebugVerbose("id: %s", id)

	filename, err := makeFilename(id)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, codeerror.NewNotFound(err.Error())
		}
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	var p Person
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, codeerror.NewInternalServerError(err.Error())
	}
	return &p, nil
}

// Remove the person with the given ID
func Remove(id string) error {
	f := debug.NewFunction(pkg, "Remove")
	f.DebugVerbose("id: %s", id)

	filename, err := makeFilename(id)
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)

	if err == nil { // File exists
		err = os.Remove(filename)
		if err != nil {
			return codeerror.NewInternalServerError(err.Error())
		}
		return nil

	} else if os.IsNotExist(err) { // File does not exist
		return codeerror.NewNotFound(fmt.Sprintf("File Not Found: %s", filename))
	}

	return codeerror.NewInternalServerError(err.Error())
}

// Size returns the number of people
func Size() (int, error) {

	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		return 0, codeerror.NewInternalServerError(err.Error())
	}

	return len(files), nil
}

// Exists returns 'true' if the person exists
func Exists(id string) bool {

	filename, err := makeFilename(id)
	if err != nil {
		return false
	}

	_, err = os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}
