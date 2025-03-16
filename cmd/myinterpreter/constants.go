package main

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

const max_args_count = 127

var currentScope = NewScope(nil)
