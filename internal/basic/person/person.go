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

const (
	// RoleAdmin constant
	RoleAdmin = "admin"

	// RoleNormal constant
	RoleNormal = "normal"

	// RoleSuspended constant
	RoleSuspended = "suspended"
)

var (
	personBaseDir string
	personListDir string

	validate = validator.New()

	pkg = debug.NewPackage("person")

	functionCreateFileStructure   = debug.NewFunction(pkg, "createFileStructure")
	functionCheckPassword         = debug.NewFunction(pkg, "CheckPassword")
	functionUpdate                = debug.NewFunction(pkg, "Update")
	functionUpdateFields          = debug.NewFunction(pkg, "updateFields")
	functionUpdatePlayer          = debug.NewFunction(pkg, "UpdatePlayer")
	functionUpdateRole            = debug.NewFunction(pkg, "UpdateRole")
	functionList                  = debug.NewFunction(pkg, "List")
	functionCanUpdatePerson       = debug.NewFunction(pkg, "CanUpdatePerson")
	functionCanUpdatePersonRole   = debug.NewFunction(pkg, "CanUpdatePersonRole")
	functionCanUpdatePersonPlayer = debug.NewFunction(pkg, "CanUpdatePersonPlayer")
	functionSave                  = debug.NewFunction(pkg, "Save")
	functionLoad                  = debug.NewFunction(pkg, "Load")
	functionRemove                = debug.NewFunction(pkg, "Remove")
	functionSize                  = debug.NewFunction(pkg, "Size")
	functionExists                = debug.NewFunction(pkg, "Exists")

	// AllRoles lists all the roles
	AllRoles []string
)

func init() {

	personBaseDir = common.RootDir + "/people"
	personListDir = personBaseDir + "/list"

	// AllRoles lists all the roles
	AllRoles = []string{RoleAdmin, RoleNormal, RoleSuspended}
}

// makeFilename function
func makeFilename(id string) (string, error) {

	err := common.CheckCharactersInID(id)
	if err != nil {
		return "", err
	}

	filename := personListDir + "/" + id + ".json"
	return filename, nil
}

// createFileStructure  creates the people directory
func createFileStructure() error {
	f := functionCreateFileStructure

	_, err := os.Stat(personListDir)
	if err != nil {
		err := os.MkdirAll(personListDir, 0755)
		if err != nil {
			message := fmt.Sprintf("could not make people list directory[%s]: %v", personListDir, err)
			f.Dump(message)
			return codeerror.NewInternalServerError(message)
		}
	}

	return nil
}

