package lzfse

/*
#include "lzfse.h"
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"io/ioutil"
	"sync"
	"testing"
	"unsafe"
)

type cgoAllocMap struct {
	mux sync.RWMutex
	m   map[unsafe.Pointer]struct{}
}

var cgoAllocsUnknown = new(cgoAllocMap)

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

func unpackPUint8String(str string) (*C.uint8_t, *cgoAllocMap) {
	h := (*stringHeader)(unsafe.Pointer(&str))
	return (*C.uint8_t)(h.Data), cgoAllocsUnknown
}

// EncodeScratchSize function as declared in lzfse.h:56
func EncodeScratchSize() uint {
	__ret := C.lzfse_encode_scratch_size()
	__v := (uint)(__ret)
	return __v
}

// EncodeBuffer function as declared in lzfse.h:87
func EncodeBuffer(dstBuffer []byte, dstSize uint, srcBuffer string, srcSize uint, scratchBuffer unsafe.Pointer) uint {
	cdstBuffer, _ := (*C.uint8_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&dstBuffer)).Data)), cgoAllocsUnknown
	cdstSize, _ := (C.size_t)(dstSize), cgoAllocsUnknown
	csrcBuffer, _ := unpackPUint8String(srcBuffer)
	csrcSize, _ := (C.size_t)(srcSize), cgoAllocsUnknown
	cscratchBuffer, _ := scratchBuffer, cgoAllocsUnknown
	__ret := C.lzfse_encode_buffer(cdstBuffer, cdstSize, csrcBuffer, csrcSize, cscratchBuffer)
	__v := (uint)(__ret)
	return __v
}

// DecodeScratchSize function as declared in lzfse.h:94
func DecodeScratchSize() uint {
	__ret := C.lzfse_decode_scratch_size()
	__v := (uint)(__ret)
	return __v
}

// DecodeBuffer function as declared in lzfse.h:126
func DecodeBuffer(srcBuffer []byte) []byte {
	compRatio := 4
	in, _ := unpackPUint8String(string(srcBuffer))
	in_size, _ := (C.size_t)(len(srcBuffer)), cgoAllocsUnknown

	dstBuffer := make([]byte, compRatio*len(srcBuffer))
	out, _ := (*C.uint8_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&dstBuffer)).Data)), cgoAllocsUnknown
	out_allocated, _ := (C.size_t)(compRatio*len(srcBuffer)), cgoAllocsUnknown

	scratch := make([]byte, DecodeScratchSize())
	aux, _ := unsafe.Pointer(&scratch[0]), cgoAllocsUnknown

	for {
		__ret := C.lzfse_decode_buffer(out, out_allocated, in, in_size, aux)
		out_size := (C.size_t)(__ret)
		// If output buffer was too small, grow and retry.
		if out_size == 0 || out_size == out_allocated {
			compRatio *= 2
			dstBuffer = make([]byte, compRatio*len(srcBuffer))
			out, _ = (*C.uint8_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&dstBuffer)).Data)), cgoAllocsUnknown
			out_allocated, _ = (C.size_t)(compRatio*len(srcBuffer)), cgoAllocsUnknown
		} else {
			return dstBuffer[:out_size]
		}
	}
}

func testDecodeBuffer(t *testing.T, encBuf, wantBuf []byte) {
	t.Run("README", func(t *testing.T) {
		if got := DecodeBuffer(encBuf); !bytes.Contains(got, wantBuf) {
			ioutil.WriteFile("fail.bin", got, 0755)
			t.Errorf("DecodeBuffer() = %v, want %v", got, wantBuf)
		}
	})
}
