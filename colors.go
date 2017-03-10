package main

import (
	"fmt"
)

var (
	FullColor = "green"
	ChargingColor = "green"
	// Discharging
	HighColor = "none"
	MediumColor = "yellow"
	LowColor = "red"
)

var colorMap = map[string]string{
	"none": "0",
	"default": "0",
	"red": "0;31",
	"green": "0;32",
	"yellow": "0;33",
	"white": "0;37",
}

type Colors interface {
	Get(key string) string
}

type asciiColors struct {}

func (c *asciiColors) Get(key string) string {
	return fmt.Sprintf("\033[%sm", colorMap[key])
}

type tmuxColors struct {}

func (c *tmuxColors) Get(key string) string {
	return key
}

func colorString(str string, clr string) string {
	var format = "%s%s%s"
	if tmuxOutput {
		format = "#[fg=%s]%s#[%s]"
	}
	return fmt.Sprintf(format, clr, str, colors.Get("default"))
}
