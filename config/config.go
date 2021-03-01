package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

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
quit_key = ctrl+z

[data]
data_file = data.json
`

var Config ConfigList

// onaji
func outErrorExit(err string) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fmt.Printf("call from '%s' function (line %d) \n", f.Name(), line)
	fmt.Printf("  err: %s\n", err)
	fmt.Print("  ")
	os.Exit(1)
}

// onaji
func isExists(name string) bool {
	_, err := os.Stat(name)
	return err != nil
}

// onaji
func writeToFile(bytes, fPath string) {
	err := ioutil.WriteFile(fPath, []byte(bytes), 0664)
	if err != nil {
		outErrorExit(err.Error())
	}
}

func init() {
	fpath := os.Getenv("HOME") + "/.cla/" + "config.ini"
	dir, _ := filepath.Split(fpath)
	// if f, err := os.Stat(dir); os.IsNotExist(err) || !f.IsDir() {
	if isExists(dir) {
		fmt.Println(dir)
		if err := os.Mkdir(dir, 0774); err != nil {
			outErrorExit(err.Error())
		}
	}

	if isExists(fpath) {
		fmt.Println("test")
		writeToFile(defaultSetting, fpath)
		// // if err := os.Rename("/home/fdr/git-mng/cla/config/config.ini", fpath); err != nil {
		// // 	outErrorExit(err.Error())
		// }
	}

	cfg, err := ini.Load(fpath)
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
