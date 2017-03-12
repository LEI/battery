package main

import (
	"fmt"
)

var (
	// DefaultColor to reset foreground and background
	DefaultColor = "default"
	// EmptyColor when empty battery charge
	EmptyColor = "red"
	// FullColor when full battery charge
	FullColor = "green"
	// ChargingColor when plugged in
	ChargingColor = "green"
	// HighColor when discharging, over 75%
	HighColor = "none"
	// MediumColor when discharging, between 25 and 75%
	MediumColor = "yellow"
	// LowColor when discharging, below 25%
	LowColor = "red"
)

var colorMap = map[string]string{
	"default": "0",
	"none":    "0",
	"red":     "0;31",
	"green":   "0;32",
	"yellow":  "0;33",
	"white":   "0;37",
}

// Colors adapter
type Colors interface {
	Get(key string) string
}

type asciiColors struct{}

// Get ASCII colors
func (c *asciiColors) Get(key string) string {
	clr, ok := colorMap[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("\033[%sm", clr)
}

type tmuxColors struct{}

// Get Tmux xolors
func (c *tmuxColors) Get(key string) string {
	return key
}
