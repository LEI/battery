package main

import (
	// "fmt"
	// "reflect"
	"testing"
	"time"

	"github.com/distatus/battery"
)

// var defaultFormat = "{{.ID}}: {{.State}}, {{.Percent}}%{{if ne .Duration \"\"}}, {{end}}{{.Duration}}"
var oneMinute = time.Duration(time.Duration(1) * time.Minute)

func TestGetAll(t *testing.T) {
	batteries, err := GetAll()
	if err != nil {
		// CI fails this test
		if err != ErrNoBatteries || len(batteries) != 0 {
			t.Error(err)
		}
	}
}

func TestParse(t *testing.T) {

	sparkFlag = true
	// tmuxFlag = true
	if colors == nil {
		colors = &tmuxColors{}
	}

	cases := []struct {
		battery *Battery
		in      string
		out     string
		// f func(int, *Battery, string) bool
	}{{
		&Battery{0, &battery.Battery{}, 0},
		"{{.ID}}",
		"BAT0",
	}, {
		&Battery{0, &battery.Battery{State: battery.Unknown}, 0},
		"{{.State}}",
		"Unknown",
	}, {
		&Battery{0, &battery.Battery{Current: 1, Full: 1}, 0},
		"{{.Bar}} {{.Percent}}%",
		"█ 100%",
	}, {
		&Battery{0, &battery.Battery{State: battery.Full, Current: 1, Full: 1, ChargeRate: 0}, 0},
		batteryTpl,
		"BAT0: Full, 100%",
	}, {
		&Battery{0, &battery.Battery{State: battery.Charging, Current: 1, Full: 2, ChargeRate: 2}, 0},
		"{{.Duration}}",
		"30m until charged",
	}, {
		&Battery{0, &battery.Battery{State: battery.Charging, Current: 1, Full: 1.01, ChargeRate: 2}, oneMinute},
		batteryTpl,
		"BAT0: Charging, 99%, 1m until charged",
	}, {
		&Battery{0, &battery.Battery{State: battery.Full, Current: 1, Full: 1, ChargeRate: 1}, oneMinute},
		"{{.Ftime \"%02d:%02d:%02d\"}}",
		"00:01:00",
	}, {
		&Battery{0, &battery.Battery{State: battery.Empty, Current: 0, Full: 1, ChargeRate: 0}, 0},
		"#[bg={{.StateColor}}]{{.Percent}}%#[default]",
		"#[bg=red]0%#[default]",
	}}
	for i, c := range cases {
		str, err := c.battery.Parse(c.in)
		// if !reflect.DeepEqual(err, c.errOut) {
		// 	t.Errorf("%d: %v != %v", i, err, c.errOut)
		// }
		if err != nil {
			t.Errorf("%d: %s", i, err)
			return
		}
		if colorFlag {
			str = ColorString(str, StateColorString(c.battery))
		}
		if str != c.out {
			t.Errorf("%d: %s != %s", i, str, c.out)
		}
	}
}
