package utils

import (
	"fmt"
	"os"
	"runtime"
)

func ErrorExit(err string) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fmt.Printf("call from '%s' function (line %d) \n", f.Name(), line)
	fmt.Printf("  err: %s\n", err)
	fmt.Print("  ")
	os.Exit(1)
}
