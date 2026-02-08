package evaluator

import (
	"bytes"
	"io"
	"os"
	"rubo-lang/src/parser"
	"strings"
	"testing"
)

func TestProfiler(t *testing.T) {
	input := `def benchmark(n)
  n
end
`
	l := parser.New(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := NewEnvironment()
	Eval(program, env)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Loop 1001 times
	for i := 0; i < 1001; i++ {
		call := &parser.CallExpression{
			Function: &parser.Identifier{Value: "benchmark"},
			Arguments: []parser.Expression{
				&parser.IntegerLiteral{Value: int64(i)},
			},
		}
		Eval(call, env)
	}

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old

	output := buf.String()
	expected := "OPTIMIZATION TRIGGERED: Function benchmark is hot. Handing off to Rust..."
	if !strings.Contains(output, expected) {
		t.Errorf("expected optimization message not found in output. got=%q", output)
	}

	if env.CallCounts["benchmark"] != 1001 {
		t.Errorf("expected count 1001, got=%d", env.CallCounts["benchmark"])
	}
}
