package sysinfo

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DockerInfo struct {
	Active  bool
	Running int
	Exited  int
}

type Fail2banInfo struct {
	Active        bool
	TotalBanned   int
	CurrentBanned int
}

type CertInfo struct {
	Expiry time.Time
	Valid  bool
}

func ServiceActive(name string) bool {
	err := exec.Command("systemctl", "is-active", "--quiet", name).Run()
	return err == nil
}

func Docker() *DockerInfo {
	if !ServiceActive("docker.service") {
		return nil
	}
	info := &DockerInfo{Active: true}
	out, err := exec.Command("docker", "ps", "-q").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			info.Running = len(lines)
		}
	}
	out, err = exec.Command("docker", "ps", "-a", "-q", "-f", "status=exited").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			info.Exited = len(lines)
		}
	}
	return info
}

func Fail2ban() *Fail2banInfo {
	if !ServiceActive("fail2ban.service") {
		return nil
	}
	out, err := exec.Command("fail2ban-client", "status", "sshd").Output()
	if err != nil {
		return &Fail2banInfo{Active: true}
	}
	text := string(out)
	info := &Fail2banInfo{Active: true}
	re1 := regexp.MustCompile(`Total banned:\s*(\d+)`)
	if m := re1.FindStringSubmatch(text); len(m) > 1 {
		info.TotalBanned, _ = strconv.Atoi(m[1])
	}
	re2 := regexp.MustCompile(`Currently banned:\s*(\d+)`)
	if m := re2.FindStringSubmatch(text); len(m) > 1 {
		info.CurrentBanned, _ = strconv.Atoi(m[1])
	}
	return info
}

func Certificate(certPath string) *CertInfo {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return nil
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil
	}
	return &CertInfo{
		Expiry: cert.NotAfter,
		Valid:  time.Now().Before(cert.NotAfter),
	}
}
