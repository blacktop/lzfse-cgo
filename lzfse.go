package lzfse

/*
#cgo CFLAGS: -I${SRCDIR}

#include <stdlib.h>
#include "lzfse.h"
#include "lzvn_decode_base.h"
#include "lzvn_encode_base.h"

size_t lzvn_decode_buffer(void *__restrict dst, size_t dst_size,
                          const void *__restrict src, size_t src_size);

size_t lzvn_decode_buffer(void *__restrict dst, size_t dst_size,
                          const void *__restrict src, size_t src_size) {
  // Init LZVN decoder state
  lzvn_decoder_state dstate;
  memset(&dstate, 0x00, sizeof(dstate));
  dstate.src = src;
  dstate.src_end = (const unsigned char*) src + src_size;

  dstate.dst_begin = dst;
  dstate.dst = dst;
  dstate.dst_end = (unsigned char*) dst + dst_size;

  dstate.d_prev = 0;
  dstate.end_of_stream = 0;

  // Run LZVN decoder
  lzvn_decode(&dstate);

  // This is how much we decompressed
  return dstate.dst - (unsigned char*) dst;
}
*/
import "C"
import (
	"bytes"
	"os"
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
func EncodeBuffer(srcBuffer []byte) []byte {
	csrcBuffer, _ := unpackPUint8String(string(srcBuffer))
	csrcSize, _ := (C.size_t)(len(srcBuffer)), cgoAllocsUnknown
	dstBuffer := make([]byte, len(srcBuffer)*2)
	cdstBuffer, _ := (*C.uint8_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&dstBuffer)).Data)), cgoAllocsUnknown
	cdstSize, _ := (C.size_t)(len(dstBuffer)), cgoAllocsUnknown
	__ret := C.lzfse_encode_buffer(cdstBuffer, cdstSize, csrcBuffer, csrcSize, nil)
	out_size := (C.size_t)(__ret)
	return dstBuffer[:out_size]
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
	size := 0

	for i := 0; i < 16 && size < 50_000_000; i++ {
		__ret := C.lzfse_decode_buffer(out, out_allocated, in, in_size, aux)
		out_size := (C.size_t)(__ret)
		// If output buffer was too small, grow and retry.
		if out_size == 0 || out_size == out_allocated {
			compRatio *= 2
			size = compRatio * len(srcBuffer)
			dstBuffer = make([]byte, size)
			out, _ = (*C.uint8_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&dstBuffer)).Data)), cgoAllocsUnknown
			out_allocated, _ = (C.size_t)(compRatio*len(srcBuffer)), cgoAllocsUnknown
		} else {
			return dstBuffer[:out_size]
		}
	}

	return dstBuffer[:0]
}

// EncodeLZVNBuffer function as declared in lzvn_encode_base.c:383
func EncodeLZVNBuffer(srcBuf, dstBuf []byte) uint {
	scratch := make([]byte, DecodeScratchSize())
	__ret := C.lzvn_encode_buffer(
		unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&dstBuf)).Data),
		(C.size_t)(len(dstBuf)),
		unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&srcBuf)).Data),
		(C.size_t)(len(srcBuf)),
		unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&scratch)).Data),
	)
	__v := (uint)(__ret)
	return __v
}

// DecodeLZVNBuffer function as declared in lzfse_internal.h:413
func DecodeLZVNBuffer(encBuf, decBuf []byte) uint {
	__ret := C.lzvn_decode_buffer(
		unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&decBuf)).Data),
		(C.size_t)(len(decBuf)),
		unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&encBuf)).Data),
		(C.size_t)(len(encBuf)),
	)
	__v := (uint)(__ret)
	return __v
}

func testDecodeBuffer(t *testing.T, encBuf, wantBuf []byte) {
	t.Run("README", func(t *testing.T) {
		if got := DecodeBuffer(encBuf); !bytes.Equal(got, wantBuf) {
			os.WriteFile("fail.out", got, 0755)
			t.Errorf("DecodeBuffer() = %v, want %v", got, wantBuf)
		}
	})
}

func testEncodeBuffer(t *testing.T, encBuf, wantBuf []byte) {
	t.Run("README", func(t *testing.T) {
		if got := EncodeBuffer(encBuf); !bytes.Equal(got, wantBuf) {
			os.WriteFile("fail.out", got, 0755)
			t.Errorf("DecodeBuffer() = %v, want %v", got, wantBuf)
		}
	})
}

func testDecodeLZVNBuffer(t *testing.T, encBuf, wantBuf []byte) {
	t.Run("test/lzvn_enc.bin", func(t *testing.T) {
		got := make([]byte, 68608)
		if DecodeLZVNBuffer(encBuf, got); !bytes.Contains(got, wantBuf) {
			if err := os.WriteFile("fail.out", got, 0755); err != nil {
				t.Errorf("failed to write fail.out: %v", err)
			}
			t.Errorf("DecodeLZVNBuffer() = %v, want %v", got, wantBuf)
		}
	})
}

func testEncodeLZVNBuffer(t *testing.T, srcBuf, wantBuf []byte) {
	t.Run("test/lzvn_dec.bin", func(t *testing.T) {
		got := make([]byte, len(srcBuf)*4)
		if EncodeLZVNBuffer(srcBuf, got); !bytes.Contains(got, wantBuf) {
			if err := os.WriteFile("fail.out", got, 0755); err != nil {
				t.Errorf("failed to write fail.out: %v", err)
			}
			t.Errorf("EncodeLZVNBuffer() = %v, want %v", got, wantBuf)
		}
	})
}
