package main

import (
	"bytes"
	"fmt"
	"strconv"
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
func findBlockEnd() (int, error) {
	cnt := 1
	blockEndLineNumber := lineNumber + 1

	for blockEndLineNumber < len(lines) {
		if bytes.Contains(lines[blockEndLineNumber], []byte("{")) {
			cnt++
		} else if bytes.Contains(lines[blockEndLineNumber], []byte("}")) {
			cnt--
		}

		if cnt == 0 {
			return blockEndLineNumber, nil
		}

		blockEndLineNumber++
	}

	return -1, fmt.Errorf("Error at end: Expect '}'")
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

// isTruthy checks if a value is truthy
func isTruthy(val Value) bool {
	switch val := val.(type) {
	case bool:
		return val
	case string:
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f != 0
		}
		if val == "nil" {
			return false
		}

		return true
	case float64:
		return val != 0
	default:
		fmt.Println("Unknown type !")
		return false
	}
}

// evaluateCondition evaluates a condition and returns the result
func evaluateCondition(condition []byte) (Value, error) {
	var conditionResult Value
	if isAssignment(condition) {
		handleAssignment(string(condition))
		conditionResult = true
	} else {
		var err error
		conditionResult, err = evaluate(condition)
		if err != nil {
			return nil, err
		}
	}

	return conditionResult, nil
}

func extractCondition(stmt []byte, prefix string) ([]byte, []byte, error) {
	// fmt.Printf("Stmt : %s\n", stmt)
	s, ok := strings.CutPrefix(string(stmt), prefix)
	if !ok {
		return []byte{}, []byte{}, fmt.Errorf("stmt should start with %s: %s", prefix, stmt)
	}

	// Find closing parenthesis
	closeParenIndex := strings.Index(s, ")")
	if closeParenIndex == -1 {
		return []byte{}, []byte{}, fmt.Errorf(") not found in %s stmt : %s", prefix, stmt)
	}

	condition := strings.TrimSpace(s[1:closeParenIndex])
	body := strings.TrimSpace(s[closeParenIndex+1:])
	if strings.HasPrefix(body, "{") {
		restOfBody := strings.TrimSpace(strings.TrimPrefix(body, "{"))
		body = "{"
		lines = append(lines[:lineNumber+1], append([][]byte{[]byte(restOfBody)}, lines[lineNumber+1:]...)...)
	}

	return []byte(condition), []byte(body), nil
}

func parseForStmt(stmt []byte) ([]byte, []byte, []byte, error) {
	parts := strings.Split(string(stmt), ";")

	if len(parts) != 3 {
		return []byte{}, []byte{}, []byte{}, fmt.Errorf("Invalid for stmt : %s", stmt)
	}

	// trim space for all parts
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return []byte(parts[0]), []byte(parts[1]), []byte(parts[2]), nil
}