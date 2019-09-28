package destination

import (
	"fmt"

	"github.com/rsmaxwell/players-api/common"
)

// Destination is the Generic Destination interface
type Destination interface{}

var (
	baseDir string
)

func init() {
	baseDir = common.RootDir
}

// Reference structure
type Reference struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// FormatReference function
func FormatReference(ref *Reference) string {
	if ref.Type == "court" {
		return fmt.Sprintf("court[%s]", ref.ID)
	}

	return "queue"
}
