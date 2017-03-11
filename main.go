package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/joliv/spark"
)

var (
	defaultFmt   = "{{.Id}}: {{.State}}, {{.Percent}}%{{if ne .Duration \"\"}}, {{end}}{{.Duration}}"
	outputSep    = "\n"
	outputFormat = "%s\n"
	tmuxOutput   bool
	colorOutput  bool
	colors       Colors
)

func init() {
	pflag.BoolVarP(&colorOutput, "color", "c", colorOutput, "Enable color output")
	pflag.BoolVarP(&tmuxOutput, "tmux", "t", tmuxOutput, "Enable tmux status bar colors")
	// pflag.StringVarP(&outputSep, "new-line", "n", outputSep, "Lines separator")
	// pflag.BoolVarP(&asciiBar, "ascii", "a", asciiBar, "Enable ascii bar left to percentage")
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
	case colorOutput:
		colors = &asciiColors{}
	case tmuxOutput:
		colors = &tmuxColors{}
	}
	var out []string
	for i, bat := range batteries {
		b := New(i, bat)
		// str := b.String()
		str, err := b.Parse(defaultFmt)
		if err != nil {
			exit(1, err)
		}
		if colorOutput || tmuxOutput {
			str = colorString(str, getStateColor(b))
		}
		out = append(out, str)
	}
	fmt.Printf(outputFormat, strings.Join(out, outputSep))
}

func getStateColor(bat *Battery) string {
	var clr string
	switch {
	case bat.IsCharging():
		clr = colors.Get(ChargingColor)
	case bat.IsDischarging():
		percent := bat.PercentFloat()
		switch {
		case percent >= 75:
			clr = colors.Get(HighColor)
		case percent >= 25: // && percent < 75:
			clr = colors.Get(MediumColor)
		case percent < 25:
			clr = colors.Get(LowColor)
		}
	default:
		clr = colors.Get(FullColor)
	}
	return clr
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

func sparkBar(val float64, max float64) string {
	sparkLine := spark.Line([]float64{0, val, max})
	runes := []rune(sparkLine)
	if len(runes) != 3 {
		panic(fmt.Errorf("invalid sparkline length (%d != 3): %s", len(runes), string(runes)))
	}
	return fmt.Sprintf("%s", string(runes[1]))
}

func exit(code int, msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}
