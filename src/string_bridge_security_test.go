package main

import (
	"fmt"
	"os"
	"path/filepath"
	"rubo-lang/src/bridge"
	"rubo-lang/src/evaluator"
	"rubo-lang/src/parser"
	"runtime"
	"strings"
	"testing"
)

func TestStringBridgeHardening(t *testing.T) {
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

	t.Run("PanicHandling", func(t *testing.T) {
		// Load a function that we know will panic in Rust
		fn, err := bridge.LoadRustStringFunction(libPath, "rubo_trigger_panic")
		if err != nil {
			t.Fatalf("Failed to load rust function: %v", err)
		}

		env := evaluator.NewEnvironment()
		// Register it in the evaluator
		env.NativeFunctions["trigger_panic"] = fn

		// Rubo code that calls the panicking function
		input := `def trigger_panic(s); end; trigger_panic("test")`
		l := parser.New(input)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		result := evaluator.Eval(program, env)

		// We expect an ERROR_OBJ
		if result.Type() != evaluator.ERROR_OBJ {
			t.Errorf("Expected ERROR_OBJ, got %s", result.Type())
		}

		errObj := result.(*evaluator.Error)
		if !strings.Contains(errObj.Message, "RuboPanic") {
			t.Errorf("Expected RuboPanic error message, got: %s", errObj.Message)
		}
		fmt.Printf("Caught expected Rust panic: %s\n", errObj.Message)
	})

	t.Run("NormalOperation", func(t *testing.T) {
		fn, err := bridge.LoadRustStringFunction(libPath, "rubo_string_len")
		if err != nil {
			t.Fatalf("Failed to load rust function: %v", err)
		}

		env := evaluator.NewEnvironment()
		env.NativeFunctions["string_len"] = fn

		input := `def string_len(s); end; string_len("Hardening Test")`
		l := parser.New(input)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		result := evaluator.Eval(program, env)

		if integer, ok := result.(*evaluator.Integer); ok {
			if integer.Value != 14 {
				t.Errorf("Expected 14, got %d", integer.Value)
			}
		} else {
			t.Errorf("Expected Integer, got %T", result)
		}
		fmt.Printf("Normal operation verified after hardening.\n")
	})
}
