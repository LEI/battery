// +android

package main

import (
// "fmt"
// "os"
// "os/exec"
)

// Health: GOOD
// Plugged: UNPLUGGED, PLUGGED_AC, PLUGGED_USB
// Status: DISCHARGING, CHARGING, FULL

var termuxBatteryStatus = []byte(`{
	 "health": "GOOD",
	 "percentage": 100,
	 "plugged": "UNPLUGGED",
	 "status": "DISCHARGING",
	 "temperature": 13.37
}`)

// func init() {
// 	fmt.Println("termux-battery-status:", string(termuxBatteryStatus))
// }
