package main

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/distatus/battery"
	"github.com/joliv/spark"
)

var (
	Charging = battery.Charging
	Discharging = battery.Discharging
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
	sparkLine := spark.Line([]float64{0, bat.Current, bat.Full})
	runes := []rune(sparkLine)
	if len(runes) != 3 {
		panic(fmt.Errorf("invalid sparkline length (%d != 3): %s", len(runes), string(runes)))
	}
	return fmt.Sprintf("%s", string(runes[1]))
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
	hours, minutes, seconds := bat.Hours(), bat.Minutes(), bat.Seconds()
	switch 0 { // %02d
	case hours + minutes + seconds:
		return "" // fully charged
	case hours + minutes:
		str = fmt.Sprintf("%ds %s", seconds, str)
	case hours:
		str = fmt.Sprintf("%dm %s", minutes, str)
	default:
		str = fmt.Sprintf("%dh%dm %s", hours, minutes, str)
	}
	return fmt.Sprintf("%s", str)
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
	case bat.IsCharging():
		if bat.ChargeRate == 0 {
			return time.Duration(-1)
		}
		timeNum = bat.Current / bat.ChargeRate
	default:
		return time.Duration(0)
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%fh", timeNum))
	if err != nil {
		panic(err)
	}
	bat.dur = duration
	return duration
}

func (bat *Battery) Hours() int {
	duration := bat.ParseDuration()
	h := duration.Hours()
	return int(h) % int(time.Hour)
}

func (bat *Battery) Fhours(format string) string {
	return fmt.Sprintf(format, bat.Hours())
}

func (bat *Battery) Minutes() int {
	duration := bat.ParseDuration()
	m := duration.Minutes()
	return int(m) % int(time.Minute)
}

func (bat *Battery) Fminutes(format string) string {
	return fmt.Sprintf(format, bat.Minutes())
}

func (bat *Battery) Seconds() int {
	duration := bat.ParseDuration()
	s := duration.Seconds()
	return int(s) % int(time.Second)
}

func (bat *Battery) Fseconds(format string) string {
	return fmt.Sprintf(format, bat.Seconds())
}
