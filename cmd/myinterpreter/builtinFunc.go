package main

import (
	"time"
)

// A built-in function type
type BuiltinFunction struct {
	name     string
	function func(args []Value) (Value, error)
}

var native_functions map[string]*BuiltinFunction

// Implement the Callable interface
func (bf *BuiltinFunction) Call(args []Value) (Value, error) {
	return bf.function(args)
}

func init_native_functions() {
	native_functions = make(map[string]*BuiltinFunction)

	native_functions["clock"] = &BuiltinFunction{
		name: "clock",
		function: func(args []Value) (Value, error) {
			return float64(time.Now().Unix()), nil
		},
	}
}
