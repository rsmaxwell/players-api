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
	"runtime"
)

// this logs the function name as well.
func handleError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}

// Dir synchronises a directory with a reference directory
func Dir(reference, copy string) error {

	// Check the reference is a directory
	fi, err := os.Stat(reference)
	if err != nil {
		handleError(err)
		return err
	}

	mode := fi.Mode()
	if !mode.IsDir() {
		handleError(err)
		return err
	}

	// Make sure the 'copy' is also a directory
	fi, err = os.Stat(copy)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(copy, 755)
			if err != nil {
				handleError(err)
				return err
			}

			fi, err = os.Stat(copy)
			if err != nil {
				if os.IsNotExist(err) {
					handleError(err)
					return err
				}
				handleError(err)
				return err
			}
		} else {
			handleError(err)
			return err
		}
	}

	if !fi.Mode().IsDir() {
		err = os.RemoveAll(copy)
		if err != nil {
			handleError(err)
			return err
		}

		err = os.MkdirAll(copy, 755)
		if err != nil {
			handleError(err)
			return err
		}
	}

	// Make sure all the files in the 'copy' also exists in the 'reference'
	files, err := ioutil.ReadDir(copy)
	if err != nil {
		handleError(err)
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
					handleError(err)
					return err
				}
			} else {
				handleError(err)
				return err
			}
		}
	}

	// Synchronise all the files in the 'reference' with the 'copy'
	files, err = ioutil.ReadDir(reference)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		reference2 := filepath.Join(reference, f.Name())
		copy2 := filepath.Join(copy, f.Name())

		fi, err := os.Stat(reference2)
		if err != nil {
			handleError(err)
			return err
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			err = Dir(reference2, copy2)
			if err != nil {
				handleError(err)
				return err
			}
		case mode.IsRegular():
			err = file(reference2, copy2)
			if err != nil {
				handleError(err)
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
			handleError(err)
			return err
		}
		handleError(err)
		return err
	}

	// If the hashes of the files do not match, then copy the reference file
	hashref, err := hashfile(reference)
	if err != nil {
		handleError(err)
		return err
	}

	hashcopy, err := hashfile(copy)
	if err != nil {
		handleError(err)
		return err
	}

	if bytes.Compare(hashref, hashcopy) != 0 {
		_, err = copyfile(reference, copy)
		if err != nil {
			handleError(err)
			return err
		}
	}

	return nil
}

func hashfile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		handleError(err)
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		handleError(err)
		return nil, err
	}

	return h.Sum(nil), nil
}

func copyfile(reference, copy string) (int64, error) {
	sourceFileStat, err := os.Stat(reference)
	if err != nil {
		handleError(err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", reference)
	}

	source, err := os.Open(reference)
	if err != nil {
		handleError(err)
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(copy)
	if err != nil {
		handleError(err)
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if err != nil {
		handleError(err)
		return 0, err
	}

	return nBytes, nil
}
