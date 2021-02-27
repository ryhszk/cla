package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
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

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read config.ini: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
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
