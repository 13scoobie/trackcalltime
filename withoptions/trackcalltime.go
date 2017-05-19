package trackcalltime

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// Define a global regex for extracting function names. Thanks to
// Jeff below for pointing out that I was calling MustCompile() from
// the function each time.
var RE_stripFnPreamble = regexp.MustCompile(`^.*\.(.*)$`)

type Options struct {
	// Setting "DisableNesting" to "true" will cause tracey to not indent
	// any messages from nested functions. The default value is "false"
	// which enables nesting by prepending "SpacesPerIndent" number of
	// spaces per level nested.
	DisableNesting  bool
	SpacesPerIndent int

	// Private member, used to keep track of how many levels of nesting
	// the current trace functions have navigated.
	currentDepth int
}

// Single entry-point to fetch trace functions
func New(opts *Options) (func(string), func() string) {

	var options Options
	if opts != nil {
		options = *opts
	}

	// If nesting is enabled, and the spaces are not specified,
	// use the "default" value of 2
	if options.DisableNesting {
		options.SpacesPerIndent = 0
	} else if options.SpacesPerIndent == 0 {
		options.SpacesPerIndent = 4
	}

	_incrementDepth := func() {
		options.currentDepth += 1
	}

	_decrementDepth := func() {
		options.currentDepth -= 1
		if options.currentDepth < 0 {
			panic("Depth is negative! Should never happen!")
		}
	}

	_spacify := func() string {
		return strings.Repeat(" ", options.currentDepth*options.SpacesPerIndent)
	}

	// Define our enter function
	_enter := func() string {
		defer _incrementDepth()

		fnName := "<unknown>"
		// Skip this function, and fetch the PC and file for its parent
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			fnName = RE_stripFnPreamble.ReplaceAllString(runtime.FuncForPC(pc).Name(), "$1")
		}

		fmt.Printf("%sEntering %s\n", _spacify(), fnName)
		return fnName
	}

	// Define the exit function
	_exit := func(s string) {
		_decrementDepth()
		fmt.Printf("%sExiting  %s\n", _spacify(), s)
	}

	// Return the trace functions to the caller
	return _exit, _enter
}
