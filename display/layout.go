package display

import (
	"strings"
	"unicode/utf8"
)

type InfoLine struct {
	Label string
	Value string
}

// StripANSI removes ANSI escape codes to measure visible string length
func StripANSI(s string) string {
	var result []byte
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		if b[i] == '\033' {
			for i < len(b) && b[i] != 'm' {
				i++
			}
			continue
		}
		result = append(result, b[i])
	}
	return string(result)
}

func VisibleLen(s string) int {
	return utf8.RuneCountInString(StripANSI(s))
}

func RenderLines(lines []InfoLine) string {
	var sb strings.Builder

	// First line: user@hostname header
	if len(lines) > 0 {
		sb.WriteString(lines[0].Value + "\n")
		sepLen := VisibleLen(lines[0].Value)
		sb.WriteString(Colored(BrightBlue, strings.Repeat("─", sepLen)) + "\n")
		lines = lines[1:]
	}

	// Find max label width for alignment
	maxLabelWidth := 0
	for _, l := range lines {
		if l.Label != "" {
			w := utf8.RuneCountInString(l.Label) + 1
			if w > maxLabelWidth {
				maxLabelWidth = w
			}
		}
	}

	for _, l := range lines {
		if l.Label == "" && l.Value == "" {
			sb.WriteString("\n")
		} else if l.Label == "" {
			sb.WriteString(l.Value + "\n")
		} else {
			colored := Label(l.Label + ":")
			visLen := utf8.RuneCountInString(l.Label) + 1
			padding := strings.Repeat(" ", max(0, maxLabelWidth-visLen))
			sb.WriteString(colored + padding + "  " + l.Value + "\n")
		}
	}

	return sb.String()
}
