package display

import "fmt"

const (
	Reset        = "\033[0m"
	Bold         = "\033[1m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Cyan         = "\033[36m"
	BrightRed    = "\033[91m"
	BrightGreen  = "\033[92m"
	BrightYellow = "\033[93m"
	BrightBlue   = "\033[94m"
	BrightCyan   = "\033[96m"
	BrightWhite  = "\033[97m"
)

func Colored(color, s string) string {
	return color + s + Reset
}

func Label(s string) string {
	return Colored(BrightCyan, s)
}

func Value(s string) string {
	return Colored(BrightWhite, s)
}

func ColorByRatio(ratio float64) string {
	switch {
	case ratio < 0.4:
		return BrightGreen
	case ratio < 0.8:
		return BrightYellow
	default:
		return BrightRed
	}
}

func ColoredRatio(ratio float64, s string) string {
	return Colored(ColorByRatio(ratio), s)
}

func ColorByTemp(temp int) string {
	switch {
	case temp < 50:
		return BrightGreen
	case temp < 70:
		return BrightYellow
	default:
		return BrightRed
	}
}

func ColoredTemp(temp int) string {
	return fmt.Sprintf("%s%d°C%s", ColorByTemp(temp), temp, Reset)
}

func StatusActive(s string) string {
	return Colored(BrightGreen, "["+s+"]")
}

func StatusInactive(s string) string {
	return Colored(BrightRed, "["+s+"]")
}
