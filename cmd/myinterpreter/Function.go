package main

import "fmt"

type Function struct {
	stmt  FuncStatement
	scope *Scope
}

func (fn *Function) String() string {
	return fmt.Sprintf("<fn %s>", fn.stmt.name.lexeme)
}

func (fn *Function) Call(args []Value) (returnVal Value, err error) {
	previousScope := currentScope
	currentScope = NewScope(fn.scope)
	defer func() {
		currentScope = previousScope
		if err := recover(); err != nil {
			if ret, ok := err.(*Return); ok {
				returnVal = ret.val
				return
			}
			panic(err)
		}
	}()

	// Set the arguments in the new scope
	for i, arg := range fn.stmt.args {
		currentScope.setScopeValue(arg.lexeme, args[i])
	}

	for _, stmt := range fn.stmt.body {
		_, err = stmt.Execute()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

type FunctionExpr struct {
	expr  Func
	scope *Scope
}

func (fn *FunctionExpr) String() string {
	return fmt.Sprintf("(fn %s %s)", fn.expr.name.lexeme, fn.expr.args)
}

func (fn *FunctionExpr) Call(args []Value) (Value, error) {
	previousScope := currentScope
	currentScope = NewScope(fn.scope)
	defer func() {
		currentScope = previousScope
	}()

	for _, arg := range fn.expr.args {
		currentScope.setScopeValue(arg.lexeme, arg)
	}

	for _, stmt := range fn.expr.body {
		_, err := stmt.Execute()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

type Return struct {
	val Value
}