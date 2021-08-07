package bs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
)

var (
	defaultStdin  io.Reader = os.Stdin
	defaultStdout io.Writer = os.Stdout
	defaultStderr io.Writer = os.Stderr

	echoFilters []string = []string{}

	colorsEnabled bool
)

func init() {
	// TODO: check for terminal, or is this enough?
	colorsEnabled = len(os.Getenv("NO_COLOR")) == 0
}

// SetStdin overrides stdin for Run*/Bash*/Read* (defaults to os.Stdin)
func SetStdin(r io.Reader) {
	defaultStdin = r
}

// SetStdout overrides stdout for Run*/Bash* (defaults to os.Stdout)
func SetStdout(w io.Writer) {
	defaultStdout = w
}

// SetStderr overrides stderr for Run*/Bash* (defaults to os.Stderr)
func SetStderr(w io.Writer) {
	defaultStderr = w
}

func PushEchoFilter(str string) {
	echoFilters = append(echoFilters, str)
}

func PopEchoFilter() {
	echoFilters = echoFilters[:len(echoFilters)-1]
}

// Echo writes to stdout, and ensures the last character written is a newline.

func Echo(str string) {
	echo(str, ensureNewline, colorEcho)
}

func Echof(format string, args ...interface{}) {
	echo(fmt.Sprintf(format, args...), ensureNewline, colorEcho)
}

func Verbose(str string) {
	if !mg.Verbose() {
		return
	}
	echo(str, ensureNewline, colorVerbose)
}

func Verbosef(format string, args ...interface{}) {
	if !mg.Verbose() {
		return
	}
	echo(fmt.Sprintf(format, args...), ensureNewline, colorVerbose)
}

func Warn(str string) {
	echo(str, ensureNewline, colorWarn)
}

func Warnf(format string, args ...interface{}) {
	echo(fmt.Sprintf(format, args...), ensureNewline, colorWarn)
}

type echoOpt byte

const (
	ensureNewline echoOpt = iota
	ignoreFilter  echoOpt = iota
	colorEcho     echoOpt = iota
	colorVerbose  echoOpt = iota
	colorAsk      echoOpt = iota
	colorWarn     echoOpt = iota
)

func echo(str string, opts ...echoOpt) {
	newline := false
	filter := true
	var color string
	for _, v := range opts {
		switch v {
		case ensureNewline:
			newline = true
		case ignoreFilter:
			filter = false
		case colorEcho:
			color = ansiWhite
		case colorVerbose:
			color = ansiCyan
		case colorAsk:
			color = ansiBlue
		case colorWarn:
			color = ansiYellow
		}
	}

	if filter {
		for _, v := range echoFilters {
			str = strings.ReplaceAll(str, v, "******")
		}
	}

	if newline && str[len(str)-1] != '\n' {
		str += "\n"
	}

	if colorsEnabled && len(color) > 0 {
		str = color + str + ansiReset
	}

	fmt.Fprint(defaultStdout, str)
}

// ScanLine reads from default stdin until a newline is encountered

func ScanLine() string {
	str, err := ScanLineErr()
	if err != nil {
		fnErrorHandler(err)
	}
	return str
}

func ScanLineErr() (string, error) {
	r := bufio.NewReader(defaultStdin)
	str, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(str, "\n"), nil
}

// Ask is a combination of echo and scanline

func Ask(msg string) string {
	echo(msg, colorAsk)
	return ScanLine()
}

func Askf(format string, args ...interface{}) string {
	echo(fmt.Sprintf(format, args...), colorAsk)
	return ScanLine()
}

// ansi color helpers

const (
	ansiCSI   = "\u001b[" // Control Sequence Introducer
	ansiReset = ansiCSI + "39m"

	ansiBlack       = ansiCSI + "30m"
	ansiDarkRed     = ansiCSI + "31m"
	ansiDarkGreen   = ansiCSI + "32m"
	ansiDarkYellow  = ansiCSI + "33m"
	ansiDarkBlue    = ansiCSI + "34m"
	ansiDarkMagenta = ansiCSI + "35m"
	ansiDarkCyan    = ansiCSI + "36m"
	ansiLightGray   = ansiCSI + "37m"

	ansiDarkGray = ansiCSI + "90m"
	ansiRed      = ansiCSI + "91m"
	ansiGreen    = ansiCSI + "92m"
	ansiYellow   = ansiCSI + "93m"
	ansiBlue     = ansiCSI + "94m"
	ansiMagenta  = ansiCSI + "95m"
	ansiCyan     = ansiCSI + "96m"
	ansiWhite    = ansiCSI + "97m"
)
