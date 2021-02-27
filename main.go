package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	cfg "github.com/ryhszk/cla/config"

	tinp "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	ter "github.com/muesli/termenv"
)

const focusedTextColor = cfg.Config.focusedTextColor
const unfocusedTextColor = "245"

var (
	color         = ter.ColorProfile().Color
	focusedPrompt = colorSetting("⇒ ", focusedTextColor)
	blurredPrompt = "  "
	// focusedSubmitButton = "[ " + ter.String("Save").Foreground(color("82")).String() + " ]"
	// blurredSubmitButton = "[ " + ter.String("Save").Foreground(color("240")).String() + " ]"
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

func execCmd(cmd string) {
	c := exec.Command(getShellName(), "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}

// var mdata = initialModel(result)

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
	// submitButton string
}

const settingFile string = ".clarc"

const limitLineNumber = 32

func readFromFile() []string {
	f, err := os.OpenFile(settingFile, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not open: %v\n", settingFile, err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", settingFile, err)
	}

	return lines
}

func writeToFile(lines string) {
	err := ioutil.WriteFile(settingFile, []byte(lines), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not write: %v\n", settingFile, err)
		os.Exit(1)
	}
}

func writeToFileWithBlankLine() {
	// fmt.Println("execute this")
	f, err := os.OpenFile(settingFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not write: %v\n", settingFile, err)
	}
	defer f.Close()

	fmt.Fprintln(f, "test")
}

func initialModel(ch chan string) model {
	tms := []tinp.Model{}
	for i, cmd := range readFromFile() {
		tm := tinp.NewModel()
		if i == 0 {
			tm.Focus()
			tm.TextColor = focusedTextColor
			tm.Prompt = focusedPrompt
		} else {
			tm.Prompt = blurredPrompt
		}
		tm.Placeholder = "Unregistered."
		tm.SetValue(cmd)
		tm.CharLimit = 64
		tm.Width = 64
		tms = append(tms, tm)
	}

	// return model{0, inputs, blurredSubmitButton}
	return model{0, ch, tms}
}

func LoggingSettings(logFile string) {
	// RDWRはreadとwrite。パーミッションで0666は読み書きができるユーザーその他。
	logfile, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogFile)
}

func (m *model) addModel() {
	// LoggingSettings("test.log")
	// log.Printf("before %v", cap(m.cmdInputs))
	tm := tinp.NewModel()
	tm.Placeholder = "Unregistered"
	tm.Prompt = blurredPrompt
	tm.CharLimit = 64
	tm.Width = 64
	tm.SetValue("test")
	m.cmdInputs = append(m.cmdInputs, tm)
	// log.Printf("after %v", cap(m.cmdInputs))
}

func (m *model) removeModel(i int) {
	// m.cmdInputs

	if i >= len(m.cmdInputs) {
		return
	}
	m.cmdInputs = append(m.cmdInputs[:i], m.cmdInputs[i+1:]...)
}

// func unset(s []string, i int) []string {
// 	if i >= len(s) {
// 		return s
// 	}
// 	append(s[:i], s[i+1:]...)
// }

func (m model) Init() tea.Cmd {
	return tinp.Blink
}

