package main

import (
	"fmt"
	"strings"
	"unicode"
)

var lines [][]byte
var lineNumber int
var isString bool
var currentScope *Scope

func readLines(fileContent []byte) [][]byte {
	var lines [][]byte
	line := make([]byte, 0)
	for i := 0; i < len(fileContent); i++ {
		ch := fileContent[i]

		if ch == '"' {
			isString = !isString
			line = append(line, ch)
			continue
		}
		if isString {
			line = append(line, ch)
			continue
		}

		if ch == ';' {
			line = append(line, ch)
			trimmed := []byte(strings.TrimSpace(string(line)))
			if len(trimmed) > 0 {
				// fmt.Printf("Line : %s, Len %d\n", trimmed, len(trimmed))
				lines = append(lines, trimmed)
			}
			line = []byte{}
			continue
		}

		if ch == '\n' {
			// Only add newline if inside a string or if current statement is empty
			if len(line) > 0 && isString {
				line = append(line, ' ')
			} else if len(line) > 0 {
				trimmed := []byte(strings.TrimSpace(string(line)))
				if len(trimmed) > 0 {
					// fmt.Printf("Line : %s, Len %d\n", trimmed, len(trimmed))
					lines = append(lines, trimmed)
				}
				line = []byte{}
				continue
			}
			continue
		}

		line = append(line, ch)
	}

	// Handle last statement if it exists
	if len(line) > 0 {
		trimmed := []byte(strings.TrimSpace(string(line)))
		if len(trimmed) > 0 {
			// fmt.Printf("Line : %s, Len %d\n", trimmed, len(trimmed))
			lines = append(lines, trimmed)
		}
	}
	return lines
}

func getPrintContents(line []byte) []byte {
	s, ok := strings.CutPrefix(string(line), "print")
	if !ok {
		fmt.Printf("Print line dosent start with Print : %s\n", line)
		return []byte{}
	}
	s, _ = strings.CutSuffix(string(s), ";")

	s = strings.TrimSpace(s)
	if val, ok := currentScope.getScopeValue(s); ok {
		s = val
	}
	return []byte(s)
}

func isPrintStmt(stmt []byte) bool {
	stmtString := string(stmt)
	return strings.HasPrefix(stmtString, "print")
}

func isVarDeclaration(stmt []byte) bool {
	stmtString := string(stmt)
	return strings.HasPrefix(stmtString, "var ")
}

func getVarDeclaration(stmt []byte) error {
	stmtString := string(stmt)
	stmtString, _ = strings.CutPrefix(stmtString, "var ")
	stmtString, _ = strings.CutSuffix(stmtString, ";")

	pos := strings.Index(stmtString, "=")
	val := "nil"
	if pos == -1 {
		pos = len(stmtString)
	} else {
		finalVal, err := handleAssignment(stmtString[pos+1:])
		if err != nil {
			return err
		}

		val = finalVal
	}
	key := strings.TrimSpace(stmtString[:pos])

	// fmt.Printf("Key : %s, Value : %s\n", key, val)
	currentScope.setScopeValue(key, val)

	return nil
}

func isBlockStart(stmt []byte) bool {
	stmtString := strings.TrimSpace(string(stmt))

	return stmtString == "{"
}

func isBlockEnd(stmt []byte) bool {
	stmtString := strings.TrimSpace(string(stmt))

	return stmtString == "}"
}

func checkBracketBalanced(lines [][]byte) error {
	openingBracket := 0
	closingBracket := 0

	for _, stmt := range lines {
		if isBlockStart(stmt) {
			openingBracket++
		} else if isBlockEnd(stmt) {
			closingBracket++
		}
	}

	if openingBracket != closingBracket {
		return fmt.Errorf("Error at end: Expect '}'")
	}
	return nil
}

func run(fileContents []byte) error {
	lineNumber = 0
	lines = readLines(fileContents)
	currentScope = NewScope(nil)
	// fmt.Println(lines)

	if err := checkBracketBalanced(lines); err != nil {
		exitCode = 65
		return err
	}

	for {
		if lineNumber >= len(lines) {
			break
		}

		stmt := lines[lineNumber]

		err := handleStmt(stmt)
		if err != nil {
			return err
		}

		lineNumber++
	}
	return nil
}

func handleStmt(stmt []byte) error {
	printStmt := false

	if isPrintStmt(stmt) {
		printStmt = true
		stmt = getPrintContents(stmt)
	} else if isVarDeclaration(stmt) {
		err := getVarDeclaration(stmt)
		if err != nil {
			return err
		}

		return nil
	} else if isBlockStart(stmt) {
		err := handleBlock()
		if err != nil {
			exitCode = 65
			return err
		}

		return nil
	}

	if strings.Contains(string(stmt), "=") {
		handleAssignment(string(stmt))
	}

	if len(stmt) == 0 {
		if printStmt {
			exitCode = 65
			return fmt.Errorf("empty print stmt")
		} else {
			return nil
		}
	}

	// fmt.Printf("Eval : %s, Len : %d\n", stmt, len(stmt))
	expr, err := evaluate(stmt)
	if err != nil {
		return err
	}

	if printStmt {
		fmt.Println(expr)
	}

	return nil
}

func handleAssignment(stmt string) (string, error) {
	stmt, _ = strings.CutSuffix(stmt, ";")

	if strings.Contains(stmt, "=") {
		pos := strings.Index(stmt, "=")

		key := strings.TrimSpace(stmt[:pos])
		val, err := handleAssignment(stmt[pos+1:])
		if err != nil {
			return val, err
		}

		// fmt.Printf("Key : %s, Value : %s\n", key, val)
		// Try to find and update existing variable
		if success := currentScope.assignScopeValue(key, val); !success {
			// If variable doesn't exist anywhere, define it in current scope
			// fmt.Printf("Not found : %s, %s\n", key, val)
			currentScope.setScopeValue(key, val)
		}
		return val, nil
	} else {
		val := strings.TrimSpace(stmt)
		if strings.ContainsAny(val, "+-*/()") {
			evalVal, err := evaluate([]byte(val))
			if err != nil {
				return val, err
			}

			if strings.HasPrefix(val, `"`) {
				val = `"` + fmt.Sprint(evalVal) + `"`
			} else {
				val = fmt.Sprint(evalVal)
			}
		} else if mapVal, ok := currentScope.getScopeValue(val); ok {
			val = mapVal
		} else if unicode.IsLetter(rune(val[0])) && val != "true" && val != "false" {
			exitCode = 70
			return val, fmt.Errorf("Undefined variable '%s'", val)
		}

		return val, nil
	}
}

func handleBlock() error {
	// push new scope
	enclosingScope := currentScope
	currentScope = NewScope(enclosingScope)

	for {
		lineNumber++
		if lineNumber >= len(lines) {
			break
		}

		stmt := lines[lineNumber]

		if isBlockEnd(stmt) {
			// pop scope
			currentScope = enclosingScope
			return nil
		}

		err := handleStmt(stmt)
		if err != nil {
			return err
		}
	}

	return fmt.Errorf("Error at end: Expect '}' .")
}
