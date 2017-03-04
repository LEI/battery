// Extended battery CLI
// Copyright (C) 2016 Karol 'Kenji Takahashi' Woźniak
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/distatus/battery"
	"github.com/joliv/spark"
	"github.com/spf13/pflag"
)

type Option struct {
	format string
	flag   bool
}

var (
	optSep       = ", "
	lineSep      = "\n"
	formatOutput = "%s"
	colorOutput  bool
	tmuxOutput   bool
	sparkLine    bool
	opts         = map[string]*Option{
		"id":       {"BAT%d: ", false},
		"state":    {"%s", false},
		"percent":  {"%.2f%%", false},
		"duration": {"%dh%dm %s", false},
	}
	colors = map[string]string{
		"green":   "0;32",
		"yellow":  "0;33",
		"red":     "0;31",
		"white":   "0;37",
		"default": "0",
		"none":    "0",
	}
	states = map[string]string{
		"charging": "none",
		"high":     "green",
		"medium":   "yellow",
		"low":      "red",
	}
)

func init() {
	pflag.BoolVarP(&colorOutput, "color", "c", colorOutput, "Enable color output")
	pflag.BoolVarP(&tmuxOutput, "tmux", "t", tmuxOutput, "Enable tmux status bar colors")
	pflag.BoolVarP(&sparkLine, "spark", "", sparkLine, "Enable sparkline left to percentage")
	// pflag.BoolVarP(&asciiBar, "ascii", "a", asciiBar, "Enable ascii bar left to percentage")
	// pflag.IntVarP(&limit, "limit", "l", limit, "Limit lines")

	pflag.BoolVarP(&opts["duration"].flag, "duration", "d", opts["duration"].flag, "Print time until (dis)charged or charge rate status")
	pflag.BoolVarP(&opts["id"].flag, "id", "i", opts["id"].flag, "Battery identifier")
	pflag.BoolVarP(&opts["percent"].flag, "percent", "p", opts["percent"].flag, "Print remaingin charge as a percentage")
	pflag.BoolVarP(&opts["state"].flag, "state", "s", opts["state"].flag, "Print state (Charging, Discharging)")

	pflag.StringVarP(&opts["duration"].format, "dfmt", "", opts["duration"].format, "Format duration")
	pflag.StringVarP(&opts["id"].format, "ifmt", "", opts["id"].format, "Format battery number")
	pflag.StringVarP(&opts["percent"].format, "pfmt", "", opts["percent"].format, "Format percentage")
	pflag.StringVarP(&opts["state"].format, "sfmt", "", opts["state"].format, "Format state")
}

func main() {
	pflag.Parse()
	nbArg := pflag.NArg()
	if nbArg > 1 {
		exit(1, fmt.Sprintf("invalid number of args: %d", nbArg))
	} else if nbArg == 1 && pflag.Arg(0) != "" {
		formatOutput = pflag.Arg(0)
	}
	nbFlag := pflag.NFlag()

	if nbFlag == 0 || (opts["id"].flag && !opts["state"].flag && !opts["percent"].flag && !opts["duration"].flag) {
		// Print full battery info
		opts["id"].flag = true
		opts["state"].flag = true
		opts["percent"].flag = true
		opts["duration"].flag = true
	}
	batteries, err := battery.GetAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(batteries) == 0 {
		fmt.Fprintln(os.Stderr, "No batteries")
		os.Exit(1)
	}
	var lines []string
	for i, bat := range batteries {
		str := getBatteryString(i, bat)
		lines = append(lines, str)
	}
	fmt.Printf(formatOutput+"\n", strings.Join(lines, lineSep))
}

func getBatteryString(idx int, bat *battery.Battery) string {
	var out []string
	for _, key := range []string{"state", "percent", "duration"} {
		val := opts[key]
		if !val.flag {
			continue
		}
		var opt string
		switch key {
		case "state":
			opt = fmt.Sprintf(val.format, bat.State)
		case "percent":
			opt = fmt.Sprintf(val.format, bat.Current/bat.Full*100)
			if sparkLine {
				sl := spark.Line([]float64{0, bat.Current, bat.Full})
				runes := []rune(sl)
				if len(runes) != 3 {
					panic(fmt.Errorf("invalid sparkline lendth (%d != 3): %s", len(runes), string(runes)))
				}
				opt = fmt.Sprintf("%s %s", string(runes[1]), opt)
			}
		case "duration":
			batteryDuration := durationFormat(val.format, bat)
			if opts["state"].flag && batteryDuration == "fully charged" {
				// Hide duration, battery state is already 'Full'
				continue
			}
			opt = batteryDuration
			// default:
			// 	opt += fmt.Sprint(val.format)
		}
		if opt != "" {
			out = append(out, opt)
		}
	}
	str := strings.Join(out, optSep)
	if opts["id"].flag {
		str = fmt.Sprintf(opts["id"].format+"%s", idx, str)
	}
	if colorOutput {
		str = applyColors(str, bat)
	}
	return str
}

func getColor(key string) string {
	switch key {
	default:
		if !tmuxOutput {
			key = fmt.Sprintf("\033[%sm", colors[key])
		}
	}
	return key
}

func applyColors(str string, bat *battery.Battery) string {
	var clr string = getColor("none")
	switch bat.State {
	case battery.Charging:
		clr = getColor(states["charging"])
	case battery.Discharging:
		percent := bat.Current / bat.Full * 100
		switch {
		case percent >= 75:
			clr = getColor(states["high"])
		case percent >= 25: // && percent < 75:
			clr = getColor(states["medium"])
		case percent < 25:
			clr = getColor(states["low"])
		}
	}
	format := "%s%s%s"
	if tmuxOutput {
		format = "#[fg=%s]%s#[%s]"
	}
	return fmt.Sprintf(format, clr, str, getColor("default"))
}

func durationFormat(format string, bat *battery.Battery) string {
	var str string
	var timeNum float64
	switch bat.State {
	case battery.Charging:
		if bat.ChargeRate == 0 {
			return "charging at zero rate - will never fully charge"
		}
		str = "until charged"
		timeNum = (bat.Full - bat.Current) / bat.ChargeRate
	case battery.Discharging:
		if bat.ChargeRate == 0 {
			return "discharging at zero rate - will never fully discharge"
		}
		str = "remaining"
		timeNum = bat.Current / bat.ChargeRate
	default: // Full charge
		return "fully charged"
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%fh", timeNum))
	if err != nil {
		exit(1, err)
	}
	hours, err := extractTime(duration, "h", "")
	if err != nil {
		exit(1, err)
	}
	minutes, err := extractTime(duration, "m", "h")
	if err != nil {
		exit(1, err)
	}
	// fmt.Printf("> timeNum: %f\n> duration: %s\n> format:%s\n", timeNum, duration, format)
	return fmt.Sprintf(format, hours, minutes, str)
}

func extractTime(duration time.Duration, unit string, after string) (int, error) {
	var str = "0"
	var from, to int
	if after != "" {
		from = strings.Index(duration.String(), after)
	}
	to = strings.Index(duration.String(), unit)
	if from > 0 && to > 0 {
			str = duration.String()[from+1:to]
	} else if to > 0 {
		str = duration.String()[:to]
	}
	integer, err := strconv.Atoi(str)
	return integer, err
}

func exit(code int, msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}
