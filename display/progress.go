package display

import (
	"fmt"
	"strings"
)

const BarWidth = 15

// ProgressBar generates a bar like ████████░░░░░░░ 53%
func ProgressBar(ratio float64, width int) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	filled := int(ratio * float64(width))
	empty := width - filled
	var sb strings.Builder
	for range filled {
		sb.WriteString("█")
	}
	for range empty {
		sb.WriteString("░")
	}
	color := ColorByRatio(ratio)
	return fmt.Sprintf("%s%s%s %s%d%%%s", color, sb.String(), Reset, color, int(ratio*100), Reset)
}
