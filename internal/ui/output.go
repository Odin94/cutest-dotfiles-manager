package ui

import (
	"fmt"
)

const (
	AnsiGreen  = "\033[32;1m"
	AnsiRed    = "\033[31;1m"
	AnsiYellow = "\033[33;1m"
	AnsiDim    = "\033[90m"
	AnsiReset  = "\033[0m"
)

const (
	ansiGreen  = AnsiGreen
	ansiRed    = AnsiRed
	ansiYellow = AnsiYellow
	ansiDim    = AnsiDim
	ansiReset  = AnsiReset
)

func PrintSuccess(msg string) {
	fmt.Println(ansiGreen + msg + ansiReset)
}

func PrintError(msg string) {
	fmt.Println(ansiRed + msg + ansiReset)
}

func PrintWarn(msg string) {
	fmt.Println(ansiYellow + msg + ansiReset)
}

func PrintDryRunBanner() {
	fmt.Println(ansiYellow + "--- dry run (no files or state will be written) ---" + ansiReset)
}

func PrintMissingVarsWarning(missing []string, affected map[string]string) {
	fmt.Println(ansiYellow + "Missing variables (mappings using them will be skipped):" + ansiReset)
	for _, v := range missing {
		fmt.Println(ansiYellow + "  - " + v + ansiReset)
	}
	if len(affected) > 0 {
		fmt.Println(ansiYellow + "Affected mappings:" + ansiReset)
		for src, dest := range affected {
			fmt.Println(ansiDim + "    " + src + " -> " + dest + ansiReset)
		}
	}
}
