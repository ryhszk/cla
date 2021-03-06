package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"

	txtinp "github.com/charmbracelet/bubbles/textinput"
	bubble "github.com/charmbracelet/bubbletea"
	ter "github.com/muesli/termenv"

	cfg "github.com/ryhszk/cla/config"
	util "github.com/ryhszk/cla/utils"

	"golang.org/x/term"
)

type model struct {
	index     int
	choice    chan string
	txtModels []txtinp.Model
}

var (
	color          = ter.ColorProfile().Color
	focusedPrompt  = colorSetting("> ", focusedColor)
	blurredPrompt  = "  "
	focusedColor   = cfg.Config.FocusedColor
	unfocusedColor = cfg.Config.UnfocusedColor
	dataFile       = os.Getenv("HOME") + "/.cla/" + cfg.Config.DataFile
	limitLine      = cfg.Config.LimitLine
	execKey        = cfg.Config.ExecKey
	saveKey        = cfg.Config.SaveKey
	delKey         = cfg.Config.DelKey
	addKey         = cfg.Config.AddKey
	quitKey        = cfg.Config.QuitKey
)

func colorSetting(srcStr, colorCode string) string {
	return ter.String(srcStr).Foreground(color(colorCode)).String()
}

func main() {
	result := make(chan string, 1)

	if err := bubble.NewProgram(initialModel(result)).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	var cmdLine string
	if cmdLine = <-result; cmdLine != "" {
		util.ExecCmd(cmdLine)
	}
}

func (m model) Init() bubble.Cmd {
	return txtinp.Blink
}

func initialModel(ch chan string) model {
	tms := []txtinp.Model{}
	for i, j := range util.FromJSON(dataFile) {
		tm := txtinp.NewModel()
		tm.Width = 99
		tm.CharLimit = 99
		tm.SetValue(j.CmdLine)
		tm.Placeholder = "Input any command."
		if i != 0 {
			tm.Prompt = blurredPrompt
		} else {
			tm.Focus()
			tm.Prompt = focusedPrompt
			tm.TextColor = focusedColor
		}
		tms = append(tms, tm)
	}
	return model{0, ch, tms}
}

func (m model) Update(msg bubble.Msg) (bubble.Model, bubble.Cmd) {
	var bubbleCmd bubble.Cmd

	switch msg := msg.(type) {
	case bubble.KeyMsg:
		switch msg.String() {

		case quitKey:
			close(m.choice)
			bubbleCmd = bubble.Quit

		case execKey:
			m.choice <- m.txtModels[m.index].Value()
			bubbleCmd = bubble.Quit

		case saveKey:
			var newDatas []util.JsonData
			var tmpData util.JsonData
			for i := 0; i < len(m.txtModels); i++ {
				tmpData.ID = i
				tmpData.CmdLine = m.txtModels[i].Value()
				newDatas = append(newDatas, tmpData)
			}
			newJsons, _ := json.Marshal(newDatas)
			util.ToFile(string(newJsons), dataFile)
			bubbleCmd = nil

		case addKey:
			if len(m.txtModels) >= limitLine {
				return m, nil
			}
			newDatas := util.FromJSON(dataFile)
			tailNumber := len(m.txtModels)
			emptyData := util.JsonData{tailNumber, ""}
			newDatas = append(newDatas, emptyData)
			newJsons, _ := json.Marshal(newDatas)
			util.ToFile(string(newJsons), dataFile)
			m.addModel()
			bubbleCmd = nil

		case delKey:
			if m.index == 0 {
				return m, nil
			}
			// Load from file again to avoid unintended saving.
			oldD := util.FromJSON(dataFile)
			newD := util.RmElem(oldD, m.index)
			newJ, _ := json.Marshal(newD)
			util.ToFile(string(newJ), dataFile)
			m.rmModel(m.index)
			// End of line case
			if m.index > len(m.txtModels)-1 {
				m.index--
			}
			m.cycleCursor()
			bubbleCmd = nil

		case "down", "tab":
			m.index++
			m.cycleCursor()
			bubbleCmd = nil

		case "up", "shift+tab":
			m.index--
			m.cycleCursor()
			bubbleCmd = nil

		default:
			// Handle character input and blinks
			m, bubbleCmd = updateInputs(msg, m)
		}
	}

	return m, bubbleCmd
}

// Pass messages and models through to text input components. Only text inputs
// with Focus() set will respond, so it's safe to simply update all of them
// here without any further logic.
func updateInputs(msg bubble.Msg, m model) (model, bubble.Cmd) {
	var (
		bubbleCmd  bubble.Cmd
		bubbleCmds []bubble.Cmd
	)

	for i := 0; i < len(m.txtModels); i++ {
		m.txtModels[i], bubbleCmd = m.txtModels[i].Update(msg)
		bubbleCmds = append(bubbleCmds, bubbleCmd)
	}

	return m, bubble.Batch(bubbleCmds...)
}

func (m *model) cycleCursor() {

	if m.index > len(m.txtModels)-1 {
		m.index = 0
	} else if m.index < 0 {
		m.index = len(m.txtModels) - 1
	}

	for i := 0; i <= len(m.txtModels)-1; i++ {
		if i == m.index {
			// Set focused state
			m.txtModels[i].Focus()
			m.txtModels[i].Prompt = focusedPrompt
			m.txtModels[i].TextColor = focusedColor
			continue
		}
		// Remove focused state
		m.txtModels[i].Blur()
		m.txtModels[i].Prompt = blurredPrompt
		m.txtModels[i].TextColor = ""
	}
}

func (m *model) addModel() {
	tm := txtinp.NewModel()
	tm.Width = 99
	tm.CharLimit = 99
	tm.SetValue("")
	tm.Placeholder = "Input any command."
	tm.Prompt = blurredPrompt
	m.txtModels = append(m.txtModels, tm)
}

func (m *model) rmModel(i int) {
	if i >= len(m.txtModels) {
		return
	}
	m.txtModels = append(m.txtModels[:i], m.txtModels[i+1:]...)
}

func (m model) View() string {
	// initial
	_, h, _ := term.GetSize(syscall.Stdin)
	inputs := []string{}
	for i := 0; i < len(m.txtModels); i++ {
		inputs = append(inputs, m.txtModels[i].View())
	}

	// header
	s := "\n"
	s += colorSetting("______________________________________________\n", unfocusedColor)

	// body
	totalLine := h - 12
	var fstLine int
	if m.index >= totalLine {
		fstLine = (m.index + 1) - totalLine
	} else {
		fstLine = 0
	}

	numDgt := strconv.Itoa(len(strconv.Itoa(limitLine)))
	cstFmt := "|%" + numDgt + "d: %s\n"
	lineCnt := 0
	for i := fstLine; i < len(inputs); i++ {
		if lineCnt > totalLine-1 {
			continue
		}
		lineCnt++
		s += fmt.Sprintf(cstFmt, i, inputs[i])
	}

	// footer
	s += colorSetting("+---------------------------------------------+\n", unfocusedColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Exit.                   |\n", quitKey), unfocusedColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Execute selected line.  |\n", execKey), unfocusedColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Save all lines.         |\n", saveKey), unfocusedColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Add a line at end.      |\n", addKey), unfocusedColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Remove current line.    |\n", delKey), unfocusedColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Move down.              |\n", "↓ [tab]"), unfocusedColor)
	s += colorSetting(fmt.Sprintf("| %-17s | Move up.                |\n", "↑ [shift+tab]"), unfocusedColor)
	s += colorSetting("+---------------------------------------------+\n", unfocusedColor)
	s += "\n"

	return s
}
