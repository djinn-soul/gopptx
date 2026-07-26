package main

/*
#include <stdlib.h>
*/
import "C"

import "unsafe"

//export deck_global_error
func deck_global_error() *C.char {
	globalErrorMu.RLock()
	defer globalErrorMu.RUnlock()
	if globalError == "" {
		return nil
	}
	return C.CString(globalError)
}

//export deck_free_string
func deck_free_string(value *C.char) {
	if value != nil {
		C.free(unsafe.Pointer(value))
	}
}
