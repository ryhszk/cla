package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	// [main]
	limitLine          int
	focusedTextColor   string
	unfocusedTextColor string
	// [keybind]
	execKey string
	saveKey string
	delKey  string
	addKey  string
	quitKey string
	// [data]
	dataFile string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read config.ini: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		// [main]
		limitLine:          cfg.Section("main").Key("limit_line").MustInt(),
		focusedTextColor:   cfg.Section("main").Key("focused_text_color").String(),
		unfocusedTextColor: cfg.Section("main").Key("unfocused_text_color").String(),
		// [keybind]
		execKey: cfg.Section("keybind").Key("exec_key").String(),
		saveKey: cfg.Section("keybind").Key("save_key").String(),
		delKey:  cfg.Section("keybind").Key("del_key").String(),
		addKey:  cfg.Section("keybind").Key("add_key").String(),
		quitKey: cfg.Section("keybind").Key("quit_key").String(),
		// [data]
		dataFile: cfg.Section("data").Key("data_file").String(),
	}
}
