package main

import (
	"fmt"
)

var (
	DefaultColor  = "default"
	EmptyColor    = "red"
	FullColor     = "green"
	ChargingColor = "green"
	HighColor     = "none"   // Discharging, over 75%
	MediumColor   = "yellow" // Discharging, between 25 and 75%
	LowColor      = "red"    // Discharging, below 25%
)

var colorMap = map[string]string{
	"default": "0",
	"none":    "0",
	"red":     "0;31",
	"green":   "0;32",
	"yellow":  "0;33",
	"white":   "0;37",
}

type Colors interface {
	Get(key string) string
}

type asciiColors struct{}

func (c *asciiColors) Get(key string) string {
	clr, ok := colorMap[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("\033[%sm", clr)
}

type tmuxColors struct{}

func (c *tmuxColors) Get(key string) string {
	return key
}
