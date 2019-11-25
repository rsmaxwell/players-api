package queue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rsmaxwell/players-api/internal/basic/destination"
	"github.com/rsmaxwell/players-api/internal/basic/peoplecontainer"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
)

// Queue type
type Queue struct {
	destination.Destination
	Container peoplecontainer.PeopleContainer `json:"container"`
}

var (
	queueBaseDir string
)

func init() {
	queueBaseDir = common.RootDir
}

// makeQueueFilename function
func makeQueueFilename() (string, error) {
	filename := queueBaseDir + "/" + "queue" + ".json"

	return filename, nil
}

// createFileStructure creates the queue file
func createFileStructure() error {

	_, err := os.Stat(queueBaseDir)
	if err != nil {
		err := os.MkdirAll(queueBaseDir, 0755)
		if err != nil {
			return codeerror.NewInternalServerError(err.Error())
		}
	}

	filename, err := makeQueueFilename()
	if err != nil {
		return err
	}

	ref := common.Reference{Type: "queue", ID: ""}

	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) { // File does not exist
			err = New("Queue").Save(&ref)
		} else {
			return codeerror.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// New initialises a Queue object
func New(name string) *Queue {
	queue := new(Queue)
	queue.Container.Name = name
	queue.Container.Players = []string{}
	return queue
}

// Update method
func Update(ref *common.Reference, fields map[string]interface{}) error {

	q, err := Load(ref)
	if err != nil {
		return err
	}

	err = q.Update(fields)
	if err != nil {
		return err
	}

	err = q.Save(ref)
	if err != nil {
		return err
	}

	return nil
}

// Update method
func (q *Queue) Update(fields map[string]interface{}) error {

	if v, ok := fields["Container"]; ok {
		if container2, ok := v.(map[string]interface{}); ok {
			err := q.Container.Update(container2)
			if err != nil {
				return err
			}
		} else {
			return codeerror.NewInternalServerError(fmt.Sprintf("Unexpected Comtainer type: %v", v))
		}
	}

	return nil
}

// Save writes a Queue to disk
func (q *Queue) Save(ref *common.Reference) error {

	if ref.Type != "queue" {
		return codeerror.NewInternalServerError("Unexpected Reference type")
	}

	queueJSON, err := json.Marshal(q)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	filename, err := makeQueueFilename()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, queueJSON, 0644)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	return nil
}

// Exists returns 'true' if the queue exists
func Exists(id string) bool {

	filename, err := makeQueueFilename()
	if err != nil {
		return false
	}

	_, err = os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}

// Load returns the Queue
func Load(ref *common.Reference) (*Queue, error) {

	if ref.Type != "queue" {
		return nil, codeerror.NewInternalServerError("Unexpected Reference type")
	}

	filename, err := makeQueueFilename()
	if err != nil {
		return nil, err
	}

	err = createFileStructure()
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

	var q Queue
	err = json.Unmarshal(data, &q)
	if err != nil {
		return nil, codeerror.NewInternalServerError(err.Error())
	}
	return &q, nil
}

// GetContainer returns the Destination
func (q *Queue) GetContainer() *peoplecontainer.PeopleContainer {
	return &q.Container
}

// CheckPlayersLocation checks the players are at this destination
func (q *Queue) CheckPlayersLocation(players []string) error {
	pc := q.GetContainer()
	return destination.CheckPlayersInContainer(pc, players)
}

// CheckSpace checks there is space in the target for the moving players
func (q *Queue) CheckSpace(players []string) error {
	return nil
}

// RemovePlayers deletes players from the destination
func (q *Queue) RemovePlayers(players []string) error {
	pc := q.GetContainer()
	return destination.RemovePlayersFromContainer(pc, players)
}

// AddPlayers adds players to the destination
func (q *Queue) AddPlayers(players []string) error {
	pc := q.GetContainer()
	return destination.AddPlayersToContainer(pc, players)
}
