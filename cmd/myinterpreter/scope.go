package main

type Scope struct {
	values    map[string]string
	enclosing *Scope
}

func NewScope(enclosing *Scope) *Scope {
	return &Scope{
		values:    make(map[string]string),
		enclosing: enclosing,
	}
}

func (scope *Scope) setScopeValue(key string, val string) {
	scope.values[key] = val
}

func (scope *Scope) getScopeValue(key string) (string, bool) {
	// check current scope
	if val, ok := scope.values[key]; ok {
		return val, true
	}

	// check enclosing scope
	if scope.enclosing != nil {
		return scope.enclosing.getScopeValue(key)
	}

	return "", false
}
