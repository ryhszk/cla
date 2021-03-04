package utils

import (
	"io/ioutil"
	"os"
)

func IsExists(name string) bool {
	_, err := os.Stat(name)
	return err != nil
}

func IsZeroSize(fp *os.File) bool {
	info, err := fp.Stat()
	if err != nil {
		ErrorExit(err.Error())
	}

	if info.Size() == 0 {
		return true
	}
	return false
}

func WriteToFile(bytes, fPath string) {
	err := ioutil.WriteFile(fPath, []byte(bytes), 0664)
	if err != nil {
		ErrorExit(err.Error())
	}
}
