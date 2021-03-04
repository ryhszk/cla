package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	tinp "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	ter "github.com/muesli/termenv"

	cfg "github.com/ryhszk/cla/config"

	util "github.com/ryhszk/cla/utils"

	"golang.org/x/term"
)

var (
	color              = ter.ColorProfile().Color
	focusedPrompt      = colorSetting("> ", focusedTextColor)
	blurredPrompt      = "  "
	focusedTextColor   = cfg.Config.FocusedTextColor
	unfocusedTextColor = cfg.Config.UnfocusedTextColor
	dataFile           = os.Getenv("HOME") + "/.cla" + cfg.Config.DataFile
	limitLineNumber    = cfg.Config.LimitLine
	execKey            = cfg.Config.ExecKey
	saveKey            = cfg.Config.SaveKey
	delKey             = cfg.Config.DelKey
	addKey             = cfg.Config.AddKey
	quitKey            = cfg.Config.QuitKey
)

type model struct {
	index     int
	choice    chan string
	cmdInputs []tinp.Model
}

func main() {
	result := make(chan string, 1)

	if err := tea.NewProgram(initialModel(result)).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	var cmd string
	if cmd = <-result; cmd != "" {
		util.Execute(cmd)
	}
}

func colorSetting(srcStr, colorCode string) string {
	return ter.String(srcStr).Foreground(color(colorCode)).String()
}

func initialModel(ch chan string) model {
	tms := []tinp.Model{}
	for i, j := range util.ReadFromJSON(dataFile) {
		tm := tinp.NewModel()
		if i == 0 {
			tm.Focus()
			tm.TextColor = focusedTextColor
			tm.Prompt = focusedPrompt
		} else {
			tm.Prompt = blurredPrompt
		}
		tm.Placeholder = "Input any command."
		tm.SetValue(j.Line)
		tm.CharLimit = 99
		tm.Width = 99
		tms = append(tms, tm)
	}
	return model{0, ch, tms}
}

func (m *model) addModel() {
	tm := tinp.NewModel()
	tm.Placeholder = "Input any command."
	tm.Prompt = blurredPrompt
	tm.CharLimit = 99
	tm.Width = 99
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
			newDatas := util.ReadFromJSON(dataFile)
			tailNumber := len(m.cmdInputs)
			emptyData := util.CmdData{tailNumber, ""}
			newDatas = append(newDatas, emptyData)
			newJsons, _ := json.Marshal(newDatas)
			util.WriteToFile(string(newJsons), dataFile)
			m.addModel()

		// Cycle between inputs
		case "tab", "shift+tab", execKey, "up", "down", saveKey, delKey:

			s := msg.String()

			if s == saveKey {
				var newDatas []util.CmdData
				var tmpData util.CmdData
				for i := 0; i < len(m.cmdInputs); i++ {
					tmpData.ID = i
					tmpData.Line = m.cmdInputs[i].Value()
					newDatas = append(newDatas, tmpData)
				}
				newJsons, _ := json.Marshal(newDatas)
				util.WriteToFile(string(newJsons), dataFile)
			} else if s == delKey {
				// Load from file again to avoid unintended saving.
				oldDatas := util.ReadFromJSON(dataFile)
				newDatas := util.RemoveElementOfData(oldDatas, m.index)
				newJsons, _ := json.Marshal(newDatas)
				util.WriteToFile(string(newJsons), dataFile)
				m.removeModel(m.index)
				if m.index > len(m.cmdInputs)-1 {
					m.index = -2 // 要素が消える分
				} else {
					m.index--
				}
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

			if m.index > len(m.cmdInputs)-1 {
				m.index = 0
			}

			if m.index < 0 {
				m.index = len(m.cmdInputs) - 1
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
	_, h, _ := term.GetSize(syscall.Stdin)
	s := "\n"
	s += colorSetting("______________________________________________\n", unfocusedTextColor)
	inputs := []string{}
	for i := 0; i < len(m.cmdInputs); i++ {
		inputs = append(inputs, m.cmdInputs[i].View())
	}

	var lineNum int
	limit := h - 10

	if m.index >= limit {
		lineNum = (m.index + 1) - limit
	} else {
		lineNum = 0
	}

	cnt := 0
	for i := lineNum; i < len(inputs); i++ {
		if cnt > limit-1 {
			continue
		}
		cnt++
		s += fmt.Sprintf("|%2d: %s\n", i, inputs[i])
	}

	s += colorSetting("+---------------------------------------------+\n", unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Execute selected line.  |\n", execKey), unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Save all lines.         |\n", saveKey), unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Remove current line.    |\n", delKey), unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Add a line at end.      |\n", addKey), unfocusedTextColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Exit.                   |\n", quitKey), unfocusedTextColor)
	s += colorSetting("+---------------------------------------------+\n", unfocusedTextColor)
	s += "\n"

	return s
}
