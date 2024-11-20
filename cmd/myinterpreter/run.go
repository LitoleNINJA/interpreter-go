package main

import (
	"fmt"
	"os"
	"strings"
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

func getVarDeclaration(stmt []byte) (string, string) {
	stmtString := string(stmt)
	stmtString, _ = strings.CutPrefix(stmtString, "var ")
	split := strings.Split(stmtString, "=")
	key := strings.TrimSpace(split[0])
	value := "nil"
	if len(split) == 2 {
		value = strings.TrimSpace(split[1])
	}

	if val, ok := values[value]; ok {
		value = val
	}

	return key, value
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
			key, val := getVarDeclaration(stmt)
			values[key] = val
			continue
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
