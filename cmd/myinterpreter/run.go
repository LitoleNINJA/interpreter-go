package main

import (
	"fmt"
	"os"
	"strings"
)

func getPrintContents(line []byte) []byte {
	s, ok := strings.CutPrefix(string(line), "print")
	if !ok {
		fmt.Printf("Print line dosent start with Print : %s\n", line)
		return []byte{}
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

func run(fileContents []byte) error {
	lines := readLines(fileContents)
	// fmt.Println(lines)
	for _, stmt := range lines {
		printStmt := false
		if isPrintStmt(stmt) {
			printStmt = true
			stmt = getPrintContents(stmt)
		}
		if len(stmt) == 0 {
			fmt.Println("Empty statement !")
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
