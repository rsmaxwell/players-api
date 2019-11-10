package sync

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	pkg = debug.NewPackage("sync")

	functionHandleDir  = debug.NewFunction(pkg, "HandleDir")
	functionHandleFile = debug.NewFunction(pkg, "handleFile")
	functionHashfile   = debug.NewFunction(pkg, "hashfile")
	functionCopyfile   = debug.NewFunction(pkg, "copyfile")
)

// HandleDir synchronises a directory with a reference directory
func HandleDir(reference, copy string) error {
	f := functionHandleDir
	f.DebugVerbose("reference: [%s], copy:[%s]", reference, copy)

	// Check the reference is a directory
	fi, err := os.Stat(reference)
	if err != nil {
		f.Dump("could not stat file [%s]\n%v", reference, err)
		return err
	}

	mode := fi.Mode()
	if !mode.IsDir() {
		err := fmt.Errorf("expected file [%s] to be a directory", reference)
		f.Dump("%v", err)
		return err
	}

	// Make sure the 'copy' is also a directory
	fi, err = os.Stat(copy)
	if err != nil {
		if os.IsNotExist(err) {

			f.DebugVerbose("calling os.MkdirAll(\"%s\", 0755)", copy)
			err = os.MkdirAll(copy, 0755)
			if err != nil {
				f.Dump("error creating directory [%s]\n%v", copy, err)
				return err
			}

			f.DebugVerbose("calling os.Chmod(\"%s\", 0755)", copy)
			err = os.Chmod(copy, 0755)
			if err != nil {
				f.Dump("error chmod directory [%s]\n%v", copy, err)
				return err
			}

			f.DebugVerbose("calling os.Stat(\"%s\")", copy)
			fi, err = os.Stat(copy)
			if err != nil {
				if os.IsNotExist(err) {
					f.Dump("could not find copy file [%s]\n%v", copy, err)
					return err
				}
				f.Dump("unexpected error on file [%s]\n%v", copy, err)
				return err
			}

			f.DebugVerbose("file perms of [%s]: %#o", copy, fi.Mode().Perm())

		} else {
			f.Dump("unexpected error on file [%s]\n%v", copy, err)
			return err
		}
	}

	if fi.Mode().IsDir() {
		f.DebugVerbose("the copy file:[%s] exists and is already a directory", copy)
	} else {
		f.DebugVerbose("the copy file:[%s] exists but is NOT a directory", copy)

		f.DebugVerbose("removing:[%s]", copy)
		err = os.RemoveAll(copy)
		if err != nil {
			f.Dump("could not remove file [%s]\n%v", copy, err)
			return err
		}

		f.DebugVerbose("calling os.MkdirAll(\"%s\", 0755)", copy)
		err = os.MkdirAll(copy, 0755)
		if err != nil {
			f.Dump("could not make directory [%s]\n%v", copy, err)
			return err
		}

		f.DebugVerbose("calling os.Chmod(\"%s\", 0755)", copy)
		err = os.Chmod(copy, 0755)
		if err != nil {
			f.Dump("error chmod directory [%s]\n%v", copy, err)
			return err
		}

		f.DebugVerbose("calling os.Stat(\"%s\")", copy)
		fi, err = os.Stat(copy)
		if err != nil {
			if os.IsNotExist(err) {
				f.Dump("could not find copy file [%s]\n%v", copy, err)
				return err
			}
			f.Dump("unexpected error on file [%s]\n%v", copy, err)
			return err
		}

		f.DebugVerbose("file perms of [%s]: %#o", copy, fi.Mode().Perm())
	}

	// Make sure all the files in the 'copy' also exists in the 'reference'
	files, err := ioutil.ReadDir(copy)
	if err != nil {
		f.Dump("could not read directory [%s]\n%v", copy, err)
		return err
	}

	for _, file := range files {

		reference2 := filepath.Join(reference, file.Name())
		copy2 := filepath.Join(copy, file.Name())

		_, err := os.Stat(reference2)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.RemoveAll(copy2)
				if err != nil {
					f.Dump("could not remove file [%s]\n %v", copy2, err)
					return err
				}
			} else {
				f.Dump("could not stat file [%s]\n%v", reference2, err)
				return err
			}
		}
	}

	// Synchronise all the files in the 'reference' with the 'copy'
	files, err = ioutil.ReadDir(reference)
	if err != nil {
		f.Dump("could not read directory [%s]\n%v", reference, err)
		return err
	}

	for _, file := range files {

		reference2 := filepath.Join(reference, file.Name())
		copy2 := filepath.Join(copy, file.Name())

		fi, err := os.Stat(reference2)
		if err != nil {
			f.Dump("could not stat directory [%s]\n%v", reference2, err)
			return err
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			err = HandleDir(reference2, copy2)
			if err != nil {
				f.Dump("could not copy [%s] to [%s]\n%v", reference2, copy2, err)
				return err
			}
		case mode.IsRegular():
			err = handleFile(reference2, copy2)
			if err != nil {
				f.Dump("could not get file info for [%s] to [%s]\n%v", reference2, copy2, err)
				return err
			}
		}

	}
	return nil
}

func handleFile(reference, copy string) error {
	f := functionHandleFile
	f.DebugVerbose("reference: [%s], copy:[%s]", reference, copy)

	// If the 'copy' does not exist, then copy the reference file
	_, err := os.Stat(copy)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = copyfile(reference, copy)
			f.Dump("could not stat file [%s]\n%v", copy, err)
			return err
		}
		f.Dump("unexpected error stating file [%s]\n%v", copy, err)
		return err
	}

	// If the hashes of the files do not match, then copy the reference file
	hashref, err := hashfile(reference)
	if err != nil {
		f.Dump("could not hash reference file [%s]\n%v", reference, err)
		return err
	}

	hashcopy, err := hashfile(copy)
	if err != nil {
		f.Dump("could not hash copy file [%s]\n%v", copy, err)
		return err
	}

	if bytes.Compare(hashref, hashcopy) != 0 {
		_, err = copyfile(reference, copy)
		if err != nil {
			f.Dump("could not compare file [%s] with [%s]\n%v", reference, copy, err)
			return err
		}
	}

	return nil
}

func hashfile(filename string) ([]byte, error) {
	f := functionHashfile
	f.DebugVerbose("filename: [%s]", filename)

	file, err := os.Open(filename)
	if err != nil {
		f.Dump("could not open file [%s]\n%v", filename, err)
		return nil, err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		f.Dump("could not copy file [%s]\n%v", filename, err)
		return nil, err
	}

	return h.Sum(nil), nil
}

func copyfile(reference, copy string) (int64, error) {
	f := functionCopyfile
	f.DebugVerbose("reference: [%s], copy: [%s]", reference, copy)

	sourceFileStat, err := os.Stat(reference)
	if err != nil {
		f.Dump("could not stat file [%s]\n%v", reference, err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", reference)
	}

	source, err := os.Open(reference)
	if err != nil {
		f.Dump("could not open file [%s]\n%v", reference, err)
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(copy)
	if err != nil {
		f.Dump("could not create file [%s]\n%v", copy, err)
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if err != nil {
		f.Dump("could not copy file [%v] to [%v]\n%v", destination, source, err)
		return 0, err
	}

	return nBytes, nil
}
