package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PrintOSln(out string) {
	fmt.Fprintln(os.Stdout, out)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Input(reader *bufio.Reader) []string {
	fmt.Fprint(os.Stdout, "\033[1m\033[36m$ ")
	fmt.Fprint(os.Stdout, "\033[1m\033[33m")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: ", err)
		return []string{}
	}

	// Reset the formatting
	fmt.Fprint(os.Stdout, "\033[0m")

	return strings.Split(strings.TrimSpace(input), " ")
}
