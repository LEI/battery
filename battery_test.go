package main

import (
	// "fmt"
	// "reflect"
	"testing"
	"time"

	"github.com/distatus/battery"
)

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
		&Battery{0, &battery.Battery{}, time.Duration(0)},
		"{{.Id}}",
		"BAT0",
	}, {
		&Battery{0, &battery.Battery{State: battery.Charging}, time.Duration(0)},
		"{{.State}}",
		"Charging",
	}, {
		&Battery{0, &battery.Battery{Current: 1, Full: 1}, time.Duration(0)},
		"{{.Spark}} {{.Percent}}%",
		"█ 100%",
	}, {
		&Battery{0, &battery.Battery{State: battery.Charging, Current: 1, Full: 2, ChargeRate: 2}, time.Duration(0)},
		"{{.Duration}}",
		"30m until charged",
	}}
	for i, c := range cases {
		str, err := c.battery.Parse(c.in)
		if err != nil {
			t.Errorf("%s", err)
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
