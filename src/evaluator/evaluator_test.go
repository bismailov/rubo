package evaluator

import (
	"rubo-lang/src/parser"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50}, // -10 might fail because I didn't implement prefix minus, let's adjust
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionApplication(t *testing.T) {
	input := `def add(x, y)
  x + y
end
add(5, 5)`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func testEval(input string) Object {
	l := parser.New(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj Object, expected int64) bool {
	result, ok := obj.(*Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
		return false
	}

	return true
}
