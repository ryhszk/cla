package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// JsonData has the following data in json format.
// ID      ... Number for command identification.
// CmdLine ... User entered and saved string (command).
type JsonData struct {
	ID      int    `json:"id"`
	CmdLine string `json:"cmd"`
}

// FromJSON returns the data read from fPath as an array of type JsonData.
// If the contents of fPath are empty, write JsonData{0, ""} to a new file
// and use it as the read file (fPath).
func FromJSON(fPath string) []JsonData {
	dir, _ := filepath.Split(fPath)

	AssumeDirExists(dir)

	fp, err := os.OpenFile(fPath, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		ErrExit(err.Error())
	}
	defer fp.Close()

	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		ErrExit(err.Error())
	}

	// When the file is created, the initial data is written in json format.
	// bytes variable the same.
	if isZero(fp) {
		data := JsonData{0, ""}
		s, _ := json.Marshal(data)
		jsonFmt := "[" + string(s) + "]"
		ToFile(jsonFmt, fPath)

		bytes = []byte(jsonFmt)
	}

	var cds []JsonData
	if err = json.Unmarshal(bytes, &cds); err != nil {
		ErrExit(err.Error())
	}

	return cds
}

// RmElem returns an array of type JsonData with
// the elements of rmIndx removed from cds.
func RmElem(cds []JsonData, rmLIdx int) []JsonData {
	var newD []JsonData
	var id = 0
	for i, d := range cds {
		if i == rmLIdx {
			continue
		}
		d.ID = id
		newD = append(newD, d)
		id++
	}
	return newD
}

func isZero(fp *os.File) bool {
	info, err := fp.Stat()
	if err != nil {
		ErrExit(err.Error())
	}

	if info.Size() == 0 {
		return true
	}
	return false
}
