package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

const DefaultPath = "/etc/fastmotd.toml"

type Config struct {
	CertPath string            `toml:"cert_path"`
	DiskPaths []string         `toml:"disk_paths"`
	Sections map[string]bool   `toml:"sections"`
}

func Load(path string) *Config {
	cfg := &Config{
		CertPath: "",
		DiskPaths: []string{"/home"},
		Sections: map[string]bool{
			"os_info":    true,
			"load":       true,
			"memory":     true,
			"swap":       true,
			"disk_usage": true,
			"cpu_temp":   true,
			"gpu_temp":   true,
			"disk_temp":  true,
			"docker":     true,
			"fail2ban":   true,
			"ssl_cert":   true,
			"network":    true,
			"wifi":       true,
			"users":      true,
			"last_login": true,
			"processes":  true,
		},
	}

	f, err := os.Open(path)
	if err != nil {
		return cfg
	}
	f.Close()

	_, err = toml.DecodeFile(path, cfg)
	if err != nil {
		return cfg
	}

	// Ensure all sections default to true if not specified
	defaults := []string{
		"os_info", "load", "memory", "swap", "disk_usage",
		"cpu_temp", "gpu_temp", "disk_temp", "docker", "fail2ban",
		"ssl_cert", "network", "wifi", "users", "last_login", "processes",
	}
	for _, s := range defaults {
		if _, ok := cfg.Sections[s]; !ok {
			cfg.Sections[s] = true
		}
	}

	return cfg
}

func (c *Config) Section(name string) bool {
	v, ok := c.Sections[name]
	return !ok || v
}
