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
)

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
func CheckPassword(id, password string) bool {

	p, err := Load(id)
	if err != nil {
		return false
	}
	if p == nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(p.HashedPassword, []byte(password))
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

	if v, ok := fields["HashedPassword"]; ok {
		value, ok := v.([]byte)
		if !ok {
			return codeerror.NewBadRequest("The type of 'Person.HashedPassword' should be a string")
		}
		p.HashedPassword = value
	}

	return nil
}

// UpdatePlayer method
func UpdatePlayer(id string, value bool) error {

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

// Add adds a person to the list
func (p *Person) Add(id string) error {

	// The first user must be made an 'admin' user
	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}
	if len(files) == 0 {
		p.Role = RoleAdmin
	}

	// Check the user does not already exist
	found := Exists(id)
	if found {
		return codeerror.NewInternalServerError(fmt.Sprintf("Person [%s] already exists", id))
	}

	// Save the updated court to disk
	err = p.Save(id)
	if err != nil {
		return err
	}

	return nil
}

// Save writes a Person to disk
func (p *Person) Save(id string) error {

	err := validate.Struct(p)
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

// List returns a list of the person IDs with one of the allowed role values
func List(filter []string) ([]string, error) {

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

		p, err := Load(id)
		if err != nil {
			return nil, err
		}
		if !common.Contains(filter, p.Role) {
			continue
		}

		list = append(list, id)
	}

	return list, nil
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

// Load returns the Person with the given ID
func Load(id string) (*Person, error) {

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

// CanLogin function
func CanLogin(id string) bool {

	p, err := Load(id)
	if err != nil {
		return false
	}

	switch p.Role {
	case RoleAdmin:
		return true
	case RoleNormal:
		return true
	}

	return false
}

// CanUpdateCourt function
func CanUpdateCourt(sessionID string) bool {

	p, err := Load(sessionID)
	if err != nil {
		return false
	}

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
func CanGetMetrics(sessionID string) bool {

	p, err := Load(sessionID)
	if err != nil {
		return false
	}

	switch p.Role {
	case RoleAdmin:
		return true
	}

	return false
}
