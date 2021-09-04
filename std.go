package bs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	initialVerboseEnvVar = "MAGEFILE_VERBOSE"
)

var (
	defaultStdin  io.Reader = os.Stdin
	defaultStdout io.Writer = os.Stdout
	defaultStderr io.Writer = os.Stderr

	echoFilters []string = []string{}

	colorsEnabled bool

	// defaults to Mage's verbose flag, since this package was original written to be used in Magefiles.
	// However, if you want to use your own VERBOSE flag here, just call SetVerboseEnvVarName.
	verboseEnvVar string = initialVerboseEnvVar
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

// SetVerboseEnvVarName allows changing the name of the environment variable that is used to
// decide if we are in Verbose mode. This function creates the new env var immediately,
// setting its value to true or false based on the value of the old env var name.
func SetVerboseEnvVarName(s string) {
	wasVerbose := IsVerbose()
	verboseEnvVar = s
	SetVerbose(wasVerbose)
}

func SetVerbose(b bool) {
	os.Setenv(verboseEnvVar, strconv.FormatBool(b))
}

func IsVerbose() bool {
	b, _ := strconv.ParseBool(os.Getenv(verboseEnvVar))
	return b
}

func Verbose(str string) {
	if !IsVerbose() {
		return
	}
	echo(str, ensureNewline, colorVerbose)
}

func Verbosef(format string, args ...interface{}) {
	if !IsVerbose() {
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
