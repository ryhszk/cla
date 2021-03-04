package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"

	util "github.com/ryhszk/cla/utils"
)

type configList struct {
	// [main]
	LimitLine          int
	FocusedTextColor   string
	UnfocusedTextColor string
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
limit_line           = 32
focused_text_color   = 82
unfocused_text_color = 245

[keybind]
exec_key = enter
save_key = ctrl+s
del_key  = ctrl+d
add_key  = ctrl+a
quit_key = ctrl+c

[data]
data_file = data.json
`

func init() {
	fpath := os.Getenv("HOME") + "/.cla/" + "config.ini"
	dir, _ := filepath.Split(fpath)
	if util.IsExists(dir) {
		if err := os.Mkdir(dir, 0774); err != nil {
			util.ErrorExit(err.Error())
		}
	}

	if util.IsExists(fpath) {
		util.WriteToFile(defaultSetting, fpath)
	}

	cfg, err := ini.Load(fpath)
	if err != nil {
		log.Printf("Failed to read config.ini: %v", err)
		os.Exit(1)
	}

	Config = configList{
		// [main]
		LimitLine:          cfg.Section("main").Key("limit_line").MustInt(),
		FocusedTextColor:   cfg.Section("main").Key("focused_text_color").String(),
		UnfocusedTextColor: cfg.Section("main").Key("unfocused_text_color").String(),
		// [keybind]
		ExecKey: cfg.Section("keybind").Key("exec_key").String(),
		SaveKey: cfg.Section("keybind").Key("save_key").String(),
		DelKey:  cfg.Section("keybind").Key("del_key").String(),
		AddKey:  cfg.Section("keybind").Key("add_key").String(),
		QuitKey: cfg.Section("keybind").Key("quit_key").String(),
		// [data]
		DataFile: cfg.Section("data").Key("data_file").String(),
	}
}
