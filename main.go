package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/distatus/battery"
	"github.com/spf13/pflag"
)

var (
	// Default battery template
	batteryTpl   = "{{.ID}}: {{.State}}, {{.Percent}}%{{if ne .Duration \"\"}}, {{end}}{{.Duration}}"
	outputSep    = "\n"
	outputFormat = "%s\n"
	colorFlag    bool
	sparkFlag    bool
	tmuxFlag     bool
	colors       Colors
)

var usageStr = `RTFM`

func init() {
	pflag.BoolVarP(&colorFlag, "color", "c", colorFlag, "Enable color output")
	pflag.BoolVarP(&sparkFlag, "spark", "s", sparkFlag, "Enable sparkline bar")
	pflag.BoolVarP(&tmuxFlag, "tmux", "t", tmuxFlag, "Enable tmux status bar colors")
	// pflag.StringVarP(&outputSep, "new-line", "n", outputSep, "Lines separator")
	// pflag.IntVarP(&limit, "limit", "l", limit, "Limit lines")
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s [flags] [format]\n", os.Args[0])
		pflag.PrintDefaults()
	}
	pflag.ErrHelp = errors.New(usageStr)
}

// ColorString according to format flags.
func ColorString(str string, clr string) string {
	var format = "%s%s%s"
	if tmuxFlag {
		format = "#[fg=%s]%s#[%s]"
	}
	return fmt.Sprintf(format, clr, str, colors.Get(DefaultColor))
}

// FormatDurationString in a human readable way.
func FormatDurationString(hours, minutes, seconds int) string {
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

// GetBar string in spark or ascii mode.
func GetBar(val float64, max float64) string {
	switch {
	case sparkFlag:
		return sparkBar(val, max)
	default:
		return asciiBar(val, max)
	}
}

// StateColorString according to battery status.
func StateColorString(bat *Battery) string {
	var clr string
	switch {
	case bat.IsEmpty():
		clr = colors.Get(EmptyColor)
	case bat.IsFull():
		clr = colors.Get(FullColor)
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
		clr = colors.Get(DefaultColor)
	}
	return clr
}

func findBatteries() ([]*battery.Battery, error) {
	batteries, err := GetAll()
	if err != nil {
		return batteries, err
	}
	if len(batteries) == 0 {
		return batteries, ErrNoBatteries
	}
	return batteries, nil
}

func main() {
	pflag.Parse()
	switch pflag.NArg() {
	case 0, 1:
		if pflag.Arg(0) != "" {
			batteryTpl = pflag.Arg(0)
		}
	default:
		fmt.Fprintf(os.Stderr, "invalid number of args: %d\n", pflag.NArg())
		os.Exit(1)
	}
	batteries, err := findBatteries()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch {
	case tmuxFlag:
		colors = &tmuxColors{}
	default: // case colorFlag:
		colors = &asciiColors{}
	}
	var out []string
	for i, bat := range batteries {
		b := &Battery{i, bat, 0}
		str := b.String()
		if colorFlag || tmuxFlag {
			str = ColorString(str, StateColorString(b))
		}
		out = append(out, str)
	}
	fmt.Printf(outputFormat, strings.Join(out, outputSep))
}
