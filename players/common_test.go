package players

import (
	"io/ioutil"
)

func writefile(filepath string, contents string) error {

	data := []byte(contents)
	err := ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func writefileInDirectory(directory string, filename string, contents string) error {
	return writefile(directory+"/"+filename, contents)
}
