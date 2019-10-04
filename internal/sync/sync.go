package sync

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Dir synchronises a directory with a reference directory
func Dir(reference, copy string) error {

	// Check the reference is a directory
	fi, err := os.Stat(reference)
	if err != nil {
		return err
	}

	mode := fi.Mode()
	if !mode.IsDir() {
		return err
	}

	// Make sure the 'copy' is also a directory
	fi, err = os.Stat(copy)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(copy, 777)
			if err != nil {
				return err
			}

			fi, err = os.Stat(copy)
			if err != nil {
				if os.IsNotExist(err) {
					return err
				}
				return err
			}
		} else {
			return err
		}
	}

	if !fi.Mode().IsDir() {
		err = os.RemoveAll(copy)
		if err != nil {
			return err
		}

		err = os.MkdirAll(copy, 777)
		if err != nil {
			return err
		}
	}

	// Make sure all the files in the 'copy' also exists in the 'reference'
	files, err := ioutil.ReadDir(copy)
	if err != nil {
		return err
	}

	for _, f := range files {

		reference2 := filepath.Join(reference, f.Name())
		copy2 := filepath.Join(copy, f.Name())

		_, err := os.Stat(reference2)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.RemoveAll(copy2)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// Synchronise all the file in the 'reference' with the 'copy'
	files, err = ioutil.ReadDir(reference)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		reference2 := filepath.Join(reference, f.Name())
		copy2 := filepath.Join(copy, f.Name())

		fi, err := os.Stat(reference2)
		if err != nil {
			return err
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			err = Dir(reference2, copy2)
			if err != nil {
				return err
			}
		case mode.IsRegular():
			err = file(reference2, copy2)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func file(reference, copy string) error {

	// If the 'copy' does not exist, then copy the reference file
	_, err := os.Stat(copy)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = copyfile(reference, copy)
			return err
		} else {
			return err
		}
	}

	// If the hashes of the files do not match, then copy the reference file
	hashref, err := hashfile(reference)
	if err != nil {
		return err
	}

	hashcopy, err := hashfile(copy)
	if err != nil {
		return err
	}

	if bytes.Compare(hashref, hashcopy) != 0 {
		_, err = copyfile(reference, copy)
		if err != nil {
			return err
		}
	}

	return nil
}

func hashfile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func copyfile(reference, copy string) (int64, error) {
	sourceFileStat, err := os.Stat(reference)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", reference)
	}

	source, err := os.Open(reference)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(copy)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