// CheckPassword - Basic check on the user calling the service
func (p *Person) CheckPassword(password string) bool {
	f := functionCheckPassword

	err := bcrypt.CompareHashAndPassword(p.HashedPassword, []byte(password))
	if err != nil {
		f.DebugVerbose("password check failed: %v", err)
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
	f := functionUpdate
	f.DebugVerbose("id: %s, fields: %v:", id, fields)

	p, err := Load(id)
	if err != nil {
		return err
	}

	err = p.updateFields(fields)
	if err != nil {
		message := fmt.Sprintf("could not update person fields. person[%s]: %v", id, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	err = p.Save(id)
	if err != nil {
		message := fmt.Sprintf("could not save person[%s]: %v", id, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	return nil
}

// updateFields method
func (p *Person) updateFields(fields map[string]interface{}) error {
	f := functionUpdateFields
	f.DebugVerbose("fields: %v:", fields)

	if v, ok := fields["FirstName"]; ok {
		value, ok := v.(string)
		if !ok {
			message := fmt.Sprintf("The type of 'Person.FirstName' should be a string")
			f.DebugVerbose(message)
			return codeerror.NewBadRequest(message)
		}
		p.FirstName = value
	}

	if v, ok := fields["LastName"]; ok {
		value, ok := v.(string)
		if !ok {
			message := fmt.Sprintf("The type of 'Person.LastName' should be a string")
			f.DebugVerbose(message)
			return codeerror.NewBadRequest(message)
		}
		p.LastName = value
	}

	if v, ok := fields["Email"]; ok {
		value, ok := v.(string)
		if !ok {
			message := fmt.Sprintf("The type of 'Person.Email' should be a string")
			f.DebugVerbose(message)
			return codeerror.NewBadRequest(message)
		}
		p.Email = value
	}

	if v, ok := fields["Password"]; ok {
		value, ok := v.(string)
		if !ok {
			message := fmt.Sprintf("The type of 'Person.HashedPassword' should be a string")
			f.DebugVerbose(message)
			return codeerror.NewBadRequest(message)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
		if err != nil {
			message := fmt.Sprintf("could not hash the password: %v", err)
			f.Dump(message)
			return codeerror.NewBadRequest(message)
		}

		p.HashedPassword = hashedPassword
	}

	return nil
}

// UpdatePlayer method
func UpdatePlayer(id string, value bool) error {
	f := functionUpdatePlayer
	f.DebugVerbose("id: %s, player: %t", id, value)

	p, err := Load(id)
	if err != nil {
		return err
	}

	err = p.updateFieldsPlayer(value)
	if err != nil {
		message := fmt.Sprintf("could not update player field for person[%s]: %v", id, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	err = p.Save(id)
	if err != nil {
		message := fmt.Sprintf("could not save person[%s]: %v", id, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
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
	f := functionUpdateRole
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
	f := functionList
	f.DebugVerbose("filter: %s", filter)

	err := createFileStructure()
	if err != nil {
		message := fmt.Sprintf("could not create person file structure: %v", err)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
	}

	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		message := fmt.Sprintf("could not read the personListDir directory [%s]: %v", personListDir, err)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
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
func (p *Person) IsPlayer() bool {
	return p.Player
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
func (p *Person) CanUpdatePerson(sessionID, userID string) bool {

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
func (p *Person) CanUpdatePersonRole(sessionID, userID string) bool {

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
func (p *Person) CanUpdatePersonPlayer() bool {

	switch p.Role {
	case RoleAdmin:
		return true
	case RoleNormal:
		return true
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
	f := functionSave
	f.DebugVerbose("id: %s", id)

	err := createFileStructure()
	if err != nil {
		message := fmt.Sprintf("could not create the person file structure: %v", err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	// The first user must be made an 'admin' user
	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		message := fmt.Sprintf("could not read the personListDir directory[%s]: %v", personListDir, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}
	if len(files) == 0 {
		p.Role = RoleAdmin
	}

	err = validate.Struct(p)
	if err != nil {
		message := fmt.Sprintf("validation failed for person[%s]: %v", id, err)
		f.DebugVerbose(message)
		return codeerror.NewBadRequest(message)
	}

	personJSON, err := json.Marshal(p)
	if err != nil {
		message := fmt.Sprintf("could not marshal person[%s]: %v", id, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	filename, err := makeFilename(id)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, personJSON, 0644)
	if err != nil {
		message := fmt.Sprintf("could not write file [%s] for person[%s]: %v", filename, id, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	return nil
}

// Load returns the Person with the given ID
func Load(id string) (*Person, error) {
	f := functionLoad
	f.DebugVerbose("id: %s", id)

	filename, err := makeFilename(id)
	if err != nil {
		message := fmt.Sprintf("could not make filename for person[%s]: %v", id, err)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			message := fmt.Sprintf("the file [%s] does not exist, for person[%s]", filename, id)
			f.Dump(message)
			return nil, codeerror.NewNotFound(message)
		}

		message := fmt.Sprintf("could not read the file [%s] for person[%s]: %v", filename, id, err)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
	}

	var p Person
	err = json.Unmarshal(data, &p)
	if err != nil {
		message := fmt.Sprintf("could not unmarshal the contents of file [%s] for person[%s]: %v", filename, id, err)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
	}
	return &p, nil
}

// Remove the person with the given ID
func Remove(id string) error {
	f := functionRemove
	f.DebugVerbose("id: %s", id)

	filename, err := makeFilename(id)
	if err != nil {
		message := fmt.Sprintf("could not make filename for person[%s]: %v", id, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	_, err = os.Stat(filename)

	if err == nil { // File exists
		err = os.Remove(filename)
		if err != nil {
			message := fmt.Sprintf("could not remove file [%s] for person[%s]: %v", filename, id, err)
			f.Dump(message)
			return codeerror.NewInternalServerError(message)
		}
		return nil

	} else if os.IsNotExist(err) { // File does not exist
		message := fmt.Sprintf("the file [%s] for person[%s] could not be found: %v", filename, id, err)
		f.Dump(message)
		return codeerror.NewNotFound(message)
	}

	message := fmt.Sprintf("could not stat the file [%s] for person[%s]: %v", filename, id, err)
	f.Dump(message)
	return codeerror.NewInternalServerError(message)
}

// Size returns the number of people
func Size() (int, error) {
	f := functionSize

	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		message := fmt.Sprintf("could not read the personListDir [%s]: %v", personListDir, err)
		f.Dump(message)
		return 0, codeerror.NewInternalServerError(message)
	}

	return len(files), nil
}

// Exists returns 'true' if the person exists
func Exists(id string) bool {
	f := functionExists

	filename, err := makeFilename(id)
	if err != nil {
		f.Dump("could not make the filename for person[%s]: %v", id, err)
		return false
	}

	_, err = os.Stat(filename)
	if err != nil {
		f.DebugVerbose("could not stat the file[%s] for person[%s]: %v", filename, id, err)
		return false
	}

	return true
}
