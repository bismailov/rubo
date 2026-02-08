package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rubo-lang/src/parser"
	"strings"
)

// GenerateRust takes a Rubo function AST and converts it to a Rust function string.
func GenerateRust(fn *parser.FunctionLiteral) string {
	var params []string
	for _, p := range fn.Parameters {
		params = append(params, fmt.Sprintf("%s: i64", p.Value))
	}

	var body string
	if fn.Body != nil && len(fn.Body.Statements) > 0 {
		// In Ruby/Rubo, the last expression is the return value.
		lastStmt := fn.Body.Statements[len(fn.Body.Statements)-1]
		if exprStmt, ok := lastStmt.(*parser.ExpressionStatement); ok {
			body = generateExpression(exprStmt.Expression)
		}
	}

	return fmt.Sprintf(`#[no_mangle]
pub extern "C" fn %s(%s) -> i64 {
    %s
}`, fn.Name.Value, strings.Join(params, ", "), body)
}

func generateExpression(expr parser.Expression) string {
	switch e := expr.(type) {
	case *parser.Identifier:
		return e.Value
	case *parser.IntegerLiteral:
		return fmt.Sprintf("%d", e.Value)
	case *parser.InfixExpression:
		return fmt.Sprintf("(%s %s %s)", generateExpression(e.Left), e.Operator, generateExpression(e.Right))
	case *parser.PrefixExpression:
		return fmt.Sprintf("(%s%s)", e.Operator, generateExpression(e.Right))
	case *parser.CallExpression:
		var args []string
		for _, arg := range e.Arguments {
			args = append(args, generateExpression(arg))
		}
		return fmt.Sprintf("%s(%s)", generateExpression(e.Function), strings.Join(args, ", "))
	}
	return ""
}

// WriteRustFile saves the generated Rust code to src/runtime/src/hot_func.rs.
func WriteRustFile(filename string, code string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(code), 0644)
}

func CompileRust() error {
	// Find the project root by looking for "src/runtime/Cargo.toml"
	cwd, _ := os.Getwd()
	var runtimePath string

	// Possible paths to check relative to cwd
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
		// Fallback to searching upwards for the root containing src/runtime/Cargo.toml
		tempDir := cwd
		for {
			p := filepath.Join(tempDir, "src", "runtime")
			if _, err := os.Stat(filepath.Join(p, "Cargo.toml")); err == nil {
				runtimePath = p
				break
			}
			parent := filepath.Dir(tempDir)
			if parent == tempDir {
				return fmt.Errorf("could not find src/runtime/Cargo.toml directory starting from %s", cwd)
			}
			tempDir = parent
		}
	}

	cmd := exec.Command("cargo", "build", "--release")
	cmd.Dir = runtimePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
