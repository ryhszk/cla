package utils

import (
	"io/ioutil"
	"os"
)

func ToFile(bytes, fPath string) {
	err := ioutil.WriteFile(fPath, []byte(bytes), 0664)
	if err != nil {
		ErrExit(err.Error())
	}
}

func AssumeDirExists(dirPath string) {
	if !Exists(dirPath) {
		return
	}

	if err := os.Mkdir(dirPath, 0774); err != nil {
		ErrExit(err.Error())
	}
}

func AssumeFileExists(contents, filePath string) {
	if !Exists(filePath) {
		return
	}
	ToFile(contents, filePath)
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return err != nil
}
