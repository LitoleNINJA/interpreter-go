package main

import (
	"fmt"
	"strings"
	"unicode"
)

var values map[string]string
var lines [][]byte
var lineNumber int

func readLines(fileContent []byte) [][]byte {
	var lines [][]byte
	line := make([]byte, 0)
	for i := 0; i <= len(fileContent); i++ {
		if i == len(fileContent) || fileContent[i] == 10 {
			line = []byte(strings.TrimSpace(string(line)))
			// fmt.Printf("Line : %s\n", line)
			lines = append(lines, line)
			line = make([]byte, 0)
		} else {
			line = append(line, fileContent[i])
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
	if val, ok := values[s]; ok {
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
	values[key] = val

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
	values = make(map[string]string)
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
	if strings.Contains(stmt, "=") {
		pos := strings.Index(stmt, "=")

		key := strings.TrimSpace(stmt[:pos])
		val, err := handleAssignment(stmt[pos+1:])
		if err != nil {
			return val, err
		}

		// fmt.Printf("Key : %s, Value : %s\n", key, val)
		values[key] = val
		return val, nil
	} else {
		val := strings.TrimSpace(stmt)
		if strings.ContainsAny(val, "+-*/()") {
			evalVal, err := evaluate([]byte(val))
			if err != nil {
				return val, err
			}

			val = fmt.Sprint(evalVal)
		} else if mapVal, ok := values[val]; ok {
			val = mapVal
		} else if unicode.IsLetter(rune(val[0])) {
			return val, fmt.Errorf("Undefined variable '%s'", val)
		}

		return val, nil
	}
}

func handleBlock() error {
	// localValues := make(map[string]string)

	for {
		lineNumber++
		if lineNumber >= len(lines) {
			break
		}

		stmt := lines[lineNumber]

		if isBlockEnd(stmt) {
			return nil
		}

		err := handleStmt(stmt)
		if err != nil {
			return err
		}
	}

	return fmt.Errorf("Error at end: Expect '}' .")
}
