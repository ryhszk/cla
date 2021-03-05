package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CmdData struct {
	ID   int    `json:"id"`
	Line string `json:"cmd"`
}

func FromJSON(fpath string) []CmdData {
	dir, _ := filepath.Split(fpath)

	AssumeDirExists(dir)

	fp, err := os.OpenFile(fpath, os.O_RDONLY|os.O_CREATE, 0664)
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
		data := CmdData{0, ""}
		s, _ := json.Marshal(data)
		jsonFmtStr := "[" + string(s) + "]"
		ToFile(jsonFmtStr, fpath)

		bytes = []byte(jsonFmtStr)
	}

	var datas []CmdData
	err = json.Unmarshal(bytes, &datas)
	if err != nil {
		ErrExit(err.Error())
	}

	return datas
}

func RemoveElement(datas []CmdData, rmLIdx int) []CmdData {
	var newDatas []CmdData
	var dataID = 0
	for i, d := range datas {
		if i == rmLIdx {
			continue
		}
		d.ID = dataID
		newDatas = append(newDatas, d)
		dataID++
	}
	return newDatas
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
