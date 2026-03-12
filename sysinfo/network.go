package sysinfo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type NetInfo struct {
	Interface string
	RXBytes   uint64
	TXBytes   uint64
}

func NetworkStats() []NetInfo {
	entries, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return nil
	}
	var results []NetInfo
	skip := []string{"lo", "veth", "br-", "docker", "virbr"}
	for _, e := range entries {
		name := e.Name()
		shouldSkip := false
		for _, prefix := range skip {
			if strings.HasPrefix(name, prefix) {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			continue
		}
		base := filepath.Join("/sys/class/net", name, "statistics")
		rxData, err := os.ReadFile(filepath.Join(base, "rx_bytes"))
		if err != nil {
			continue
		}
		txData, err := os.ReadFile(filepath.Join(base, "tx_bytes"))
		if err != nil {
			continue
		}
		rx, _ := strconv.ParseUint(strings.TrimSpace(string(rxData)), 10, 64)
		tx, _ := strconv.ParseUint(strings.TrimSpace(string(txData)), 10, 64)
		if rx == 0 && tx == 0 {
			continue
		}
		results = append(results, NetInfo{Interface: name, RXBytes: rx, TXBytes: tx})
	}
	return results
}

type WifiInfo struct {
	Interface string
	Signal    int // dBm
	Quality   int // link quality (0-70)
}

func WifiSignal() []WifiInfo {
	f, err := os.Open("/proc/net/wireless")
	if err != nil {
		return nil
	}
	defer f.Close()

	var results []WifiInfo
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// skip header lines
		if strings.Contains(line, "|") || strings.TrimSpace(line) == "" {
			continue
		}
		// format: " wlan0: 0000   66.  -44.  -256 ..."
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		iface := strings.TrimSuffix(fields[0], ":")
		quality, _ := strconv.ParseFloat(strings.TrimSuffix(fields[2], "."), 64)
		signal, _ := strconv.ParseFloat(strings.TrimSuffix(fields[3], "."), 64)
		results = append(results, WifiInfo{
			Interface: iface,
			Quality:   int(quality),
			Signal:    int(signal),
		})
	}
	return results
}

func WifiBar(quality int) string {
	// quality is 0-70 from /proc/net/wireless
	pct := float64(quality) / 70.0
	if pct > 1 {
		pct = 1
	}
	bars := int(pct * 5)
	full := "▂▄▆█"
	empty := "▂"
	var sb strings.Builder
	for i := 0; i < 4; i++ {
		if i < bars {
			sb.WriteRune(rune(full[i]))
		} else {
			sb.WriteString(empty)
		}
	}
	return fmt.Sprintf("%s (%d dBm)", sb.String(), quality)
}
