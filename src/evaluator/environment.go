package evaluator

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]int)
	n := make(map[string]interface{})
	return &Environment{store: s, CallCounts: c, NativeFunctions: n, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	// Native functions are typically global or shared
	if outer != nil {
		env.NativeFunctions = outer.NativeFunctions
		env.CallCounts = outer.CallCounts
	}
	return env
}

type Environment struct {
	store           map[string]Object
	outer           *Environment
	CallCounts      map[string]int
	NativeFunctions map[string]interface{}
	OnHotFunction   func(fn *Function)
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string) Object {
	// Note: Basic set, usually you'd pass the value too. Let's fix the signature.
	return nil
}

func (e *Environment) SetValue(name string, val Object) Object {
	e.store[name] = val
	return val
}