var selectCmd string

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// inputs := m.cmdInputs
	// isUnderLimit := false
	// if len(inputs) > limitLineNumber {
	// 	isUnderLimit = true
	// }
	// LoggingSettings("test.log")
	// log.Printf("in the Update: %v %v", cap(m.cmdInputs), len(m.cmdInputs))

	// writeToFileWithBlankLine()
	// m.addModel()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+q", "ctrl+z":
			close(m.choice)
			return m, tea.Quit

		case "ctrl+a":
			if len(m.cmdInputs) < limitLineNumber {
				writeToFileWithBlankLine()
				m.addModel()
			}

		// 	// return m, nil

		// case "ctrl+d":

		// Cycle between inputs
		case "tab", "shift+tab", "enter", "up", "down", "ctrl+s", "ctrl+d":

			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			var cmdlines string
			// if s == "enter" && m.index == len(inputs) {
			// 	for i := 0; i < len(inputs); i++ {
			// 		cmdlines += inputs[i].Value() + "\n"
			// 	}
			// 	writeToFile(cmdlines)
			// 	//return m, tea.Quit
			// } else if s == "enter" || s == "ctrl+s" {
			// 	selectCmd = inputs[m.index].Value()
			// 	return m, tea.Quit
			// }
			if s == "ctrl+s" {
				for i := 0; i < len(m.cmdInputs); i++ {
					cmdlines += m.cmdInputs[i].Value() + "\n"
				}
				writeToFile(cmdlines)
				//return m, tea.Quit
			} else if s == "ctrl+d" {
				for i, cmd := range readFromFile() {
					if i == m.index {
						continue
					}
					cmdlines += cmd + "\n"
				}
				// log.Printf("in the Update: %v", cmdlines)
				writeToFile(cmdlines)
				// log.Printf("in the Update: inputs %v", len(m.cmdInputs))
				m.removeModel(m.index)
				// log.Printf("in the Update: inputs %v", len(m.cmdInputs))

				//return m, tea.Quit
			}

			if s == "enter" {
				m.choice <- m.cmdInputs[m.index].Value()
				return m, tea.Quit
			}

			// Cycle indexes
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
			// if s == "up" || s == "shift+tab" {
			// 	// if m.index < 0 {
			// 	// 	m.index = 0
			// 	// } else {
			// 	m.index--
			// 	// }
			// } else if s == "down" || s == "tab" {
			// 	m.index++
			// }

			// if m.index >= len(m.cmdInputs) {
			// 	m.index = len(m.cmdInputs)
			// } else if m.index < 0 {
			// 	m.index = 0
			// }

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

			// for i := 0; i < len(m.cmdInputs); i++ {
			// 	m.cmdInputs[i] = m.cmdInputs[i]
			// }
			// m.nameInput = inputs[0]
			// m.cmdInputs[0] = inputs[0]
			// m.cmdInputs[1] = inputs[1]

			// if m.index == len(inputs) {
			// 	m.submitButton = focusedSubmitButton
			// } else {
			// 	m.submitButton = blurredSubmitButton
			// }

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
	// s += "* Which command do you want to run ?\n\n"

	inputs := []string{}
	for i := 0; i < len(m.cmdInputs); i++ {
		inputs = append(inputs, m.cmdInputs[i].View())
	}

	for i := 0; i < len(inputs); i++ {
		s += fmt.Sprintf("%2d: %s\n", i, inputs[i])
		// if i < len(inputs)-1 {
		// 	s += "\n"
		// }
	}
	s += "\n"
	s += colorSetting("+---------------------------------------+\n", unfocusedTextColor)
	s += colorSetting("| enter      ... Execute selected line. |\n", unfocusedTextColor)
	s += colorSetting("| ctrl+[q|z] ... Exit.                  |\n", unfocusedTextColor)
	s += colorSetting("| ctrl+s     ... Save lines.            |\n", unfocusedTextColor)
	s += colorSetting("| ctrl+d     ... Remove current line.   |\n", unfocusedTextColor)
	s += colorSetting("| ctrl+a     ... Add a line at the end. |\n", unfocusedTextColor)
	s += colorSetting("+---------------------------------------+\n", unfocusedTextColor)
	s += "\n"
	// s += "\n\n  " + m.submitButton + "\n\n"
	// s += string(len(inputs))
	// s += "\n\n  " + string(m.index) + "\n\n"
	// mdata.cmdInputs = append(mdata.cmdInputs, tinp.NewModel())
	// mdata.cmdInputs[1].SetValue("aaaaa")
	// s += mdata.cmdInputs[1].Value()
	// s += "\n\n  " + m.Value + "\n\n"

	return s
}
