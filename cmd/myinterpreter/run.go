package main

import (
	"fmt"
	"strings"
)

func getPrintContents(line string) string {
	index := strings.Index(line, "print ")
	if index == -1 {
		return ""
	}

	return line[index+6:]
}

func fomatLine(line string) string {
	line = strings.TrimSpace(line)

	return line
}

func readPrintStmt(fileContent []byte) [][]byte {
	fileString := string(fileContent)
	len := len(fileString)

	var lines [][]byte
	for i := 0; i < len; i++ {
		if i < len-5 && fileString[i:i+5] == "print" {
			index := strings.Index(fileString[i:], ";")
			line := fileString[i : i+index]
			line = getPrintContents(line)
			line = fomatLine(line)
			lines = append(lines, []byte(line))
		}
	}

	return lines
}

func run(fileContents []byte) error {
	lines := readPrintStmt(fileContents)

	for _, stmt := range lines {
		expr, err := evaluate(stmt)
		if err != nil {
			return err
		}

		fmt.Println(expr)
	}
	return nil
}
