package main

import (
	"fmt"
)

var (
	DefaultColor  = "none"
	EmptyColor    = "red"
	FullColor     = "green"
	ChargingColor = "green"
	// Discharging, over 75%
	HighColor = "none"
	// Discharging, between 25 and 75%
	MediumColor = "yellow"
	// Discharging, below 25%
	LowColor = "red"
)

var colorMap = map[string]string{
	"none":   "0",
	"red":    "0;31",
	"green":  "0;32",
	"yellow": "0;33",
	"white":  "0;37",
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
