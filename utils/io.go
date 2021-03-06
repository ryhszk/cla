package utils

import (
	"io/ioutil"
	"os"
)

// ToFile write the bytes to fPath.
// No return.
func ToFile(bytes, fPath string) {
	err := ioutil.WriteFile(fPath, []byte(bytes), 0664)
	if err != nil {
		ErrExit(err.Error())
	}
}

// AssumeDirExists create the dPath if dPath not exists.
// So, the fPath is treated as if it exists.
// No return.
func AssumeDirExists(dPath string) {
	if !Exists(dPath) {
		return
	}

	if err := os.Mkdir(dPath, 0774); err != nil {
		ErrExit(err.Error())
	}
}

// AssumeFileExists create the fPath with cont if fPath not exists.
// So, the fPath is treated as if it exists.
// No return.
func AssumeFileExists(cont, fPath string) {
	if !Exists(fPath) {
		return
	}
	ToFile(cont, fPath)
}

// Exists returns true if the path (file or directory) exists,
// false otherwise.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}
