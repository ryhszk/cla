package utils

import (
	"os"
	"os/exec"
	"runtime"
)

// ExecCmd is Execute a command (cmd) on shell.
// This function supports interactive commands (e.g. top, dstat) and,
// piped commands (e.g. echo hello | wc).
// No retun.
func ExecCmd(cmd string) {
	c := exec.Command(shellName(), "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}

func shellName() string {
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
