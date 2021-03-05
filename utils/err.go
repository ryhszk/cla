package utils

import (
	"fmt"
	"os"
	"runtime"
)

// ErrExit is Print error (err) and, call to os.Exit(1).
// The printout includes the name of the calling function,
// the number of lines, and the error (err) that was passed.
// No return.
func ErrExit(err string) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fmt.Printf("call from '%s' function (line %d) \n", f.Name(), line)
	fmt.Printf("  err: %s\n", err)
	fmt.Print("  ")
	os.Exit(1)
}
