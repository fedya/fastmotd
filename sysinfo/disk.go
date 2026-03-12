package sysinfo

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type DiskInfo struct {
	Total      uint64
	Used       uint64
	Available  uint64
	Ratio      float64
	MountPoint string
}

func DiskUsage(path string) DiskInfo {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return DiskInfo{MountPoint: path}
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)
	used := total - free
	ratio := 0.0
	if total > 0 {
		ratio = float64(used) / float64(total)
	}
	return DiskInfo{
		Total:      total,
		Used:       used,
		Available:  available,
		Ratio:      ratio,
		MountPoint: path,
	}
}

type DiskTemp struct {
	Device string
	Temp   int // degrees C
}

func DiskTemps() []DiskTemp {
	hwmonDirs, err := filepath.Glob("/sys/class/hwmon/hwmon*")
	if err != nil {
		return nil
	}
	var temps []DiskTemp
	for _, dir := range hwmonDirs {
		nameBytes, err := os.ReadFile(filepath.Join(dir, "name"))
		if err != nil {
			continue
		}
		name := strings.TrimSpace(string(nameBytes))
		if name != "drivetemp" && name != "nvme" {
			continue
		}
		tempBytes, err := os.ReadFile(filepath.Join(dir, "temp1_input"))
		if err != nil {
			continue
		}
		millideg, err := strconv.Atoi(strings.TrimSpace(string(tempBytes)))
		if err != nil {
			continue
		}
		dev := diskDevName(dir, name)
		temps = append(temps, DiskTemp{Device: dev, Temp: millideg / 1000})
	}
	return temps
}

func diskDevName(hwmonDir, driverName string) string {
	devicePath, err := filepath.EvalSymlinks(filepath.Join(hwmonDir, "device"))
	if err != nil {
		return "?"
	}

	if driverName == "nvme" {
		// try to read model name
		if model, err := os.ReadFile(filepath.Join(devicePath, "model")); err == nil {
			return strings.TrimSpace(string(model))
		}
		return filepath.Base(devicePath)
	}

	// drivetemp: use model name + block device
	model := ""
	if m, err := os.ReadFile(filepath.Join(devicePath, "model")); err == nil {
		model = strings.TrimSpace(string(m))
	}
	blocks, err := filepath.Glob(filepath.Join(devicePath, "block", "*"))
	if err != nil || len(blocks) == 0 {
		if model != "" {
			return model
		}
		return fmt.Sprintf("disk(%s)", filepath.Base(devicePath))
	}
	blockName := filepath.Base(blocks[0])
	if model != "" {
		return fmt.Sprintf("%s (%s)", model, blockName)
	}
	return blockName
}
