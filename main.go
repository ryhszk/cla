package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

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

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

var mdata = model{
	cursor:   1,
	choices:  []string{"ls -la", "free -h", "top", "./count", "dstat"},
	selected: make(map[int]struct{}),
}

func (m model) Init() tea.Cmd {
	return nil
}

var isexeccmd bool = false

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			isexeccmd = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Please select a command from next list.\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("[%s] %s\n", cursor, choice)
	}
	// s += fmt.Sprintf("%d\n", m.cursor)
	s += "\nPress q to quit.\n"
	return s
}

func (m model) getCursor() int {
	return m.cursor
}

func main() {
	p := tea.NewProgram(&mdata)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	if isexeccmd {
		cmd := mdata.choices[mdata.cursor]
		execCmd(cmd)
	}
}
