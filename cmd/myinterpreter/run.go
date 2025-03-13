package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var lines [][]byte
var lineNumber int
var isString bool
var currentScope *Scope
var mappedStmts map[int]string

type StatementType int

const (
	PrintStatement StatementType = iota
	VarDeclarationStatement
	AssignmentStatement
	BlockStartStatement
	BlockEndStatement
	IfStatement
	ElseStatement
	CommentStatement
	EmptyStatement
	ComplexStatement
	WhileStatement
	ForStatement
)

type Statement struct {
	Type StatementType
	Stmt []byte
}

// DetermineStatementType determines the type of a statement
func DetermineStatementType(stmt []byte) StatementType {
	stmtString := string(stmt)

	if strings.HasPrefix(stmtString, "//") {
		return CommentStatement
	} else if strings.HasPrefix(stmtString, "(") {
		return ComplexStatement
	} else if strings.HasPrefix(stmtString, "print") {
		return PrintStatement
	} else if strings.HasPrefix(stmtString, "var ") {
		return VarDeclarationStatement
	} else if strings.TrimSpace(stmtString) == "{" {
		return BlockStartStatement
	} else if strings.TrimSpace(stmtString) == "}" {
		return BlockEndStatement
	} else if strings.HasPrefix(stmtString, "if ") {
		return IfStatement
	} else if strings.HasPrefix(stmtString, "else ") || strings.HasPrefix(stmtString, "} else") {
		return ElseStatement
	} else if strings.HasPrefix(stmtString, "while ") {
		return WhileStatement
	} else if strings.HasPrefix(stmtString, "for ") {
		return ForStatement
	} else if isAssignment(stmt) {
		return AssignmentStatement
	}

	return -1 // Unknown statement type
}

