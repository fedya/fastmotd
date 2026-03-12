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

func RenderNeofetch(logo []string, logoWidth int, lines []InfoLine) string {
	var sb strings.Builder

	var infoStrs []string

	// First line: user@hostname header
	if len(lines) > 0 {
		infoStrs = append(infoStrs, lines[0].Value)
		sepLen := VisibleLen(lines[0].Value)
		infoStrs = append(infoStrs, Colored(BrightBlue, strings.Repeat("─", sepLen)))
		lines = lines[1:]
	}

	// Find max label width for alignment
	maxLabelWidth := 0
	for _, l := range lines {
		if l.Label != "" {
			w := utf8.RuneCountInString(l.Label) + 1 // +1 for ":"
			if w > maxLabelWidth {
				maxLabelWidth = w
			}
		}
	}

	for _, l := range lines {
		if l.Label == "" && l.Value == "" {
			infoStrs = append(infoStrs, "")
		} else if l.Label == "" {
			infoStrs = append(infoStrs, l.Value)
		} else {
			colored := Label(l.Label + ":")
			visLen := utf8.RuneCountInString(l.Label) + 1
			padding := strings.Repeat(" ", max(0, maxLabelWidth-visLen))
			infoStrs = append(infoStrs, colored+padding+"  "+l.Value)
		}
	}

	maxLines := max(len(logo), len(infoStrs))

	gap := "   "

	for i := range maxLines {
		logoLine := ""
		if i < len(logo) {
			logoLine = logo[i]
		}
		pad := max(0, logoWidth-VisibleLen(logoLine))

		infoLine := ""
		if i < len(infoStrs) {
			infoLine = infoStrs[i]
		}

		sb.WriteString(logoLine + strings.Repeat(" ", pad) + gap + infoLine + "\n")
	}

	return sb.String()
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
