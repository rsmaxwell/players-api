package destination

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/common"
)

// Queue type
type Queue struct {
	Destination
	Container Container `json:"container"`
}

func init() {
	baseDir = common.RootDir
}

// ClearQueue the Queue
func ClearQueue() error {

	filename, err := makeQueueFilename()
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

	err = createQueueFiles()
	if err != nil {
		return err
	}

	return nil
}

// makeQueueFilename function
func makeQueueFilename() (string, error) {
	filename := baseDir + "/" + "queue" + ".json"
	return filename, nil
}

// createDirs creates the queue file
func createQueueFiles() error {

	_, err := os.Stat(baseDir)
	if err != nil {
		err := os.MkdirAll(baseDir, 0755)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	filename, err := makeQueueFilename()
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) { // File does not exist
			err = NewQueue("Queue").Save()
		} else {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// NewQueue initialises a Queue object
func NewQueue(name string) *Queue {
	queue := new(Queue)
	queue.Container.Name = name
	queue.Container.Players = []string{}
	return queue
}

// Update method
func (q Queue) Update(fields map[string]interface{}) error {

	if v, ok := fields["Container"]; ok {
		if container2, ok := v.(map[string]interface{}); ok {
			err := q.Container.Update(container2)
			if err != nil {
				return err
			}
		} else {
			return codeError.NewInternalServerError(fmt.Sprintf("Unexpected Comtainer type: %v", v))
		}
	}

	return nil
}

// Save writes a Queue to disk
func (q Queue) Save() error {

	queueJSON, err := json.Marshal(q)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	filename, err := makeQueueFilename()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, queueJSON, 0644)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	return nil
}

// QueueExists returns 'true' if the queue exists
func QueueExists(id string) bool {

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

// LoadQueue returns the Queue
func LoadQueue() (*Queue, error) {

	filename, err := makeQueueFilename()
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

// GetContainer returns the Destination
func (q Queue) GetContainer() *Container {

	fmt.Printf("GetContainer: &q=%p\n", &q)
	fmt.Printf("GetContainer: &q.Container=%p\n", &q.Container)

	return &q.Container
}
