package parser

import (
	"fmt"
	"strings"
)

type ANSIColor string

const (
	ANSIRed    ANSIColor = "\033[31m"
	ANSIGreen  ANSIColor = "\033[32m"
	ANSIYellow ANSIColor = "\033[33m"
	ANSICyan   ANSIColor = "\033[36m"
	ANSIWhite  ANSIColor = "\033[37m"
	ANSIReset  ANSIColor = "\033[0m"
)

type Severity string

const (
	SeverityInfo  Severity = "info"
	SeverityDebug Severity = "debug"
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

func parseSeverity(line string) Severity {
	upper := strings.ToUpper(line)

	switch {
	case strings.Contains(upper, "ERROR"):
		return SeverityError

	case strings.Contains(upper, "WARN"):
		return SeverityWarn

	case strings.Contains(upper, "INFO"):
		return SeverityInfo

	case strings.Contains(upper, "DEBUG"):
		return SeverityDebug

	default:
		return SeverityInfo
	}
}

func colorForSeverity(severity Severity) ANSIColor {
	switch severity {
	case SeverityError:
		return ANSIRed
	case SeverityWarn:
		return ANSIYellow
	case SeverityInfo:
		return ANSIGreen
	case SeverityDebug:
		return ANSICyan
	default:
		return ANSIWhite
	}
}

func colorPrintf(color ANSIColor, text string) {
	fmt.Printf("%s%s%s", color, text, ANSIReset)
}

func PrintLog(line string, args ...any) {
	text := fmt.Sprintf(line, args...)
	severity := parseSeverity(text)
	color := colorForSeverity(severity)
	colorPrintf(color, text)
}

func colorSprintf(color ANSIColor, format string, args ...any) string {
	return fmt.Sprintf("%s%s%s", color, fmt.Sprintf(format, args...), ANSIReset)
}

func GenerateLog(line string, args ...any) string {
	text := fmt.Sprintf(line, args...)
	severity := parseSeverity(text)
	color := colorForSeverity(severity)
	return colorSprintf(color, text)
}
