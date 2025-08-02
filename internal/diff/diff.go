package diff

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// ShowDiff displays a colored line-by-line diff between two strings.
func ShowDiff(original, modified string) {
	dmp := diffmatchpatch.New()
	charArray1, charArray2, lineArray := dmp.DiffLinesToChars(original, modified)
	diffs := dmp.DiffMain(charArray1, charArray2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	for _, diff := range diffs {
		lines := strings.Split(strings.TrimSuffix(diff.Text, "\n"), "\n")

		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			printLines(lines, "\033[31m-", "\033[0m")
		case diffmatchpatch.DiffInsert:
			printLines(lines, "\033[32m+", "\033[0m")
		case diffmatchpatch.DiffEqual:
			printContextLines(lines)
		}
	}
}

func printLines(lines []string, prefix, suffix string) {
	for _, line := range lines {
		if line != "" {
			fmt.Printf("%s%s%s\n", prefix, line, suffix)
		}
	}
}

func printContextLines(lines []string) {
	if len(lines) <= 6 {
		// Show all lines if there are 6 or fewer
		for _, line := range lines {
			if line != "" {
				fmt.Printf(" %s\n", line)
			}
		}
		return
	}

	// Show first 3 lines
	for i := 0; i < 3 && i < len(lines); i++ {
		if lines[i] != "" {
			fmt.Printf(" %s\n", lines[i])
		}
	}

	fmt.Println(" ...")

	// Show last 3 lines
	for i := len(lines) - 3; i < len(lines); i++ {
		if i >= 0 && lines[i] != "" {
			fmt.Printf(" %s\n", lines[i])
		}
	}
}

