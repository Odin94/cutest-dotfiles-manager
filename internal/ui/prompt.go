package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ConfirmTraverseUp() (bool, error) {
	fmt.Print("No .cdm.toml here. Traverse up to find one (up to 5 levels)? [y/N]: ")
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		return false, sc.Err()
	}
	line := strings.TrimSpace(strings.ToLower(sc.Text()))
	return line == "y" || line == "yes", nil
}
