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
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/distatus/battery"
)

var (
	// ErrEmpty state
	ErrEmpty = errors.New("empty battery")
	// ErrFull state
	ErrFull = errors.New("full battery")
	// ErrIsNotCharging state
	ErrIsNotCharging = errors.New("will never fully charge")
	// ErrIsNotDischarging state
	ErrIsNotDischarging = errors.New("will never fully discharge")
	// ErrIsUnknown state
	ErrIsUnknown = errors.New("unkwnown state")
	// ErrNoBatteries detected
	ErrNoBatteries = errors.New("no batteries")
)

// GetAll batteries
func GetAll() ([]*battery.Battery, error) {
	batteries, err := battery.GetAll()
	return batteries, err
}

// func NewBattery(idx int, bat *battery.Battery) *Battery {
// 	return &Battery{idx, bat, time.Duration(0)}
// }

/* battery.Battery:
// Current battery state.
State State
// Current (momentary) capacity (in mWh).
Current float64
// Last known full capacity (in mWh).
Full float64
// Reported design capacity (in mWh).
Design float64
// Current (momentary) charge rate (in mW).
// It is always non-negative, consult .State field to check
// whether it means charging or discharging.
ChargeRate float64*/

// Battery type
type Battery struct {
	idx int
	*battery.Battery
	dur time.Duration
	// fmt string // template.Template
}

func (bat *Battery) String() string {
	str, err := bat.Template(batteryTpl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid battery %s", err)
		os.Exit(1)
	}
	return str // fmt.Sprintf("%s", str)
}

// Template battery string
func (bat *Battery) Template(tpl string) (string, error) {
	str, err := bat.Parse(tpl)
	if err != nil || str == "" {
		return str, err
	}
	return str, nil // fmt.Sprintf("%s", str), nil
}

// Parse battery template
func (bat *Battery) Parse(tpl string) (string, error) {
	tmpl, err := template.New("bat" + string(bat.idx)).Parse(tpl)
	// tmpl = tmpl.Option("missingkey=zero")
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, bat)
	str := buf.String()
	if err != nil {
		return str, err
	}
	return str, nil
}

// ID battery
func (bat *Battery) ID() string {
	return bat.Fid("BAT%d")
}

// Idx index
func (bat *Battery) Idx() int {
	return bat.idx
}

// Fid formats the battery index
func (bat *Battery) Fid(format string) string {
	return fmt.Sprintf(format, bat.idx)
}

// State string
func (bat *Battery) State() string {
	return bat.Battery.State.String()
}

// StateColor string
func (bat *Battery) StateColor() string {
	return StateColorString(bat)
}

// IsEmpty battery
func (bat *Battery) IsEmpty() bool {
	return bat.Battery.State == battery.Empty
}

// IsFull battery
func (bat *Battery) IsFull() bool {
	return bat.Battery.State == battery.Full
}

// IsCharging battery
func (bat *Battery) IsCharging() bool {
	return bat.Battery.State == battery.Charging
}

// IsDischarging battery
func (bat *Battery) IsDischarging() bool {
	return bat.Battery.State == battery.Discharging
}

// func (bat *Battery) ResetColor() string {
// 	return colors.Get(DefaultColor)
// }

// Percent remaining battery percent string
func (bat *Battery) Percent() string {
	return fmt.Sprintf("%.0f", bat.PercentFloat())
}

// func (bat *Battery) PercentInt() int {
// 	return int(fmt.PercentFloat())
// }

// PercentFloat remaining battery
func (bat *Battery) PercentFloat() float64 {
	return bat.Current / bat.Full * 100
}

// Fpercent formats remaining battery percent string
func (bat *Battery) Fpercent(format string) string {
	return fmt.Sprintf(format, bat.Percent())
}

// Bar string
func (bat *Battery) Bar() string {
	return GetBar(bat.Current, bat.Full)
}

// Duration string formatted or empty
func (bat *Battery) Duration() string {
	str, err := bat.Remaining()
	switch err {
	case nil: // continue
	case ErrEmpty, ErrFull, ErrIsNotCharging, ErrIsNotDischarging, ErrIsUnknown:
		return str
	}
	dur := FormatDurationString(bat.Hours(), bat.Minutes(), bat.Seconds())
	if dur != "" {
		str = fmt.Sprintf("%s %s", dur, str)
	}
	return str
}

// Remaining state string
func (bat *Battery) Remaining() (string, error) {
	var str string
	switch {
	case bat.IsEmpty():
		return "", ErrEmpty
	case bat.IsFull():
		return "", ErrFull
	case bat.IsCharging():
		if bat.ChargeRate == 0 {
			return "discharging at zero rate", ErrIsNotCharging
		}
		str = "until charged"
	case bat.IsDischarging():
		if bat.ChargeRate == 0 {
			return "discharging at zero rate", ErrIsNotDischarging
		}
		str = "remaining"
	default:
		return "unknown", ErrIsUnknown
	}
	return str, nil
}

// Time duration string
func (bat *Battery) Time() string {
	return bat.Ftime("%02d:%02d:%02d")
}

// Ftime formats duration string
func (bat *Battery) Ftime(format string) string {
	return fmt.Sprintf(format, bat.Hours(), bat.Minutes(), bat.Seconds())
}

// Charge rate
func (bat *Battery) Charge() float64 {
	var timeNum float64
	switch {
	case bat.IsCharging():
		if bat.ChargeRate == 0 {
			return -1
		}
		timeNum = (bat.Full - bat.Current) / bat.ChargeRate
	case bat.IsDischarging():
		if bat.ChargeRate == 0 {
			return -1
		}
		timeNum = bat.Current / bat.ChargeRate
	default:
		return 0
	}
	return timeNum
}

// ParseDuration battery
func (bat *Battery) ParseDuration() (time.Duration, error) {
	if bat.dur != time.Duration(0) {
		return bat.dur, nil
	}
	timeNum := bat.Charge()
	dur := fmt.Sprintf("%fh", timeNum)
	duration, err := time.ParseDuration(dur)
	if err != nil {
		return duration, err
	}
	bat.dur = duration
	return duration, nil
}

// Hours duration
func (bat *Battery) Hours() int {
	duration, err := bat.ParseDuration()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	h := duration.Hours()
	return int(h) % 60
}

// Fhours formats hours duration
func (bat *Battery) Fhours(format string) string {
	return fmt.Sprintf(format, bat.Hours())
}

// Minutes duration
func (bat *Battery) Minutes() int {
	duration, err := bat.ParseDuration()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	m := duration.Minutes()
	return int(m) % 60
}

// Fminutes formats minutes duration
func (bat *Battery) Fminutes(format string) string {
	return fmt.Sprintf(format, bat.Minutes())
}

// Seconds duration
func (bat *Battery) Seconds() int {
	duration, err := bat.ParseDuration()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	s := duration.Seconds()
	return int(s) % 60
}

// Fseconds formats seconds duration
func (bat *Battery) Fseconds(format string) string {
	return fmt.Sprintf(format, bat.Seconds())
}
