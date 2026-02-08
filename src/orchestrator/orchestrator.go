package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"rubo-lang/src/bridge"
	"rubo-lang/src/codegen"
	"rubo-lang/src/evaluator"
	"rubo-lang/src/parser"
	"runtime"
)

// Optimize takes a function literal and an environment, and promotes the function to a native Rust implementation.
func Optimize(fn *parser.FunctionLiteral, env *evaluator.Environment) error {
	fmt.Printf("[Orchestrator] Starting optimization for function: %s\n", fn.Name.Value)

	// 1. Generate Rust code
	rustCode := codegen.GenerateRust(fn)

	// 2. Find the project root by looking for "src/runtime/Cargo.toml"
	cwd, _ := os.Getwd()
	var runtimePath string

	checks := []string{
		filepath.Join(cwd, "src", "runtime"),
		filepath.Join(cwd, "runtime"),
		filepath.Join(cwd, "..", "src", "runtime"),
	}

	for _, p := range checks {
		if _, err := os.Stat(filepath.Join(p, "Cargo.toml")); err == nil {
			runtimePath = p
			break
		}
	}

	if runtimePath == "" {
		tempDir := cwd
		for {
			p := filepath.Join(tempDir, "src", "runtime")
			if _, err := os.Stat(filepath.Join(p, "Cargo.toml")); err == nil {
				runtimePath = p
				break
			}
			parent := filepath.Dir(tempDir)
			if parent == tempDir {
				return fmt.Errorf("could not find src/runtime directory starting from %s", cwd)
			}
			tempDir = parent
		}
	}

	hotFuncPath := filepath.Join(runtimePath, "src", "hot_func.rs")

	err := codegen.WriteRustFile(hotFuncPath, rustCode)
	if err != nil {
		return fmt.Errorf("failed to write rust file: %v", err)
	}

	// 3. Compile Rust
	err = codegen.CompileRust()
	if err != nil {
		return fmt.Errorf("failed to compile rust: %v", err)
	}

	// 4. Load Rust function
	libExt := ".so"
	if runtime.GOOS == "darwin" {
		libExt = ".dylib"
	} else if runtime.GOOS == "windows" {
		libExt = ".dll"
	}

	libPath := filepath.Join(runtimePath, "target", "release", "libruntime"+libExt)

	nativeFn, err := bridge.LoadRustFunction(libPath, fn.Name.Value)
	if err != nil {
		return fmt.Errorf("failed to load rust function: %v", err)
	}

	// 5. Update the Evaluator's NativeFunctions map
	env.NativeFunctions[fn.Name.Value] = nativeFn

	fmt.Printf("[Orchestrator] Function %s successfully optimized!\n", fn.Name.Value)
	return nil
}
