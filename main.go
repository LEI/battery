// battery
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
	"strings"

	"github.com/spf13/pflag"
)

var (
	defaultFmt   = "{{.Id}}: {{.State}}, {{.Percent}}%{{if ne .Duration \"\"}}, {{end}}{{.Duration}}"
	outputSep    = "\n"
	outputFormat = "%s\n"
	tmuxOutput   bool
	colorOutput  bool
)

var colors = map[string]string{
	"green":   "0;32",
	"yellow":  "0;33",
	"red":     "0;31",
	"white":   "0;37",
	"default": "0",
	"none":    "0",
}

var (
	FullColor = "red"
	ChargingColor = "green"
	// Discharging
	HighColor = "none"
	MediumColor = "yellow"
	LowColor = "red"
)

type Adapter interface {
	// GetAll() ([]*battery.Battery, error)
	// GetPrimary() *battery.Battery
	GetHealth() string
	GetPercent() float64
	GetStatus() string // (dis)charging
	GetState() string  // (un)plugged
	// duration
	// temperature
}

// type Battery struct {}

// type Color struct {
// 	name string
// 	ascii string
// }
// func (c *Color) Wrap(str string) string {
// }

func init() {
	pflag.BoolVarP(&colorOutput, "color", "c", colorOutput, "Enable color output")
	pflag.BoolVarP(&tmuxOutput, "tmux", "t", tmuxOutput, "Enable tmux status bar colors")
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
	if tmuxOutput && !colorOutput {
		colorOutput = true
	}
	batteries, err := GetAll()
	if err != nil {
		exit(1, err)
	}
	if len(batteries) == 0 {
		exit(1, fmt.Errorf("No batteries"))
	}
	var out []string
	for i, bat := range batteries {
		b := New(i, bat)
		// str := b.String()
		str, err := b.Parse(defaultFmt)
		if err != nil {
			exit(1, err)
		}
		if colorOutput {
			str = colorString(str, b)
		}
		out = append(out, str)
	}
	fmt.Printf(outputFormat, strings.Join(out, outputSep))
}

func colorString(str string, bat *Battery) string {
	var format = "%s%s%s"
	if tmuxOutput {
		format = "#[fg=%s]%s#[%s]"
	}
	return fmt.Sprintf(format, getStateColor(bat), str, getColor("default"))
}

func getStateColor(bat *Battery) string {
	var clr string
	switch {
	case bat.IsCharging():
		clr = getColor(ChargingColor)
	case bat.IsDischarging():
		percent := bat.PercentFloat()
		switch {
		case percent >= 75:
			clr = getColor(HighColor)
		case percent >= 25: // && percent < 75:
			clr = getColor(MediumColor)
		case percent < 25:
			clr = getColor(LowColor)
		}
	default:
		clr = getColor("full")
	}
	return clr
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

func exit(code int, msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}
