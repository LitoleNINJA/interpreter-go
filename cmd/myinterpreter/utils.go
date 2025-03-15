package main

import (
	"fmt"
	"strconv"
)

func isAlphaNum(s string) bool {
	return isIndentifierStart(s) || isStringDigit(s)
}

func getStringType(s string) string {
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return "number"
	} else if s == "true" || s == "false" {
		return "bool"
	} else if s == "nil" {
		return "nil"
	} else {
		return "string"
	}
}

func readString(index *int) (string, error) {
	j := *index + 1
	for j < len(fileContentString) && fileContentString[j] != '"' {
		j++
	}

	if j < len(fileContentString) && fileContentString[j] == '"' {
		str := fileContentString[*index+1 : j]
		*index = j
		return str, nil
	}

	*index = j
	return "", fmt.Errorf("Error: Unterminated string.")
}

func readNumber(index *int) (string, string) {
	str := ""
	frac := ""
	for *index < len(fileContentString) && isStringDigit(string(fileContentString[*index])) {
		str += string(fileContentString[*index])
		*index++
	}

	if *index < len(fileContentString)-1 && fileContentString[*index] == '.' && isStringDigit(string(fileContentString[*index+1])) {
		str += "."
		*index++
		for *index < len(fileContentString) && isStringDigit(string(fileContentString[*index])) {
			str += string(fileContentString[*index])
			frac += string(fileContentString[*index])
			*index++
		}
	}

	if i, _ := strconv.Atoi(frac); i == 0 {
		frac = "0"
	}

	return str, frac
}

func readIdentifier(index *int) string {
	str := ""
	for *index < len(fileContentString) && isAlphaNum(string(fileContentString[*index])) {
		str += string(fileContentString[*index])
		*index++
	}

	return str
}

func isStringDigit(s string) bool {
	if len(s) == 1 && s >= "0" && s <= "9" {
		return true
	}
	return false
}

func isIndentifierStart(s string) bool {
	return len(s) == 1 && (s >= "a" && s <= "z") || (s >= "A" && s <= "Z") || (s == "_")
}

func add(left Value, right Value) (Value, error) {
	switch left.(type) {
	case float64:
		leftVal := left.(float64)
		if rightVal, ok := right.(float64); !ok {
			exitCode = 70
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		} else {
			return leftVal + rightVal, nil
		}
	case string:
		leftVal := left.(string)
		if rightVal, ok := right.(string); !ok {
			exitCode = 70
			return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
		} else {
			return leftVal + rightVal, nil
		}
	default:
		exitCode = 70
		return nil, fmt.Errorf("operands must be numbers.\n[line 1]")
	}
}

func checkEqual(leftVal Value, rightVal Value) bool {
	switch left := leftVal.(type) {
	case float64:
		if right, ok := rightVal.(float64); !ok {
			return false
		} else {
			return left == right
		}
	case string:
		if right, ok := rightVal.(string); !ok {
			return false
		} else {
			return left == right
		}
	case bool:
		if right, ok := rightVal.(bool); !ok {
			return false
		} else {
			return left == right
		}
	default:
		fmt.Println("Type mismatch !")
		return false
	}
}

func checkBothNumber(leftVal Value, rightVal Value) error {
	switch leftVal.(type) {
	case float64:
		if _, ok := rightVal.(float64); !ok {
			exitCode = 70
			return fmt.Errorf("operands must be numbers.\n[line 1]")
		} else {
			return nil
		}
	default:
		exitCode = 70
		return fmt.Errorf("operands must be numbers.\n[line 1]")
	}
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
