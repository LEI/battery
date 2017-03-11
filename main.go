package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var (
	defaultFmt   = "{{.Id}}: {{.State}}, {{.Percent}}%{{if ne .Duration \"\"}}, {{end}}{{.Duration}}"
	outputSep    = "\n"
	outputFormat = "%s\n"
	colorFlag    bool
	sparkFlag    bool
	tmuxFlag     bool
	colors       Colors
)

func init() {
	pflag.BoolVarP(&colorFlag, "color", "c", colorFlag, "Enable color output")
	pflag.BoolVarP(&sparkFlag, "spark", "s", sparkFlag, "Enable sparkline bar")
	pflag.BoolVarP(&tmuxFlag, "tmux", "t", tmuxFlag, "Enable tmux status bar colors")
	// pflag.StringVarP(&outputSep, "new-line", "n", outputSep, "Lines separator")
	// pflag.IntVarP(&limit, "limit", "l", limit, "Limit lines")
}

func main() {
	pflag.Parse()
	switch pflag.NArg() {
	case 0, 1:
		if pflag.Arg(0) != "" {
			defaultFmt = pflag.Arg(0)
		}
	default:
		exit(1, fmt.Sprintf("Invalid number of args: %d", pflag.NArg()))
	}
	batteries, err := GetAll()
	if err != nil {
		exit(1, err)
	}
	if len(batteries) == 0 {
		exit(1, fmt.Errorf("No batteries"))
	}
	switch {
	case colorFlag:
		colors = &asciiColors{}
	case tmuxFlag:
		colors = &tmuxColors{}
	}
	var out []string
	for i, bat := range batteries {
		b := &Battery{i, bat, 0}
		// str := b.String()
		str, err := b.Parse(defaultFmt)
		if err != nil {
			exit(1, err)
		}
		if colorFlag || tmuxFlag {
			str = ColorString(str, b.StateColor())
		}
		out = append(out, str)
	}
	fmt.Printf(outputFormat, strings.Join(out, outputSep))
}

func ColorString(str string, clr string) string {
	var format = "%s%s%s"
	if tmuxFlag {
		format = "#[fg=%s]%s#[%s]"
	}
	return fmt.Sprintf(format, clr, str, colors.Get(DefaultColor))
}

func GetBar(val float64, max float64) string {
	switch {
	case sparkFlag:
		return sparkBar(val, max)
	default:
		return asciiBar(val, max)
	}
}

func formatTime(hours, minutes, seconds int) string {
	var str string // Pad int with zero: %02d, Truncate string: %.0s
	switch 0 {
	case hours + minutes + seconds:
		return "" // fully charged
	case hours + minutes:
		str = fmt.Sprintf("%ds", seconds)
	case hours:
		str = fmt.Sprintf("%dm", minutes)
	default:
		str = fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return str
}

func exit(code int, msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}
