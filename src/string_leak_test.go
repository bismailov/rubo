package main

import (
	"fmt"
	"os"
	"path/filepath"
	"rubo-lang/src/bridge"
	"runtime"
	"testing"
)

// Run with: go test -v src/string_leak_test.go
func TestStringBridgeLeakage(t *testing.T) {
	cwd, _ := os.Getwd()
	// Check if we are in src or root
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

	fn, err := bridge.LoadRustStringFunction(libPath, "rubo_string_len")
	if err != nil {
		t.Fatalf("Failed to load string function: %v", err)
	}

	const iterations = 1000000
	testStr := "Hello, Rubo String Bridge!"

	fmt.Printf("--- Starting String Bridge Leakage Test: %d iterations ---\n", iterations)

	for i := 0; i < iterations; i++ {
		length, _ := fn(testStr)
		if i == 0 && length != int32(len(testStr)) {
			t.Errorf("Expected length %d, got %d", len(testStr), length)
		}
		if i%200000 == 0 {
			printMemUsage(i)
		}
	}
	printMemUsage(iterations)
	fmt.Printf("--- Test Complete ---\n")
}

func printMemUsage(iter int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Iter %-7d | Alloc = %-4v MiB | TotalAlloc = %-4v MiB | NumGC = %-3v\n",
		iter, bToMb(m.Alloc), bToMb(m.TotalAlloc), m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
