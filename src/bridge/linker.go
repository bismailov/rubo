package bridge

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdint.h>
#include <stdlib.h>

typedef int64_t (*rust_fn_2)(int64_t, int64_t);

int64_t call_rust_fn_2(void* f, int64_t a, int64_t b) {
    return ((rust_fn_2)f)(a, b);
}
*/
import "C"
import (
	"fmt"
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
