package sysinfo

import (
	"os/exec"
	"os/user"
	"strings"
)

type UserInfo struct {
	Current     string
	ActiveUsers []string
	ActiveCount int
}

type LastLoginInfo struct {
	User string
	From string
	When string
}

func Users() UserInfo {
	u, err := user.Current()
	current := "unknown"
	if err == nil {
		current = u.Username
	}
	out, err := exec.Command("who").Output()
	if err != nil {
		return UserInfo{Current: current}
	}
	seen := make(map[string]bool)
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] != "" {
			seen[fields[0]] = true
		}
	}
	var users []string
	for u := range seen {
		users = append(users, u)
	}
	return UserInfo{Current: current, ActiveUsers: users, ActiveCount: len(users)}
}

func LastLogin() *LastLoginInfo {
	out, err := exec.Command("last", "-n", "1", "-w").Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 1 {
		return nil
	}
	fields := strings.Fields(lines[0])
	if len(fields) < 4 || fields[0] == "wtmp" || fields[0] == "reboot" {
		return nil
	}
	info := &LastLoginInfo{User: fields[0]}
	if len(fields) >= 3 {
		info.From = fields[2]
	}
	if len(fields) >= 7 {
		info.When = strings.Join(fields[3:7], " ")
	}
	return info
}
