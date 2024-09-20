package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	LEFT_PAREN    = "LEFT_PAREN"
	RIGHT_PAREN   = "RIGHT_PAREN"
	LEFT_BRACE    = "LEFT_BRACE"
	RIGHT_BRACE   = "RIGHT_BRACE"
	STAR          = "STAR"
	DOT           = "DOT"
	COMMA         = "COMMA"
	PLUS          = "PLUS"
	MINUS         = "MINUS"
	SEMICOLON     = "SEMICOLON"
	EQUAL         = "EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	BANG          = "BANG"
	BANG_EQUAL    = "BANG_EQUAL"
	LESS          = "LESS"
	LESS_EQUAL    = "LESS_EQUAL"
	GREATER       = "GREATER"
	GREATER_EQUAL = "GREATER_EQUAL"
	SLASH         = "SLASH"
	STRING        = "STRING"
	NUMBER        = "NUMBER"
	IDENTIFIER    = "IDENTIFIER"
	AND           = "AND"
	CLASS         = "CLASS"
	ELSE          = "ELSE"
	FALSE         = "FALSE"
	FOR           = "FOR"
	FUN           = "FUN"
	IF            = "IF"
	NIL           = "NIL"
	OR            = "OR"
	PRINT         = "PRINT"
	RETURN        = "RETURN"
	SUPER         = "SUPER"
	THIS          = "THIS"
	TRUE          = "TRUE"
	VAR           = "VAR"
	WHILE         = "WHILE"
)

var keywords = map[string]string{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Token struct {
	TokenType string
	lexeme    string
	literal   string
}

func (token *Token) setToken(args ...string) {
	token.TokenType = args[0]
	token.lexeme = args[1]
	if len(args) > 2 {
		token.literal = args[2]
	} else {
		token.literal = "null"
	}
}

func (token *Token) printToken() {
	fmt.Printf("%s %s %s\n", token.TokenType, token.lexeme, token.literal)
}

var exitCode = 0
var fileContentString string
var line = 1

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tokenize":
		tokens := tokenizeFile(fileContents)

		for _, token := range tokens {
			if token != (Token{}) {
				token.printToken()
			}
		}

		fmt.Println("EOF  null")
	case "parse":
		expr, err := parseFile(fileContents)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(expr)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func tokenizeFile(fileContents []byte) []Token {
	fileContentString = string(fileContents)

	tokens := []Token{}
	for i := 0; i < len(fileContentString); i++ {
		newToken := addToken(string(fileContentString[i]), &i)
		if newToken != (Token{}) {
			tokens = append(tokens, newToken)
		}
	}

	return tokens
}

func addToken(ch string, index *int) Token {
	token := Token{}

	switch ch {
	case "(":
		token.setToken(LEFT_PAREN, ch)
	case ")":
		token.setToken(RIGHT_PAREN, ch)
	case "{":
		token.setToken(LEFT_BRACE, ch)
	case "}":
		token.setToken(RIGHT_BRACE, ch)
	case ",":
		token.setToken(COMMA, ch)
	case ".":
		token.setToken(DOT, ch)
	case "*":
		token.setToken(STAR, ch)
	case "+":
		token.setToken(PLUS, ch)
	case "-":
		token.setToken(MINUS, ch)
	case ";":
		token.setToken(SEMICOLON, ch)
	case "=":
		if nextToken(*index) == "=" {
			token.setToken(EQUAL_EQUAL, "==")
			*index++
		} else {
			token.setToken(EQUAL, ch)
		}
	case "!":
		if nextToken(*index) == "=" {
			token.setToken(BANG_EQUAL, "!=")
			*index++
		} else {
			token.setToken(BANG, ch)
		}
	case "<":
		if nextToken(*index) == "=" {
			token.setToken(LESS_EQUAL, "<=")
			*index++
		} else {
			token.setToken(LESS, ch)
		}
	case ">":
		if nextToken(*index) == "=" {
			token.setToken(GREATER_EQUAL, ">=")
			*index++
		} else {
			token.setToken(GREATER, ch)
		}
	case "/":
		if nextToken(*index) != "/" {
			token.setToken(SLASH, ch)
		} else {
			for *index < len(fileContentString) && fileContentString[*index] != '\n' {
				*index++
			}
			line++
		}
	case " ", "\t":
		break
	case "\n":
		line++
	case `"`:
		str := readString(index)
		if str == "" {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line)
			exitCode = 65
			break
		} else {
			token.setToken(STRING, `"`+str+`"`, str)
		}
	default:
		if isStringDigit(ch) {
			str, frac := readNumber(index)
			if frac == "0" {
				floatVal, _ := strconv.ParseFloat(str, 64)
				intVal := int64(floatVal)
				token.setToken(NUMBER, str, strconv.FormatInt(intVal, 10)+"."+frac)
			} else {
				token.setToken(NUMBER, str, strings.Trim(str, "0"))
			}
			*index--
		} else if isIndentifierStart(ch) {
			str := readIdentifier(index)
			if _, isKeyword := keywords[str]; isKeyword {
				token.setToken(keywords[str], str)
			} else {
				token.setToken(IDENTIFIER, str)
			}
			*index--
		} else {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", line, ch)
			exitCode = 65
		}
	}

	return token
}

func nextToken(index int) string {
	if index < len(fileContentString)-1 {
		return string(fileContentString[index+1])
	}
	return ""
}

func readString(index *int) string {
	str := ""
	*index++
	for *index < len(fileContentString) && fileContentString[*index] != '"' {
		str += string(fileContentString[*index])
		*index++
	}

	if *index < len(fileContentString) && fileContentString[*index] == '"' {
		return str
	}
	return ""
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
	return (s >= "a" && s <= "z") || (s >= "A" && s <= "Z") || (s == "_")
}

func isAlphaNum(s string) bool {
	return isIndentifierStart(s) || isStringDigit(s)
}
