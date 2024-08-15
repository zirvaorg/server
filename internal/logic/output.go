package logic

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func formatText(text string, spaces int) string {
	lines := strings.Split(text, "\n")
	padding := strings.Repeat(" ", spaces)

	for i := 1; i < len(lines); i++ {
		lines[i] = padding + lines[i]
	}
	return strings.Join(lines, "\n")
}

func Output(tp string, text string) {
	var c func(a ...interface{}) string
	var prefix string

	switch tp {
	case "warn":
		c = color.New(color.FgHiYellow).SprintFunc()
		prefix = "[warn]"
	case "error":
		c = color.New(color.FgHiRed).SprintFunc()
		prefix = "[err]"
	case "info":
		c = color.New(color.FgHiGreen).SprintFunc()
		prefix = "[info]"
	case "ok":
		c = color.New(color.FgHiGreen).SprintFunc()
		prefix = "[ok]"
	default:
		prefix = ""
	}

	if prefix != "" {
		fmt.Printf("%s %s\n", c(prefix), formatText(text, len(prefix)+1))
	} else {
		fmt.Printf("%s\n", formatText(text, 0))
	}
}
