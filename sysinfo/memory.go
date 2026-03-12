package sysinfo

import (
	"os"
	"strconv"
	"strings"
)

type MemInfo struct {
	Total     uint64
	Used      uint64
	Available uint64
	Ratio     float64
}

type SwapInfo struct {
	Total uint64
	Used  uint64
	Ratio float64
}

func parseMemValue(line string) uint64 {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return 0
	}
	val, _ := strconv.ParseUint(parts[1], 10, 64)
	return val * 1024 // convert from kB to bytes
}

func Memory() MemInfo {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return MemInfo{}
	}
	var total, available uint64
	for line := range strings.SplitSeq(string(data), "\n") {
		switch {
		case strings.HasPrefix(line, "MemTotal:"):
			total = parseMemValue(line)
		case strings.HasPrefix(line, "MemAvailable:"):
			available = parseMemValue(line)
		}
	}
	used := total - available
	ratio := 0.0
	if total > 0 {
		ratio = float64(used) / float64(total)
	}
	return MemInfo{Total: total, Used: used, Available: available, Ratio: ratio}
}

func Swap() SwapInfo {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return SwapInfo{}
	}
	var swapTotal, swapFree uint64
	for line := range strings.SplitSeq(string(data), "\n") {
		switch {
		case strings.HasPrefix(line, "SwapTotal:"):
			swapTotal = parseMemValue(line)
		case strings.HasPrefix(line, "SwapFree:"):
			swapFree = parseMemValue(line)
		}
	}
	used := swapTotal - swapFree
	ratio := 0.0
	if swapTotal > 0 {
		ratio = float64(used) / float64(swapTotal)
	}
	return SwapInfo{Total: swapTotal, Used: used, Ratio: ratio}
}
