package queue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/common"
	"github.com/rsmaxwell/players-api/container"
)

// Queue type
type Queue struct {
	Container container.Container `json:"container"`
}

var (
	baseDir string
)

func init() {
	baseDir = common.RootDir
}

// Clear the Queue
func Clear() error {

	filename, err := makeFilename()
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)
	if err == nil {
		err = os.Remove(filename)
		if err != nil {
			return err
		}
	}

	err = createDirs()
	if err != nil {
		return err
	}

	return nil
}

// makeFilename function
func makeFilename() (string, error) {
	filename := baseDir + "/" + "queue" + ".json"
	return filename, nil
}

// createDirs creates the queue file
func createDirs() error {

	_, err := os.Stat(baseDir)
	if err != nil {
		err := os.MkdirAll(baseDir, 0755)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	filename, err := makeFilename()
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) { // File does not exist
			err = Save(New("Queue"))
		} else {
			return codeError.NewInternalServerError(err.Error())
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
func Update(id string, court2 map[string]interface{}) (*Queue, error) {

	queue, err := Load()
	if err != nil {
		return nil, codeError.NewNotFound(fmt.Sprintf("Queue not found"))
	}

	if v, ok := court2["Container"]; ok {
		if container2, ok := v.(map[string]interface{}); ok {
			err = container.Update(&queue.Container, container2)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, codeError.NewInternalServerError(fmt.Sprintf("Unexpected Comtainer type: %v", v))
		}
	}

	// Save the updated court to disk
	err = Save(queue)
	if err != nil {
		return nil, err
	}

	return queue, nil
}

// Save writes a Queue to disk
func Save(queue *Queue) error {

	queueJSON, err := json.Marshal(queue)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	filename, err := makeFilename()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, queueJSON, 0644)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	return nil
}

// Exists returns 'true' if the queue exists
func Exists(id string) bool {

	filename, err := makeFilename()
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
func Load() (*Queue, error) {

	filename, err := makeFilename()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, codeError.NewNotFound(err.Error())
		}
		return nil, codeError.NewInternalServerError(err.Error())
	}

	var q Queue
	err = json.Unmarshal(data, &q)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}
	return &q, nil
}

// CheckConsistency function
func CheckConsistency() error {
	return nil
}
