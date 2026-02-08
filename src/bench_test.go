package main

import (
	"fmt"
	"os/exec"
	"rubo-lang/src/evaluator"
	"rubo-lang/src/orchestrator"
	"rubo-lang/src/parser"
	"testing"
	"time"
)

// Run with: go test -v src/bench_test.go
func TestBenchmark(t *testing.T) {
	// Skip if cargo is not available
	if _, err := exec.LookPath("cargo"); err != nil {
		t.Skip("cargo not found in PATH, skipping benchmark")
	}

	// 1. Setup the heavy mathematical function
	// We use a complex infix expression to put some load on the AST interpreter
	input := `def heavy_math(x, y)
  (x + y) * (x - y) + (x * y) * (x + y) - (y * y)
end`
	l := parser.New(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	env := evaluator.NewEnvironment()
	// Register the orchestrator hook for automatic promotion
	env.OnHotFunction = func(fn *evaluator.Function) {
		lit := &parser.FunctionLiteral{
			Name:       &parser.Identifier{Value: fn.Name},
			Parameters: fn.Parameters,
			Body:       fn.Body,
		}
		err := orchestrator.Optimize(lit, env)
		if err != nil {
			fmt.Printf("Optimization error: %v\n", err)
		}
	}

	// Define the function
	evaluator.Eval(program, env)

	// 2. Prepare the call
	callInput := `heavy_math(100, 50)`
	callL := parser.New(callInput)
	callP := parser.NewParser(callL)
	callProgram := callP.ParseProgram()
	callStmt := callProgram.Statements[0]

	const totalCalls = 5000
	const measureCount = 100
	const hotThreshold = 1000

	durations := make([]time.Duration, totalCalls)

	fmt.Printf("--- Starting Benchmark: 5,000 calls ---\n")

	for i := 0; i < totalCalls; i++ {
		start := time.Now()
		evaluator.Eval(callStmt, env)
		durations[i] = time.Since(start)

		if i == hotThreshold {
			fmt.Printf(">>> Optimization should have triggered at call %d\n", hotThreshold)
		}
	}

	// 3. Capture Cold execution (first 100)
	var coldTotal time.Duration
	for i := 0; i < measureCount; i++ {
		coldTotal += durations[i]
	}
	coldAvg := coldTotal / time.Duration(measureCount)

	// 4. Capture Hot execution (last 100)
	var hotTotal time.Duration
	for i := totalCalls - measureCount; i < totalCalls; i++ {
		hotTotal += durations[i]
	}
	hotAvg := hotTotal / time.Duration(measureCount)

	// 5. Output comparison table
	fmt.Println("\n+---------------------------+---------------------------+")
	fmt.Println("| Phase                     | Average Execution Time    |")
	fmt.Println("+---------------------------+---------------------------+")
	fmt.Printf("| Cold (Go Interpretation)  | %-25v |\n", coldAvg)
	fmt.Printf("| Hot (Native Rust)         | %-25v |\n", hotAvg)
	fmt.Println("+---------------------------+---------------------------+")

	speedup := float64(coldAvg.Nanoseconds()) / float64(hotAvg.Nanoseconds())
	fmt.Printf("\nSpeedup Factor: %.2fx faster\n", speedup)
}
