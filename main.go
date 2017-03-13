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
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s [flags] [format]\n", os.Args[0])
		pflag.PrintDefaults()
	}
}

func main() {
	pflag.Parse()
	switch pflag.NArg() {
	case 0, 1:
		if pflag.Arg(0) != "" {
			defaultFmt = pflag.Arg(0)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid number of args: %d\n", pflag.NArg())
		os.Exit(1)
	}
	batteries, err := GetAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	if len(batteries) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", "No batteries")
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
		if colorFlag {
			str = ColorString(str, StateColorString(b))
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

func StateDurationString(bat *Battery) string {
	var str string
	switch {
	case bat.IsEmpty(), bat.IsFull():
		return ""
	case bat.IsCharging():
		if bat.ChargeRate == 0 {
			return "charging at zero rate - will never fully charge"
		}
		str = "until charged"
	case bat.IsDischarging():
		if bat.ChargeRate == 0 {
			return "discharging at zero rate - will never fully discharge"
		}
		str = "remaining"
	default:
		return "unknown state"
	}
	dur := FormatDurationString(bat.Hours(), bat.Minutes(), bat.Seconds())
	if dur != "" {
		str = fmt.Sprintf("%s %s", dur, str)
	}
	return str
}

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

func GetBar(val float64, max float64) string {
	switch {
	case sparkFlag:
		return sparkBar(val, max)
	default:
		return asciiBar(val, max)
	}
}

// func exit(code int, msg interface{}) {
// 	fmt.Fprintln(os.Stderr, msg)
// 	os.Exit(code)
// }
