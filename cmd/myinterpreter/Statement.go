package main

import "fmt"

type Statement interface {
	String() string
	Execute() (Value, error)
}

/* Types of statements
- Expression statement
- Print statement
- Variable declaration statement
- Block statement
- If statement
- While statement (includes for statement)
*/

type ExprStatement struct {
	expr Expr
}

func (e *ExprStatement) String() string {
	return e.expr.String()
}

func (e *ExprStatement) Execute() (Value, error) {
	val, err := e.expr.Evaluate()
	if err != nil {
		return nil, err
	}

	// expr should end with ';'
	// if
	return val, nil
}

type PrintStatement struct {
	expr Expr
}

func (p *PrintStatement) String() string {
	return fmt.Sprintf("print %s", p.expr)
}

func (p *PrintStatement) Execute() (Value, error) {
	val, err := p.expr.Evaluate()
	if err != nil {
		return nil, err
	}

	if val == nil {
		fmt.Println("nil")
	} else {
		fmt.Printf("%v\n", val)
	}
	return nil, nil
}

type VarStatement struct {
	name Token
	init Expr
}

func (v *VarStatement) String() string {
	initVal := "nil"
	if v.init != nil {
		initVal = v.init.String()
	}
	return fmt.Sprintf("var %s = %s", v.name.lexeme, initVal)
}

func (v *VarStatement) Execute() (Value, error) {
	var val Value
	if v.init != nil {
		var err error
		val, err = v.init.Evaluate()
		if err != nil {
			return nil, err
		}
	}

	// Try to find and update existing variable
	currentScope.setScopeValue(v.name.lexeme, val)

	return val, nil
}

type BlockStatement struct {
	stmts []Statement
}

func (b *BlockStatement) String() string {
	return fmt.Sprintf("{%v}", b.stmts)
}

func (b *BlockStatement) Execute() (Value, error) {
	// Create a new scope for the block
	enclosingScope := currentScope
	currentScope = NewScope(enclosingScope)

	for _, stmt := range b.stmts {
		_, err := stmt.Execute()
		if err != nil {
			return nil, err
		}
	}

	// Restore the enclosing scope
	currentScope = enclosingScope

	return nil, nil
}

type IfStatement struct {
	condition Expr
	thenStmt  Statement
	elseStmt  Statement
}

func (i *IfStatement) String() string {
	return fmt.Sprintf("(if %s then %s, else %s)", i.condition, i.thenStmt, i.elseStmt)
}

func (i *IfStatement) Execute() (Value, error) {
	conditionResult, err := i.condition.Evaluate()
	if err != nil {
		return nil, err
	}

	if isTruthy(conditionResult) && i.thenStmt != nil {
		_, err = i.thenStmt.Execute()
	} else if i.elseStmt != nil {
		_, err = i.elseStmt.Execute()
	}

	return nil, err
}

type WhileStatement struct {
	condition Expr
	body      Statement
}

func (w *WhileStatement) String() string {
	return fmt.Sprintf("(while (%s) do %s)", w.condition, w.body)
}

func (w *WhileStatement) Execute() (Value, error) {
	conditionResult, err := w.condition.Evaluate()
	if err != nil {
		return nil, err
	}

	for isTruthy(conditionResult) {
		_, err = w.body.Execute()
		if err != nil {
			return nil, err
		}

		conditionResult, err = w.condition.Evaluate()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
