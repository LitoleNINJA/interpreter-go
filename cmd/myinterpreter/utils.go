package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// checkBracketBalanced checks if the number of opening and closing brackets are equal
func checkBracketBalanced(lines [][]byte) error {
	openingBracket := 0
	closingBracket := 0

	for _, stmt := range lines {
		if strings.Contains(string(stmt), "{") {
			openingBracket++
		}
		if strings.Contains(string(stmt), "}") {
			closingBracket++
		}
	}

	if openingBracket != closingBracket {
		return fmt.Errorf("Error at end: Expect '}'")
	}
	return nil
}

// findBlockEnd finds the end of a block of code enclosed in curly braces
func findBlockEnd() {
	cnt := 1
	lineNumber++
	for lineNumber < len(lines) {
		if bytes.Contains(lines[lineNumber], []byte("{")) {
			cnt++
		} else if bytes.Contains(lines[lineNumber], []byte("}")) {
			cnt--
		}

		if cnt == 0 {
			return
		}

		lineNumber++
	}
	fmt.Fprintf(os.Stderr, "Error at end: Expect '}'")
	os.Exit(69)
}

// mapComplexStmt replaces nested parenthesized expressions with placeholders.
// It converts complex nested expressions like "(a) or (b)" to "%map1 or %map2"
// and maintains a mapping of placeholders to their original expressions.
func mapComplexStmt(stmt []byte) ([]byte, map[int]string) {
	cnt := 1
	start := -1
	level := 0

	simpleStmt := string(stmt)

	mapStmt := map[int]string{}

	for i := 0; i < len(simpleStmt); i++ {
		if simpleStmt[i] == '(' {
			// First opening parenthesis
			if level == 0 {
				start = i
			}
			level++
		} else if simpleStmt[i] == ')' {
			level--
			// Found matching closing parenthesis
			if level == 0 && start != -1 {
				// Extract content between parentheses
				content := simpleStmt[start+1 : i]
				// Store mapping
				mapStmt[cnt] = content

				// Replace with counter value
				replacement := fmt.Sprintf("%%map%d", cnt)
				simpleStmt = simpleStmt[:start] + replacement + simpleStmt[i+1:]

				// Adjust index to account for replacement
				i = start + len(replacement) - 1
				// Increment counter
				cnt++
				start = -1
			}
		}
	}

	return []byte(simpleStmt), mapStmt
}
