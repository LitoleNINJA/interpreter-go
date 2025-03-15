package main

type Scope struct {
	values    map[string]any
	enclosing *Scope
}

func NewScope(enclosing *Scope) *Scope {
	return &Scope{
		values:    make(map[string]any),
		enclosing: enclosing,
	}
}

func (scope *Scope) setScopeValue(key string, val any) {
	scope.values[key] = val
}

// getScopeValue returns the value of the key
// in this environment or its enclosing environments
func (scope *Scope) getScopeValue(key string) (any, bool) {
	// check for nil scope
	if scope == nil {
		return "", false
	}

	// check current scope
	if val, ok := scope.values[key]; ok {
		return val, true
	}

	// check enclosing scope
	if scope.enclosing != nil {
		return scope.enclosing.getScopeValue(key)
	}

	return nil, false
}

// assignScopeValue sets a key-value pair in the environment if the key already exists.
// If it doesn't exist in this environment, it checks the enclosing environment
// and tries to assign the value there
func (scope *Scope) assignScopeValue(key string, val any) bool {
	// check if key exists in current scope
	if _, ok := scope.values[key]; ok {
		// if found, assign new value
		scope.setScopeValue(key, val)
		return true
	}

	// if not found, check enclosing
	if scope.enclosing != nil {
		return scope.enclosing.assignScopeValue(key, val)
	}

	return false
}
