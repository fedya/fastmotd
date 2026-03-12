# fastmotd

Fast system information displayed at login (MOTD) for Linux servers.

## What it shows

- OS, kernel, uptime, hostname, public IP
- Load average, memory, swap, disk usage (with progress bars)
- CPU, GPU, and disk temperatures (NVMe + HDD)
- Docker containers status
- Fail2ban status
- SSL certificate expiry
- Network interfaces traffic
- WiFi signal strength
- Logged in users, last login, process count

## Configuration

Copy the example config to `/etc/fastmotd.toml`:

```bash
sudo cp fastmotd.toml /etc/fastmotd.toml
```

All sections are enabled by default. To disable a section, set it to `false`:

```toml
cert_path = "/path/to/cert.pem"
disk_paths = ["/", "/home", "/var"]

[sections]
docker = false
fail2ban = false
ssl_cert = false
```

If no config file exists, all sections are enabled with default values.

## Requirements

- Go 1.24+
- Linux (reads from `/proc`, `/sys`)

### Disk temperatures

**NVMe** — works out of the box, temperatures read via `/sys/class/hwmon`.

**HDD (SATA)** — requires kernel module `drivetemp`:

```bash
# load once
sudo modprobe drivetemp

# autoload on boot
echo drivetemp | sudo tee /etc/modules-load.d/drivetemp.conf
```

### Docker status

Requires access to the Docker socket. Run `fastmotd` as a user in the `docker` group or as root.

### Fail2ban status

Requires `fail2ban-client` to be available and the user to have permission to run it.

### SSL certificate monitoring

Set `cert_path` in `/etc/fastmotd.toml` to your certificate file.

## Build

```bash
make
```

## Install

```bash
sudo make install
```

This installs:
- `/usr/bin/fastmotd` — the binary
- `/etc/profile.d/fastmotd.sh` — shell hook that runs `fastmotd` at login
- `/etc/fastmotd.toml` — configuration file

## Manual run

```bash
fastmotd
```

## Example output

```
kilogra@gridpoint6
──────────
OS:          ROSA Server
Kernel:      6.12.47-5-x86_64
Uptime:      2d 23h 20m
Hostname:    gridpoint6
Public IP:   192.168.1.120

Load:        ██░░░░░░░░░░░░░ 19%  0.79 / 0.82 / 0.77
Memory:      ███░░░░░░░░░░░░ 22%  6.9 GiB / 30.5 GiB
Swap:        ███░░░░░░░░░░░░ 23%  7.4 GiB / 32.0 GiB
/home:       ████████░░░░░░░ 55%  503.3 GiB / 914.3 GiB
CPU Temp:    56°C
Disk Temp:   Samsung SSD 980 PRO 1TB  45°C
Disk Temp:   ADATA LEGEND 710  29°C
Disk Temp:   TOSHIBA DT01ACA1 (sdb)  33°C
Disk Temp:   Hitachi HDS72101 (sda)  38°C

Docker:      [active]  6 running, 0 exited

Network:     wlan0  ↑ 983.1 GiB  ↓ 94.6 GiB
Users:       kilorga (1 active)
Last login:  kilogra from 172.22.2.68, Thu Mar 12 12:58
Processes:   718
```

## Uninstall

```bash
sudo rm /usr/bin/fastmotd /etc/profile.d/fastmotd.sh /etc/fastmotd.toml
```
