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

type Mode int

const (
	_ Mode = iota
	Normal
	Edit
	Search
)

type model struct {
	cursor    int
	choice    chan string
	txtModels []txtinp.Model
	mode      Mode
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
	return model{0, ch, tms, Normal}
}

func (m model) Update(msg bubble.Msg) (bubble.Model, bubble.Cmd) {
	var bubbleCmd bubble.Cmd = nil

	switch msg := msg.(type) {
	case bubble.KeyMsg:
		switch msg.String() {

		case quitKey:
			close(m.choice)
			bubbleCmd = bubble.Quit

		case execKey:
			m.choice <- m.txtModels[m.cursor].Value()
			bubbleCmd = bubble.Quit

		case saveKey:
			m.mode = Normal
			var newDatas []util.JsonData
			var tmpData util.JsonData
			for i := 0; i < len(m.txtModels); i++ {
				tmpData.ID = i
				tmpData.CmdLine = m.txtModels[i].Value()
				newDatas = append(newDatas, tmpData)
			}
			newJsons, _ := json.Marshal(newDatas)
			util.ToFile(string(newJsons), dataFile)

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

		case delKey:
			// If there is only one line, deletion is prohibited.
			// (Since m.cursor starts at 0, adjust with len()-1)
			if m.cursor == 0 && m.cursor == len(m.txtModels)-1 {
				return m, nil
			}

			// Load from file again to avoid unintended saving.
			oldD := util.FromJSON(dataFile)
			newD := util.RmElem(oldD, m.cursor)
			newJ, _ := json.Marshal(newD)
			util.ToFile(string(newJ), dataFile)
			m.rmModel(m.cursor)
			// End of line case
			if m.cursor > len(m.txtModels)-1 {
				m.cursor--
			}
			m.cycleCursor()

		case "down", "tab":
			m.cursor++
			m.cycleCursor()

		case "up", "shift+tab":
			m.cursor--
			m.cycleCursor()

		case "ctrl+e":
			m.mode = Edit

		case "esc":
			m.mode = Normal

		default:
			// Handle character input and blinks
			if m.mode == Edit {
				m, bubbleCmd = updateInputs(msg, m)
			}
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

	if m.cursor > len(m.txtModels)-1 {
		m.cursor = 0
	} else if m.cursor < 0 {
		m.cursor = len(m.txtModels) - 1
	}

	for i := 0; i <= len(m.txtModels)-1; i++ {
		if i == m.cursor {
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
	adjh := h - 13
	inputs := []string{}
	for i := 0; i < len(m.txtModels); i++ {
		inputs = append(inputs, m.txtModels[i].View())
	}

	// header
	s := "\n"

	s += "+--------------+\n"
	s += "| "
	switch m.mode {
	case Normal:
		s += colorSetting("MODE: Normal ", unfocusedColor)
	case Edit:
		s += colorSetting("MODE: Edit   ", focusedColor)
	case Search:
		s += colorSetting("MODE: Search ", unfocusedColor)
	}
	s += "| \n"
	s += "+--------------+\n"
	// numDgt := strconv.Itoa(len(strconv.Itoa(limitLine)))
	// cstFmt := "Line: %" + numDgt + "d/%" + numDgt + "d"
	// s += fmt.Sprintf(cstFmt, m.index, len(m.txtModels)-1)
	// s += "\n+--------------+--------------+\n"

	// s += colorSetting("______________________________________________\n", unfocusedColor)

	// body
	totalLine := adjh
	var fstLine int
	if m.cursor >= totalLine {
		fstLine = (m.cursor + 1) - totalLine
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

	return s
}
