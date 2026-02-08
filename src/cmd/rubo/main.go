package main

/*
#cgo LDFLAGS: -L../../runtime/target/debug -lruntime -lpthread -ldl
#include <stdint.h>

int32_t rubo_add(int32_t a, int32_t b);
*/
import "C"
import "fmt"

func main() {
	a := int32(10)
	b := int32(32)
	result := C.rubo_add(C.int32_t(a), C.int32_t(b))
	fmt.Printf("Rubo Core PoC\n")
	fmt.Printf("%d + %d = %d (calculated by Rust)\n", a, b, int32(result))
}
