package utils

import (
	"os"
	"os/exec"
	"runtime"
)

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

func Execute(cmd string) {
	c := exec.Command(shellName(), "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}
