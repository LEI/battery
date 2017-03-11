package main

import (
	// "fmt"
	// "reflect"
	"testing"
	"time"

	"github.com/distatus/battery"
)

// var defaultFormat = "{{.Id}}: {{.State}}, {{.Percent}}%{{if ne .Duration \"\"}}, {{end}}{{.Duration}}"
var oneMinute = time.Duration(time.Duration(1) * time.Minute)

func TestGetAll(t *testing.T) {
	batteries, err := GetAll()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(batteries) == 0 {
		t.Errorf("No batteries")
		return
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		battery *Battery
		in string
		out string
		// f func(int, *Battery, string) bool
	}{{
		&Battery{0, &battery.Battery{}, 0},
		"{{.Id}}",
		"BAT0",
	}, {
		&Battery{0, &battery.Battery{State: battery.Charging}, 0},
		"{{.State}}",
		"Charging",
	}, {
		&Battery{0, &battery.Battery{Current: 1, Full: 1}, 0},
		"{{.Spark}} {{.Percent}}%",
		"█ 100%",
	}, {
		&Battery{0, &battery.Battery{State: battery.Full, Current: 1, Full: 1, ChargeRate: 0}, 0},
		defaultFmt,
		"BAT0: Full, 100%",
	}, {
		&Battery{0, &battery.Battery{State: battery.Charging, Current: 1, Full: 2, ChargeRate: 2}, 0},
		"{{.Duration}}",
		"30m until charged",
	}, {
		&Battery{0, &battery.Battery{State: battery.Charging, Current: 1, Full: 1.01, ChargeRate: 2}, oneMinute},
		defaultFmt,
		"BAT0: Charging, 99%, 1m until charged",
	}, {
		&Battery{0, &battery.Battery{State: battery.Full, Current: 1, Full: 1, ChargeRate: 1}, oneMinute},
		"{{.Ftime \"%02d:%02d:%02d\"}}",
		"00:01:00",
	}}
	for i, c := range cases {
		str, err := c.battery.Parse(c.in)
		if err != nil {
			t.Errorf("%d: %s", i, err)
			return
		}
		// if colorOutput || tmuxOutput {
		// 	str = colorString(str, getStateColor(b))
		// }
		// if !reflect.DeepEqual(err, c.errorsOut) {
		// 	t.Errorf("%d: %v != %v", i, err, c.errorsOut)
		// }
		if str != c.out {
			t.Errorf("%d: %s != %s", i, str, c.out)
		}
	}
}