// readLines reads the file content and splits it into stmt lines
func readLines(fileContent []byte) [][]byte {
	var lines [][]byte
	line := make([]byte, 0)
	cnt := 0

	for i := 0; i < len(fileContent); i++ {
		ch := fileContent[i]

		if ch == ')' {
			cnt++
		} else if ch == '(' {
			cnt--
		}

		if ch == '"' {
			isString = !isString
			line = append(line, ch)
			continue
		}
		if isString {
			line = append(line, ch)
			continue
		}

		if (ch == ';' && cnt == 0) || ch == '}' {
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

// isAssignment checks if a statement is an assignment statement
func isAssignment(stmt []byte) bool {
	stmtString := string(stmt)

	if !strings.Contains(stmtString, "=") {
		return false
	}
	if strings.Contains(stmtString, "==") ||
		strings.Contains(stmtString, ">=") ||
		strings.Contains(stmtString, "<=") ||
		strings.Contains(stmtString, "!=") {
		return false
	}
	return true
}

// getPrintContents extracts the content to be printed from a print statement
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

// getVarDeclaration extracts the key and value from a var declaration statement
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

// getElseStmt extracts the else statement, similar to getIfStmt
func getElseStmt(stmt []byte) ([]byte, bool, error) {
	// fmt.Printf("stmt : %s\n", stmt)

	if s, ok := strings.CutPrefix(string(stmt), "else "); ok {
		if strings.HasPrefix(s, "if ") {
			_, body, err := extractCondition([]byte(s), "if ")
			if err != nil {
				return []byte{}, false, err
			}

			return []byte(strings.TrimSpace(s)), string(body) == "{", nil
		} else if strings.HasPrefix(s, "{") {
			restOfBody := strings.TrimSpace(strings.TrimPrefix(s, "{"))
			s = "{"
			lines = append(lines[:lineNumber+1], append([][]byte{[]byte(restOfBody)}, lines[lineNumber+1:]...)...)
			return []byte(s), true, nil
		}
		return []byte(strings.TrimSpace(s)), false, nil
	}

	return []byte{}, false, fmt.Errorf("else stmt not found")
}

func isBlockEnd(stmt []byte) bool {
	return DetermineStatementType(stmt) == BlockEndStatement
}

func isIfStmt(stmt []byte) bool {
	return DetermineStatementType(stmt) == IfStatement
}

func isElseStmt(stmt []byte) bool {
	return DetermineStatementType(stmt) == ElseStatement
}

// main entry point of the interpreter
func run(fileContents []byte) error {
	lineNumber = 0
	lines = readLines(fileContents)
	currentScope = NewScope(nil)

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

// handleStmt processes a statement and executes it accordingly
func handleStmt(stmt []byte) error {
	printStmt := false
	stmtType := DetermineStatementType(stmt)

	if stmtType == CommentStatement {
		return nil
	} else if stmtType == ComplexStatement {
		return handleComplexStmt(stmt)
	} else if stmtType == PrintStatement {
		printStmt = true
		stmt = getPrintContents(stmt)
	} else if stmtType == VarDeclarationStatement {
		return getVarDeclaration(stmt)
	} else if stmtType == BlockStartStatement {
		return handleBlock()
	} else if stmtType == IfStatement {
		return handleIfBlock(stmt)
	} else if stmtType == WhileStatement {
		return handleWhileBlock(stmt)
	} else if stmtType == ForStatement {
		return handleForBlock(stmt)
	}

	if stmtType == ElseStatement {
		_, t, _ := getElseStmt(stmt)
		if t {
			for lineNumber < len(lines) && !isBlockEnd(lines[lineNumber]) {
				// fmt.Printf("skipping else stmt : %s\n", lines[lineNumber])
				lineNumber++
			}
		}
		return nil
	}

	if isAssignment(stmt) {
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

// handleAssignment processes an assignment statement and assigns the value to a variable, and returns the assigned value
func handleAssignment(stmt string) (string, error) {
	stmt, _ = strings.CutSuffix(stmt, ";")

	if isAssignment([]byte(stmt)) {
		pos := strings.Index(stmt, "=")

		key := strings.TrimSpace(stmt[:pos])
		val, err := handleAssignment(stmt[pos+1:])
		if err != nil {
			return val, err
		}

		// fmt.Printf("Key : %s, Value : %s\n", key, val)
		// Try to find and update existing variable
		if success := currentScope.assignScopeValue(key, val); !success {
			currentScope.setScopeValue(key, val)
		}
		return val, nil
	} else {
		val := strings.TrimSpace(stmt)
		if strings.ContainsAny(val, "+-*/()><") {
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

// handleBlock processes a block of code enclosed in curly braces {}.
// It creates a new scope for the block, executes all statements within the block,
// and restores the enclosing scope when the block ends.
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

// handleIfBlock processes an if statement including its condition and body.
func handleIfBlock(stmt []byte) error {
	condition, body, err := extractCondition(stmt, "if ")
	if err != nil {
		return err
	}

	// Evaluate condition
	conditionResult, err := evaluateCondition(condition)
	if err != nil {
		return err
	}

	// Handle if/else logic
	if isTruthy(conditionResult) {
		return handleStmt(body)
	} else {
		// If condition is false, skip the if block
		if bytes.Equal(body, []byte("{")) {
			lineNumber, _ = findBlockEnd()
		}
		lineNumber++

		// Check for else or else if statement
		if lineNumber < len(lines) && isElseStmt(lines[lineNumber]) {
			elseBody, _, err := getElseStmt(lines[lineNumber])
			if err != nil {
				return err
			}

			// Handle else if block
			if isIfStmt(elseBody) {
				lines = append(lines[:lineNumber+1], append([][]byte{elseBody}, lines[lineNumber+1:]...)...)
				lineNumber++
				return handleIfBlock(elseBody)
			}

			// Handle regular else block
			return handleStmt(elseBody)
		} else if lineNumber < len(lines) {
			lineNumber--
		}
	}
	return nil
}

// handleComplexStmt processes complex statements that contain nested parentheses.
// It maps nested expressions to simple placeholders and processes them accordingly.
func handleComplexStmt(stmt []byte) error {
	var simpleStmt []byte
	simpleStmt, mappedStmts = mapComplexStmt(stmt)
	// fmt.Printf("Mapped : %v\nStmt : %s\n", mappedStmts, simpleStmt)

	if strings.Contains(string(simpleStmt), "or") {
		return handleComplexOrStmt(simpleStmt)
	} else if strings.Contains(string(simpleStmt), "and") {
		return handleComplexAndStmt(simpleStmt)
	}

	_, err := evaluate(stmt)

	return err
}

// handleComplexOrStmt processes a complex OR expression with short-circuit evaluation.
func handleComplexOrStmt(stmt []byte) error {
	parts := strings.Split(string(stmt), "or")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, "%map") {
			mapKey, err := strconv.Atoi(part[4:])
			if err != nil {
				return err
			}

			mapVal, ok := mappedStmts[mapKey]
			if !ok {
				return fmt.Errorf("key %d not found in mappedStmts", mapKey)
			}

			if isAssignment([]byte(mapVal)) {
				result, err := handleAssignment(mapVal)
				if err != nil {
					return err
				}

				// If this result is truthy, we can short-circuit
				if isTruthy(result) {
					return nil
				}
			} else {
				// Evaluate as a normal expression
				result, err := evaluate([]byte(mapVal))
				if err != nil {
					return err
				}

				// If this result is truthy, we can short-circuit
				if isTruthy(result) {
					return nil
				}
			}
		} else {
			result, err := evaluate([]byte(part))
			if err != nil {
				return err
			}

			if isTruthy(result) {
				return nil
			}
		}

	}
	return nil
}

// handleComplexAndStmt processes a complex AND expression with short-circuit evaluation.
func handleComplexAndStmt(stmt []byte) error {
	parts := strings.Split(string(stmt), "and")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, "%map") {
			mapKey, err := strconv.Atoi(part[4:])
			if err != nil {
				return err
			}

			mapVal, ok := mappedStmts[mapKey]
			if !ok {
				return fmt.Errorf("key %d not found in mappedStmts", mapKey)
			}

			if isAssignment([]byte(mapVal)) {
				result, err := handleAssignment(mapVal)
				if err != nil {
					return err
				}

				// If this result is falsy, we can short-circuit
				if !isTruthy(result) {
					return nil
				}
			} else {
				// Evaluate as a normal expression
				result, err := evaluate([]byte(mapVal))
				if err != nil {
					return err
				}

				// If this result is falsy, we can short-circuit
				if !isTruthy(result) {
					return nil
				}
			}
		} else {
			result, err := evaluate([]byte(part))
			if err != nil {
				return err
			}

			if !isTruthy(result) {
				return nil
			}
		}

	}
	return nil
}

func handleWhileBlock(stmt []byte) error {
	condition, stmt, err := extractCondition(stmt, "while ")
	if err != nil {
		return err
	}

	// fmt.Printf("While Condition : %s, Stmt : %s\n", condition, stmt)

	conditionResult, err := evaluateCondition(condition)
	if err != nil {
		return err
	}

	// check for block while
	blockEndLineNumber := 0
	if bytes.Equal(stmt, []byte("{")) {
		blockEndLineNumber, err = findBlockEnd()
		if err != nil {
			return err
		}
	} else {
		// it is a single line while statement
		for isTruthy(conditionResult) {
			err := handleStmt(stmt)
			if err != nil {
				return err
			}

			conditionResult, err = evaluateCondition(condition)
			if err != nil {
				return err
			}
		}

		return nil
	}

	blockSize := blockEndLineNumber - lineNumber
	// evaluate stmts till while block ends
	for isTruthy(conditionResult) {
		// process the for block
		err = handleBlock()
		if err != nil {
			return err
		}

		// skip back to while start
		lineNumber -= blockSize

		conditionResult, err = evaluateCondition(condition)
		if err != nil {
			return err
		}
	}

	lineNumber = blockEndLineNumber

	return nil
}

func handleForBlock(stmt []byte) error {
	body, stmt, err := extractCondition(stmt, "for ")
	if err != nil {
		return err
	}

	// fmt.Printf("For Condition : %s, Stmt : %s\n", body, stmt)
	init, condition, updation, err := parseForStmt(body)
	if err != nil {
		return err
	}

	err = handleStmt(init)
	if err != nil {
		return err
	}

	conditionResult, err := evaluateCondition(condition)
	if err != nil {
		return err
	}

	// check for block for
	blockEndLineNumber := 0
	if bytes.Equal(stmt, []byte("{")) {
		blockEndLineNumber, err = findBlockEnd()
		if err != nil {
			return err
		}
	} else {
		// it is a single line for statement
		for isTruthy(conditionResult) {
			err := handleStmt(stmt)
			if err != nil {
				return err
			}

			err = handleStmt(updation)
			if err != nil {
				return err
			}

			conditionResult, err = evaluateCondition(condition)
			if err != nil {
				return err
			}
		}

		return nil
	}

	blockSize := blockEndLineNumber - lineNumber
	// evaluate stmts till for block ends
	for isTruthy(conditionResult) {
		// process the for block
		err = handleBlock()
		if err != nil {
			return err
		}

		// skip back to for start
		lineNumber -= blockSize

		err = handleStmt(updation)
		if err != nil {
			return err
		}

		conditionResult, err = evaluateCondition(condition)
		if err != nil {
			return err
		}
	}

	lineNumber = blockEndLineNumber

	return nil
}
