package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fedya/fastmotd/config"
	"github.com/fedya/fastmotd/display"
	"github.com/fedya/fastmotd/sysinfo"
)

func main() {
	cfg := config.Load(config.DefaultPath)

	var lines []display.InfoLine

	// Header: user@hostname
	users := sysinfo.Users()
	hostname := sysinfo.Hostname()
	header := fmt.Sprintf("%s%s%s@%s%s%s",
		display.BrightCyan, users.Current, display.Reset,
		display.BrightCyan, hostname, display.Reset)
	lines = append(lines, display.InfoLine{Value: header})

	// System info
	if cfg.Section("os_info") {
		lines = append(lines, display.InfoLine{Label: "OS", Value: sysinfo.OSRelease()})
		lines = append(lines, display.InfoLine{Label: "Kernel", Value: sysinfo.Kernel()})
		lines = append(lines, display.InfoLine{Label: "Uptime", Value: sysinfo.Uptime()})
		lines = append(lines, display.InfoLine{Label: "Hostname", Value: hostname})
		lines = append(lines, display.InfoLine{Label: "Public IP", Value: sysinfo.PublicIP()})
		lines = append(lines, display.InfoLine{}) // separator
	}

	// Resources with progress bars
	if cfg.Section("load") {
		load := sysinfo.LoadAvg()
		loadRatio := load.Load1 / 4.0
		if loadRatio > 1 {
			loadRatio = 1
		}
		lines = append(lines, display.InfoLine{
			Label: "Load",
			Value: fmt.Sprintf("%s %.2f / %.2f / %.2f",
				display.ProgressBar(loadRatio, display.BarWidth),
				load.Load1, load.Load5, load.Load15),
		})
	}

	if cfg.Section("memory") {
		mem := sysinfo.Memory()
		lines = append(lines, display.InfoLine{
			Label: "Memory",
			Value: fmt.Sprintf("%s %s / %s",
				display.ProgressBar(mem.Ratio, display.BarWidth),
				sysinfo.HumanizeBytes(mem.Used), sysinfo.HumanizeBytes(mem.Total)),
		})
	}

	if cfg.Section("swap") {
		swap := sysinfo.Swap()
		lines = append(lines, display.InfoLine{
			Label: "Swap",
			Value: fmt.Sprintf("%s %s / %s",
				display.ProgressBar(swap.Ratio, display.BarWidth),
				sysinfo.HumanizeBytes(swap.Used), sysinfo.HumanizeBytes(swap.Total)),
		})
	}

	if cfg.Section("disk_usage") {
		for _, dp := range cfg.DiskPaths {
			disk := sysinfo.DiskUsage(dp)
			lines = append(lines, display.InfoLine{
				Label: dp,
				Value: fmt.Sprintf("%s %s / %s",
					display.ProgressBar(disk.Ratio, display.BarWidth),
					sysinfo.HumanizeBytes(disk.Used), sysinfo.HumanizeBytes(disk.Total)),
			})
		}
	}

	if cfg.Section("cpu_temp") {
		if temp, ok := sysinfo.CPUTemp(); ok {
			lines = append(lines, display.InfoLine{Label: "CPU Temp", Value: display.ColoredTemp(temp)})
		}
	}

	if cfg.Section("gpu_temp") {
		if gpu, ok := sysinfo.GPUTemp(); ok {
			lines = append(lines, display.InfoLine{
				Label: "GPU Temp",
				Value: fmt.Sprintf("%s [%s]", display.ColoredTemp(gpu.Temp), gpu.Driver),
			})
		}
	}

	if cfg.Section("disk_temp") {
		diskTemps := sysinfo.DiskTemps()
		for _, dt := range diskTemps {
			lines = append(lines, display.InfoLine{
				Label: "Disk Temp",
				Value: fmt.Sprintf("%s [%s]", display.ColoredTemp(dt.Temp), dt.Device),
			})
		}
	}

	lines = append(lines, display.InfoLine{}) // separator

	// Services
	if cfg.Section("docker") {
		docker := sysinfo.Docker()
		if docker != nil {
			lines = append(lines, display.InfoLine{
				Label: "Docker",
				Value: fmt.Sprintf("%s %d running, %d exited",
					display.StatusActive("active"), docker.Running, docker.Exited),
			})
		}
	}

	if cfg.Section("fail2ban") {
		f2b := sysinfo.Fail2ban()
		if f2b != nil {
			lines = append(lines, display.InfoLine{
				Label: "Fail2ban",
				Value: fmt.Sprintf("%s %d total banned, %d current",
					display.StatusActive("active"), f2b.TotalBanned, f2b.CurrentBanned),
			})
		}
	}

	if cfg.Section("ssl_cert") && cfg.CertPath != "" {
		cert := sysinfo.Certificate(cfg.CertPath)
		if cert != nil {
			certColor := display.BrightGreen
			status := "valid"
			if !cert.Valid {
				certColor = display.BrightRed
				status = "EXPIRED"
			} else if time.Until(cert.Expiry) < 30*24*time.Hour {
				certColor = display.BrightYellow
				status = "expiring soon"
			}
			lines = append(lines, display.InfoLine{
				Label: "SSL Cert",
				Value: fmt.Sprintf("%s[%s]%s  expires %s",
					certColor, status, display.Reset,
					cert.Expiry.Format("Jan 02 2006")),
			})
		}
	}

	lines = append(lines, display.InfoLine{}) // separator

	// Network
	if cfg.Section("network") {
		nets := sysinfo.NetworkStats()
		for _, n := range nets {
			lines = append(lines, display.InfoLine{
				Label: "Network",
				Value: fmt.Sprintf("%s ↑ %s ↓ %s",
					n.Interface,
					sysinfo.HumanizeBytes(n.TXBytes),
					sysinfo.HumanizeBytes(n.RXBytes)),
			})
		}
	}

	if cfg.Section("wifi") {
		wifiList := sysinfo.WifiSignal()
		for _, w := range wifiList {
			lines = append(lines, display.InfoLine{
				Label: "WiFi",
				Value: fmt.Sprintf("%s %d dBm (quality %d/70)", w.Interface, w.Signal, w.Quality),
			})
		}
	}

	// Users
	if cfg.Section("users") {
		if len(users.ActiveUsers) > 0 {
			lines = append(lines, display.InfoLine{
				Label: "Users",
				Value: fmt.Sprintf("%s (%d active)", strings.Join(users.ActiveUsers, " "), users.ActiveCount),
			})
		}
	}

	if cfg.Section("last_login") {
		lastLogin := sysinfo.LastLogin()
		if lastLogin != nil {
			lines = append(lines, display.InfoLine{
				Label: "Last login",
				Value: fmt.Sprintf("%s from %s, %s", lastLogin.User, lastLogin.From, lastLogin.When),
			})
		}
	}

	if cfg.Section("processes") {
		lines = append(lines, display.InfoLine{
			Label: "Processes",
			Value: fmt.Sprintf("%d", sysinfo.ProcessCount()),
		})
	}

	if cfg.Section("battery") {
		bat := sysinfo.Battery()
		if bat != nil {
			// Invert ratio so low battery = red (like high usage)
			barRatio := 1.0 - bat.Ratio
			statusColor := display.ColorByRatio(barRatio)
			lines = append(lines, display.InfoLine{
				Label: "Battery",
				Value: fmt.Sprintf("%s %s[%s]%s",
					display.ProgressBarWithColor(bat.Ratio, display.BarWidth, display.ColorByRatio(barRatio)),
					statusColor, bat.Status, display.Reset),
			})
		}
	}



	// Render
	fmt.Print("\n")
	fmt.Print(display.RenderLines(lines))
	fmt.Print("\n")
}
