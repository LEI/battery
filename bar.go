package main

import (
	"fmt"

	"github.com/joliv/spark"
)

var BarLength = 8

// type Bar func(val float64, max float64) string

func asciiBar(val float64, max float64) string {
	var format = "[%s]"
	var sign = "="
	var barLen = 0
	if val > 0 {
		barLen = int(float64(BarLength) / (max / val))
	}
	lvlFmt := fmt.Sprintf("%%.%ds", BarLength)
	barLvl := fmt.Sprintf(lvlFmt, printN(sign, barLen))
	barFmt := fmt.Sprintf("%%- %ds", BarLength)
	bar := fmt.Sprintf(barFmt, barLvl)
	return fmt.Sprintf(format, bar)
	// barFmt := fmt.Sprintf("%%-%.0fs", barLen)
	// return fmt.Sprintf("[% 12s]", bar)
	// return fmt.Sprintf("[%-.8s]", "")
}

func printN(char string, times int) string {
	if times < 0 {
		panic(fmt.Errorf("printN times should not be negative: %d", times))
	}
	if times == 0 {
		return ""
	}
	str := fmt.Sprintf("%s", char)
	str += printN(char, times-1)
	return str
}

func sparkBar(val float64, max float64) string {
	sparkLine := spark.Line([]float64{0, val, max})
	runes := []rune(sparkLine)
	if len(runes) != 3 {
		panic(fmt.Errorf("invalid sparkline length (%d != 3): %s", len(runes), string(runes)))
	}
	return fmt.Sprintf("%s", string(runes[1]))
}
