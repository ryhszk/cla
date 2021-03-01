package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tinp "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	ter "github.com/muesli/termenv"
	cfg "github.com/ryhszk/cla/config"
)

var (
	color              = ter.ColorProfile().Color
	focusedPrompt      = colorSetting("⇒ ", focusedTextColor)
	blurredPrompt      = "  "
	focusedTextColor   = cfg.Config.FocusedTextColor
	unfocusedTextColor = cfg.Config.UnfocusedTextColor
	dataFile           = os.Getenv("HOME") + "/.cla/" + cfg.Config.DataFile
	limitLineNumber    = cfg.Config.LimitLine
	execKey            = cfg.Config.ExecKey
	saveKey            = cfg.Config.SaveKey
	delKey             = cfg.Config.DelKey
	addKey             = cfg.Config.AddKey
	quitKey            = cfg.Config.QuitKey
)

func colorSetting(srcStr, colorCode string) string {
	return ter.String(srcStr).Foreground(color(colorCode)).String()
}

func getShellName() string {
	var shn string
	switch runtime.GOOS {
	case "windows":
		shn = "bash.exe"
	case "linux":
		shn = "sh"
	default:
		shn = "sh"
	}
	return shn
}

type Data struct {
	ID  int    `json:"id"`
	Cmd string `json:"cmd"`
}

func outErrorExit(err string) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fmt.Printf("call from '%s' function (line %d) \n", f.Name(), line)
	fmt.Printf("  err: %s\n", err)
	fmt.Print("  ")
	os.Exit(1)
}

