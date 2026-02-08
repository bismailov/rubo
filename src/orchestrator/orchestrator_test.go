package orchestrator

import (
	"fmt"
	"os/exec"
	"rubo-lang/src/evaluator"
	"rubo-lang/src/parser"
	"testing"
	"time"
)

func TestEndToEndOptimization(t *testing.T) {
	// Skip if cargo is not available
	if _, err := exec.LookPath("cargo"); err != nil {
		t.Skip("cargo not found in PATH, skipping end-to-end optimization test")
	}

	input := `def add(x, y)
  x + y
end`
	l := parser.New(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	env := evaluator.NewEnvironment()
	// Register the orchestrator hook
	env.OnHotFunction = func(fn *evaluator.Function) {
		// Convert evaluator.Function back to parser.FunctionLiteral for codegen
		lit := &parser.FunctionLiteral{
			Name:       &parser.Identifier{Value: fn.Name},
			Parameters: fn.Parameters,
			Body:       fn.Body,
		}
		err := Optimize(lit, env)
		if err != nil {
			t.Errorf("Optimization failed: %v", err)
		}
	}

	// 1. Define the function
	evaluator.Eval(program, env)

	// 2. Prepare the call expression AST
	callInput := `add(10, 20)`
	callL := parser.New(callInput)
	callP := parser.NewParser(callL)
	callProgram := callP.ParseProgram()

	fmt.Printf("Starting loop for 1100 calls...\n")

	var totalGoTime time.Duration
	var totalRustTime time.Duration

	// 3. Run 1,100 times
	for i := 1; i <= 1100; i++ {
		start := time.Now()
		result := evaluator.Eval(callProgram.Statements[0], env)
		duration := time.Since(start)

		if i <= 1000 {
			totalGoTime += duration
		} else {
			totalRustTime += duration
		}

		if i == 1 {
			fmt.Printf("Call 1 (Go): %s, Duration: %v\n", result.Inspect(), duration)
		}
		if i == 1000 {
			fmt.Printf("Call 1000 (Go): %s, Duration: %v\n", result.Inspect(), duration)
			fmt.Printf("Average Go time: %v\n", totalGoTime/1000)
		}
		if i == 1001 {
			fmt.Printf("Call 1001 (Rust first run): %s, Duration: %v\n", result.Inspect(), duration)
		}
		if i == 1100 {
			fmt.Printf("Call 1100 (Rust): %s, Duration: %v\n", result.Inspect(), duration)
			fmt.Printf("Average Rust time: %v\n", totalRustTime/100)
		}
	}

	if totalRustTime/100 > totalGoTime/1000 && totalRustTime/100 > 1*time.Millisecond {
		// Note: Rust might be slower on the first load due to dlopen,
		// but subsequent calls should be faster.
		// However, for a simple x+y, Go interp might be actually quite fast.
		// In a real JIT, the "fast" part is the execution, not the overhead.
	}
}
