package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

var values map[string]string

func getPrintContents(line []byte) []byte {
	s, ok := strings.CutPrefix(string(line), "print")
	if !ok {
		fmt.Printf("Print line dosent start with Print : %s\n", line)
		return []byte{}
	}

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

func readLines(fileContent []byte) [][]byte {
	var lines [][]byte
	line := make([]byte, 0)
	for i := 0; i < len(fileContent); i++ {
		if fileContent[i] == 59 {
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

func isVarDeclaration(stmt []byte) bool {
	stmtString := string(stmt)
	return strings.HasPrefix(stmtString, "var ")
}

func getVarDeclaration(stmt []byte) error {
	stmtString := string(stmt)
	stmtString, _ = strings.CutPrefix(stmtString, "var ")

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

func run(fileContents []byte) error {
	lines := readLines(fileContents)
	values = make(map[string]string)
	// fmt.Println(lines)
	for _, stmt := range lines {
		printStmt := false
		if isPrintStmt(stmt) {
			printStmt = true
			stmt = getPrintContents(stmt)
		} else if isVarDeclaration(stmt) {
			err := getVarDeclaration(stmt)
			if err != nil {
				return err
			}
			continue
		}

		if strings.Contains(string(stmt), "=") {
			handleAssignment(string(stmt))
		}

		if len(stmt) == 0 {
			os.Exit(65)
		}

		// fmt.Printf("Eval : %s, Len : %d\n", stmt, len(stmt))
		expr, err := evaluate(stmt)
		if err != nil {
			return err
		}

		if printStmt {
			fmt.Println(expr)
		}
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