func execCmd(cmd string) {
	c := exec.Command(getShellName(), "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}

func main() {
	result := make(chan string, 1)

	if err := tea.NewProgram(initialModel(result)).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	var cmd string
	if cmd = <-result; cmd != "" {
		execCmd(cmd)
	}
}

type model struct {
	index     int
	choice    chan string
	cmdInputs []tinp.Model
}

func isZeroSize(fp *os.File) bool {
	info, err := fp.Stat()
	if err != nil {
		outErrorExit(err.Error())
	}

	if info.Size() == 0 {
		return true
	}
	return false
}

func readFromJSON(fpath string) []Data {

	if f, err := os.Stat(fpath); os.IsNotExist(err) || f.IsDir() {
		dir, _ := filepath.Split(fpath)
		if err := os.Mkdir(dir, 0774); err != nil {
			outErrorExit(err.Error())
		}

	}

	fp, err := os.OpenFile(fpath, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		outErrorExit(err.Error())
	}
	defer fp.Close()

	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		outErrorExit(err.Error())
	}

	// When the file is created, the initial data is written in json format.
	// bytes variable the same.
	if isZeroSize(fp) {
		data := Data{0, ""}
		s, _ := json.Marshal(data)
		jsonFmtStr := "[" + string(s) + "]"
		writeToFile(jsonFmtStr, dataFile)

		bytes = []byte(jsonFmtStr)
	}

	var datas []Data
	err = json.Unmarshal(bytes, &datas)
	if err != nil {
		outErrorExit(err.Error())
	}

	return datas
}

func removeElementOfData(datas []Data, rmLIdx int) []Data {
	var newDatas []Data
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

func readFromFile() []string {
	f, err := os.OpenFile(dataFile, os.O_RDONLY, 0664)
	if err != nil {
		outErrorExit(err.Error())
	}
	defer f.Close()

	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > limitLineNumber {
			break
		}
	}
	if serr := scanner.Err(); serr != nil {
		outErrorExit(err.Error())
	}

	return lines
}

func writeToFile(bytes, fPath string) {
	err := ioutil.WriteFile(fPath, []byte(bytes), 0664)
	if err != nil {
		outErrorExit(err.Error())
	}
}

func writeToFileWithBlankLine(fPath string) {
	f, err := os.OpenFile(fPath, os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		outErrorExit(err.Error())
	}
	defer f.Close()

	fmt.Fprintln(f, "test")
}

func initialModel(ch chan string) model {
	tms := []tinp.Model{}
	for i, j := range readFromJSON(dataFile) {
		tm := tinp.NewModel()
		if i == 0 {
			tm.Focus()
			tm.TextColor = focusedTextColor
			tm.Prompt = focusedPrompt
		} else {
			tm.Prompt = blurredPrompt
		}
		tm.Placeholder = "Unregistered."
		tm.SetValue(j.Cmd)
		tm.CharLimit = 64
		tm.Width = 64
		tms = append(tms, tm)
	}

	// return model{0, inputs, blurredSubmitButton}
	return model{0, ch, tms}
}

// func LoggingSettings(logFile string) {
// 	// RDWRはreadとwrite。パーミッションで0666は読み書きができるユーザーその他。
// 	logfile, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	multiLogFile := io.MultiWriter(os.Stdout, logfile)
// 	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
// 	log.SetOutput(multiLogFile)
// }

func (m *model) addModel() {
	tm := tinp.NewModel()
	tm.Placeholder = "Unregistered"
	tm.Prompt = blurredPrompt
	tm.CharLimit = 64
	tm.Width = 64
	tm.SetValue("")
	m.cmdInputs = append(m.cmdInputs, tm)
}

func (m *model) removeModel(i int) {
	if i >= len(m.cmdInputs) {
		return
	}
	m.cmdInputs = append(m.cmdInputs[:i], m.cmdInputs[i+1:]...)
}

func (m model) Init() tea.Cmd {
	return tinp.Blink
}

var selectCmd string

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case quitKey:
			close(m.choice)
			return m, tea.Quit

		case addKey:
			if len(m.cmdInputs) >= limitLineNumber {
				return m, nil
			}
			newDatas := readFromJSON(dataFile)
			tailNumber := len(m.cmdInputs)
			emptyData := Data{tailNumber, ""}
			newDatas = append(newDatas, emptyData)
			newJsons, _ := json.Marshal(newDatas)
			writeToFile(string(newJsons), dataFile)
			m.addModel()

		// Cycle between inputs
		case "tab", "shift+tab", execKey, "up", "down", saveKey, delKey:

			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			// var cmdlines string

			if s == saveKey {
				var newDatas []Data
				var tmpData Data
				for i := 0; i < len(m.cmdInputs); i++ {
					tmpData.ID = i
					tmpData.Cmd = m.cmdInputs[i].Value()
					newDatas = append(newDatas, tmpData)
				}
				newJsons, _ := json.Marshal(newDatas)
				writeToFile(string(newJsons), dataFile)
			} else if s == delKey {
				// Load from file again to avoid unintended saving.
				oldDatas := readFromJSON(dataFile)
				newDatas := removeElementOfData(oldDatas, m.index)
				newJsons, _ := json.Marshal(newDatas)
				writeToFile(string(newJsons), dataFile)
				m.removeModel(m.index)

				// oldDatas := readFromJSON(dataFile)
				// newDatas := append(oldDatas[:m.index], oldDatas[m.index+1:]...)
				// newJsons, _ := json.Marshal(newDatas)
				// writeToFile(string(newJsons), dataFile)
				// m.removeModel(m.index)
			}

			if s == execKey {
				m.choice <- m.cmdInputs[m.index].Value()
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.index--
			} else {
				m.index++
			}

			if m.index > len(m.cmdInputs) {
				m.index = 0
			} else if m.index < 0 {
				m.index = len(m.cmdInputs)
			}

			for i := 0; i <= len(m.cmdInputs)-1; i++ {
				if i == m.index {
					// Set focused state
					m.cmdInputs[i].Focus()
					m.cmdInputs[i].Prompt = focusedPrompt
					m.cmdInputs[i].TextColor = focusedTextColor
					continue
				}
				// Remove focused state
				m.cmdInputs[i].Blur()
				m.cmdInputs[i].Prompt = blurredPrompt
				m.cmdInputs[i].TextColor = ""
			}

			return m, nil
		}
	}

	// Handle character input and blinks
	m, cmd = updateInputs(msg, m)
	return m, cmd
}

// Pass messages and models through to text input components. Only text inputs
// with Focus() set will respond, so it's safe to simply update all of them
// here without any further logic.
func updateInputs(msg tea.Msg, m model) (model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	for i := 0; i < len(m.cmdInputs); i++ {
		m.cmdInputs[i], cmd = m.cmdInputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := "\n"

	inputs := []string{}
	for i := 0; i < len(m.cmdInputs); i++ {
		inputs = append(inputs, m.cmdInputs[i].View())
	}

	for i := 0; i < len(inputs); i++ {
		s += fmt.Sprintf("%2d: %s\n", i, inputs[i])
	}

	s += "\n"
	s += colorSetting("______________________________________________\n", unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Execute selected line.  |\n", execKey), unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Save lines.             |\n", saveKey), unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Remove current line.    |\n", delKey), unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Add a line at end.      |\n", addKey), unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Exit.                   |\n", quitKey), unfocusedTextColor)
	s += colorSetting("|_____________________________________________|\n", unfocusedTextColor)
	s += "\n"

	return s
}
