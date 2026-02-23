package bridge

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdint.h>
#include <stdlib.h>

typedef int64_t (*rust_fn_2)(int64_t, int64_t);
typedef int32_t (*rust_fn_string)(const char*);

int64_t call_rust_fn_2(void* f, int64_t a, int64_t b) {
    return ((rust_fn_2)f)(a, b);
}

int32_t call_rust_fn_string(void* f, const char* s) {
    return ((rust_fn_string)f)(s);
}
*/
import "C"
import (
	"errors"
	"fmt"
	"rubo-lang/internal/arena"
	"unsafe"
)

// LoadRustFunction loads a function from a shared library.
// For now, it specifically supports functions with 2 int64 arguments.
func LoadRustFunction(libPath string, funcName string) (func(...int64) int64, error) {
	cLibPath := C.CString(libPath)
	defer C.free(unsafe.Pointer(cLibPath))

	handle := C.dlopen(cLibPath, C.RTLD_NOW)
	if handle == nil {
		return nil, fmt.Errorf("failed to open library: %s", C.GoString(C.dlerror()))
	}

	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))

	symbol := C.dlsym(handle, cFuncName)
	if symbol == nil {
		return nil, fmt.Errorf("failed to find symbol: %s", C.GoString(C.dlerror()))
	}

	return func(args ...int64) int64 {
		if len(args) == 2 {
			return int64(C.call_rust_fn_2(symbol, C.int64_t(args[0]), C.int64_t(args[1])))
		}
		// Fallback for other arities (to be implemented as needed)
		if len(args) == 0 {
			return 0
		}
		return 0
	}, nil
}

// LoadRustStringFunction loads a function that takes a string and returns an int32.
func LoadRustStringFunction(libPath string, funcName string) (func(string) (int32, error), error) {
	cLibPath := C.CString(libPath)
	defer C.free(unsafe.Pointer(cLibPath))

	handle := C.dlopen(cLibPath, C.RTLD_NOW)
	if handle == nil {
		return nil, fmt.Errorf("failed to open library: %s", C.GoString(C.dlerror()))
	}

	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))

	symbol := C.dlsym(handle, cFuncName)
	if symbol == nil {
		return nil, fmt.Errorf("failed to find symbol: %s", C.GoString(C.dlerror()))
	}

	// Phase 9: Use Arena for string passing to avoid malloc overhead.
	// We allocate a 2MB buffer once per function load.
	// This assumes single-threaded access to the returned function.
	a := arena.NewArena(2 * 1024 * 1024)

	return func(input string) (int32, error) {
		// Reset the arena at the start of the execution cycle (per call)
		a.Reset()

		// Copy string to arena (no malloc)
		cStr := (*C.char)(a.Copy(input))
		// No need to free cStr, it lives in the arena which is reset next time.

		res := int32(C.call_rust_fn_string(symbol, cStr))

		switch res {
		case -1:
			return 0, errors.New("RuboRuntimeError: Received null pointer")
		case -2:
			return 0, errors.New("RuboTypeError: Invalid UTF-8 sequence")
		case -3:
			return 0, errors.New("RuboPanic: The Rust runtime encountered an unrecoverable error")
		default:
			return res, nil
		}
	}, nil
}
