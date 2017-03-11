package main

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/distatus/battery"
)

func GetAll() ([]*battery.Battery, error) {
	batteries, err := battery.GetAll()
	return batteries, err
}

type Battery struct {
	idx int
	*battery.Battery
	dur time.Duration
	// fmt string // template.Template
}

func New(idx int, bat *battery.Battery) *Battery {
	return &Battery{idx, bat, time.Duration(0)}
}

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

func (bat *Battery) Id() string {
	return bat.Fid("BAT%d")
}

func (bat *Battery) Idx() int {
	return bat.idx
}

func (bat *Battery) Fid(format string) string {
	return fmt.Sprintf(format, bat.idx)
}

func (bat *Battery) State() string {
	return bat.Battery.State.String()
}

func (bat *Battery) IsCharging() bool {
	return bat.Battery.State == battery.Charging
}

func (bat *Battery) IsDischarging() bool {
	return bat.Battery.State == battery.Discharging
}

func (bat *Battery) Percent() string {
	return fmt.Sprintf("%.0f", bat.PercentFloat())
}

// func (bat *Battery) PercentInt() int {
// 	return int(fmt.PercentFloat())
// }

func (bat *Battery) PercentFloat() float64 {
	return bat.Current / bat.Full * 100
}

func (bat *Battery) Fpercent(format string) string {
	return fmt.Sprintf(format, bat.Percent())
}

func (bat *Battery) Spark() string {
	return sparkBar(bat.Current, bat.Full)
}

func (bat *Battery) Duration() string {
	var str string
	switch {
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
		return "not charging, not discharging"
	}
	dur := formatDuration(bat)
	if dur != "" {
		str = fmt.Sprintf("%s %s", dur, str)
	}
	return str
}

func (bat *Battery) ParseDuration() time.Duration {
	if bat.dur != time.Duration(0) {
		return bat.dur
	}
	var timeNum float64
	switch {
	case bat.IsCharging():
		if bat.ChargeRate == 0 {
			return time.Duration(-1)
		}
		timeNum = (bat.Full - bat.Current) / bat.ChargeRate
	case bat.IsDischarging():
		if bat.ChargeRate == 0 {
			return time.Duration(-1)
		}
		timeNum = bat.Current / bat.ChargeRate
	default:
		return time.Duration(-1)
	}
	dur := fmt.Sprintf("%fh", timeNum)
	duration, err := time.ParseDuration(dur)
	if err != nil {
		panic(err)
	}
	bat.dur = duration
	return duration
}

func (bat *Battery) Hours() int64 {
	duration := bat.ParseDuration()
	h := duration.Hours()
	return int64(h) % 60 // int64(time.Hour / time.Minute)
}

func (bat *Battery) Fhours(format string) string {
	return fmt.Sprintf(format, bat.Hours())
}

func (bat *Battery) Minutes() int64 {
	duration := bat.ParseDuration()
	m := duration.Minutes()
	return int64(m) % 60 // int64(time.Minute / time.Second)
}

func (bat *Battery) Fminutes(format string) string {
	return fmt.Sprintf(format, bat.Minutes())
}

func (bat *Battery) Seconds() int64 {
	duration := bat.ParseDuration()
	s := duration.Seconds()
	return int64(s) % 60 // int64(time.Minute / time.Second)
}

func (bat *Battery) Fseconds(format string) string {
	return fmt.Sprintf(format, bat.Seconds())
}
