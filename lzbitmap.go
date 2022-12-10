package lzfse

/*
#cgo CFLAGS: -I${SRCDIR}

#include <stdlib.h>
#include "libzbitmap.h"

int zbm_decompress(void *dest, size_t dest_size, const void *src, size_t src_size, size_t *out_len);
*/
import "C"
import (
	"unsafe"
)

func LzBitMapDecompress(src, dst []byte) int {
	var outLen C.size_t
	return int(C.zbm_decompress(
		unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&dst)).Data),
		(C.size_t)(len(dst)),
		unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&src)).Data),
		(C.size_t)(len(src)),
		(*C.size_t)(unsafe.Pointer(&outLen)),
	))
}
