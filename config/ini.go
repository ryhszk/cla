package config

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"

	util "github.com/ryhszk/cla/utils"
)

type configList struct {
	// [main]
	LimitLine      int
	FocusedColor   string
	UnfocusedColor string
	// [keybind]
	ExecKey string
	SaveKey string
	DelKey  string
	AddKey  string
	QuitKey string
	// [data]
	DataFile string
}

// For reading config.ini
var Config configList

const defaultSetting = `
[main]
limitLine      = 101
focusedColor   = 82
unfocusedColor = 245

[keybind]
execKey = enter
saveKey = ctrl+s
delKey  = ctrl+d
addKey  = ctrl+a
quitKey = ctrl+c

[data]
dataFile = data.json
`

func init() {
	fpath := os.Getenv("HOME") + "/.cla/" + "config.ini"
	dir, _ := filepath.Split(fpath)
	util.AssumeDirExists(dir)
	util.AssumeFileExists(defaultSetting, fpath)

	cfg, err := ini.Load(fpath)
	if err != nil {
		util.ErrExit(err.Error())
	}

	Config = configList{
		// [main]
		LimitLine:      cfg.Section("main").Key("limitLine").MustInt(),
		FocusedColor:   cfg.Section("main").Key("focusedColor").String(),
		UnfocusedColor: cfg.Section("main").Key("unfocusedColor").String(),
		// [keybind]
		ExecKey: cfg.Section("keybind").Key("execKey").String(),
		SaveKey: cfg.Section("keybind").Key("saveKey").String(),
		DelKey:  cfg.Section("keybind").Key("delKey").String(),
		AddKey:  cfg.Section("keybind").Key("addKey").String(),
		QuitKey: cfg.Section("keybind").Key("quitKey").String(),
		// [data]
		DataFile: cfg.Section("data").Key("dataFile").String(),
	}
}
