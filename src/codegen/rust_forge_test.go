package codegen

import (
	"os/exec"
	"rubo-lang/src/parser"
	"strings"
	"testing"
)

func TestGenerateRust(t *testing.T) {
	input := `def add(x, y)
  x + y
end`

	l := parser.New(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	stmt := program.Statements[0].(*parser.ExpressionStatement)
	fn := stmt.Expression.(*parser.FunctionLiteral)

	got := GenerateRust(fn)

	expectedPieces := []string{
		"#[no_mangle]",
		"pub extern \"C\" fn add",
		"(x: i64, y: i64) -> i64",
		"x + y",
	}

	for _, piece := range expectedPieces {
		if !strings.Contains(got, piece) {
			t.Errorf("expected generated code to contain %q, but it didn't.\nGot:\n%s", piece, got)
		}
	}
}

func TestWriteRustFile(t *testing.T) {
	// Simple test to ensure WriteRustFile doesn't crash and creates the file
	code := `#[no_mangle]
pub extern "C" fn test_func() -> i64 { 42 }`

	// Use a temporary file for testing
	tmpFile := "/tmp/rubo_test_hot_func.rs"
	err := WriteRustFile(tmpFile, code)
	if err != nil {
		t.Fatalf("WriteRustFile failed: %v", err)
	}
}

func TestCompileRust(t *testing.T) {
	if _, err := exec.LookPath("cargo"); err != nil {
		t.Skip("cargo not found in PATH, skipping compilation test")
	}

	err := CompileRust()
	if err != nil {
		t.Fatalf("CompileRust failed: %v", err)
	}
}
