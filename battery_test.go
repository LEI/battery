package main

import (
	"fmt"
	// "reflect"
	"testing"
)

func TestParse(t *testing.T) {
	batteries, err := GetAll()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(batteries) == 0 {
		t.Errorf("No batteries")
		return
	}
	cases := []struct {
		in  string // template
		out func(int, *Battery, string) bool

	}{{
		"{{.Id}}", func(i int, b *Battery, r string) bool { return fmt.Sprintf("BAT%d", i) == r },
	},
	{
		"{{.Spark}}", func(i int, b *Battery, r string) bool { return len(r) > 0 },
	}}
	for i, c := range cases {
		for j, bat := range batteries {
			b := New(j, bat)
			str, err := b.Parse(c.in)
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
			if !c.out(j, b, str) {
				t.Errorf("%d: %s !", i, str)
			}
		}
	}
}
