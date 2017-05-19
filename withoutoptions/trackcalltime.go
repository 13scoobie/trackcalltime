package trackcalltime

import (
	"fmt"
	"regexp"
	"runtime"
)

// Regex to extract just the function name (and not the module path)
var RE_stripFnPreamble = regexp.MustCompile(`^.*\.(.*)$`)

func Enter() string {
	fnName := "<unknown>"
	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		fnName = RE_stripFnPreamble.ReplaceAllString(runtime.FuncForPC(pc).Name(), "$1")
	}

	fmt.Printf("Entering %s\n", fnName)
	return fnName
}

func Exit(s string) {
	fmt.Printf("Exiting  %s\n", s)
}
