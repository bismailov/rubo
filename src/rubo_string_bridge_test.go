package main

import (
	"fmt"
	"os"
	"path/filepath"
	"rubo-lang/src/bridge"
	"rubo-lang/src/evaluator"
	"rubo-lang/src/parser"
	"runtime"
	"testing"
)

// Run with: go test -v src/rubo_string_bridge_test.go
func TestRuboStringBridge(t *testing.T) {
	// 1. Setup paths
	cwd, _ := os.Getwd()
	var runtimePath string
	if filepath.Base(cwd) == "src" {
		runtimePath = filepath.Join(cwd, "runtime")
	} else {
		runtimePath = filepath.Join(cwd, "src", "runtime")
	}

	libExt := ".so"
	if runtime.GOOS == "darwin" {
		libExt = ".dylib"
	}
	libPath := filepath.Join(runtimePath, "target", "release", "libruntime"+libExt)

	// 2. Load Rubo script
	var scriptPath string
	if filepath.Base(cwd) == "src" {
		scriptPath = filepath.Base(cwd) // This is wrong, let's fix it
	}

	// Hardcoded for simplicity in this specific test environment
	scriptPath = filepath.Join(cwd, "tests", "string_test.rubo")
	if _, err := os.Stat(scriptPath); err != nil {
		scriptPath = filepath.Join(cwd, "..", "tests", "string_test.rubo")
	}

	input, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script at %s: %v", scriptPath, err)
	}

	// 3. Setup evaluator and register NATIVE string function
	l := parser.New(string(input))
	p := parser.NewParser(l)
	program := p.ParseProgram()

	env := evaluator.NewEnvironment()

	// Manually register the native function to simulate "Hot-Swap" result
	nativeFn, err := bridge.LoadRustStringFunction(libPath, "rubo_string_len")
	if err != nil {
		t.Fatalf("Failed to load rust function: %v", err)
	}
	env.NativeFunctions["rubo_string_len"] = nativeFn

	// 4. Eval
	result := evaluator.Eval(program, env)

	// 5. Verify
	if result == nil {
		t.Fatalf("Result is nil")
	}

	expected := int64(len("Hello from Rubo!"))
	if integer, ok := result.(*evaluator.Integer); ok {
		fmt.Printf("Rubo Script Result: %d (Expected: %d)\n", integer.Value, expected)
		if integer.Value != expected {
			t.Errorf("Expected %d, got %d", expected, integer.Value)
		}
	} else {
		t.Errorf("Expected Integer result, got %T (%+v)", result, result)
	}
}
