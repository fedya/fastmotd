package sysinfo

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type BatteryInfo struct {
	Status   string  // Charging, Discharging, Full, Not charging
	Capacity int     // percent 0-100
	Ratio    float64 // 0.0 - 1.0
}

func Battery() *BatteryInfo {
	// Find first battery in /sys/class/power_supply/
	dirs, err := filepath.Glob("/sys/class/power_supply/BAT*")
	if err != nil || len(dirs) == 0 {
		return nil
	}

	dir := dirs[0]

	capBytes, err := os.ReadFile(filepath.Join(dir, "capacity"))
	if err != nil {
		return nil
	}
	cap, err := strconv.Atoi(strings.TrimSpace(string(capBytes)))
	if err != nil {
		return nil
	}

	status := "Unknown"
	if sb, err := os.ReadFile(filepath.Join(dir, "status")); err == nil {
		status = strings.TrimSpace(string(sb))
	}

	return &BatteryInfo{
		Status:   status,
		Capacity: cap,
		Ratio:    float64(cap) / 100.0,
	}
}
