package sysinfo

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type LoadInfo struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

func LoadAvg() LoadInfo {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return LoadInfo{}
	}
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return LoadInfo{}
	}
	l1, _ := strconv.ParseFloat(fields[0], 64)
	l5, _ := strconv.ParseFloat(fields[1], 64)
	l15, _ := strconv.ParseFloat(fields[2], 64)
	return LoadInfo{Load1: l1, Load5: l5, Load15: l15}
}

func CPUTemp() (int, bool) {
	hwmonDirs, err := filepath.Glob("/sys/class/hwmon/hwmon*")
	if err != nil {
		return 0, false
	}
	for _, dir := range hwmonDirs {
		nameBytes, err := os.ReadFile(filepath.Join(dir, "name"))
		if err != nil {
			continue
		}
		name := strings.TrimSpace(string(nameBytes))
		if name == "k10temp" || name == "coretemp" {
			tempBytes, err := os.ReadFile(filepath.Join(dir, "temp1_input"))
			if err != nil {
				continue
			}
			millideg, err := strconv.Atoi(strings.TrimSpace(string(tempBytes)))
			if err != nil {
				continue
			}
			return millideg / 1000, true
		}
	}
	return 0, false
}

type GPUTempInfo struct {
	Temp   int
	Driver string
}

func GPUTemp() (*GPUTempInfo, bool) {
	hwmonDirs, err := filepath.Glob("/sys/class/hwmon/hwmon*")
	if err != nil {
		return nil, false
	}
	for _, dir := range hwmonDirs {
		nameBytes, err := os.ReadFile(filepath.Join(dir, "name"))
		if err != nil {
			continue
		}
		name := strings.TrimSpace(string(nameBytes))
		if name == "amdgpu" || name == "nouveau" || name == "nvidia" {
			tempBytes, err := os.ReadFile(filepath.Join(dir, "temp1_input"))
			if err != nil {
				continue
			}
			millideg, err := strconv.Atoi(strings.TrimSpace(string(tempBytes)))
			if err != nil {
				continue
			}
			return &GPUTempInfo{Temp: millideg / 1000, Driver: name}, true
		}
	}
	return nil, false
}

func ProcessCount() int {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			if _, err := strconv.Atoi(e.Name()); err == nil {
				count++
			}
		}
	}
	return count
}
