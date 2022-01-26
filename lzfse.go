package lzfse

/*
#cgo CFLAGS: -I${SRCDIR}

#include <stdlib.h>
#include "lzfse.h"
#include "lzvn_decode_base.h"
#include "lzvn_encode_base.h"

void lzvn_decode(lzvn_decoder_state *state);
void lzvn_encode(lzvn_encoder_state *state);

lzvn_decoder_state *lzvn_decoder_state_init(uint8_t * dst, size_t dst_size, const uint8_t * src, size_t src_size);
lzvn_encoder_state *lzvn_encoder_state_init(uint8_t * dst, size_t dst_size, const uint8_t * src, size_t src_size);

lzvn_decoder_state *lzvn_decoder_state_init(uint8_t * dst, size_t dst_size, const uint8_t * src, size_t src_size) {
	lzvn_decoder_state *state = malloc(sizeof(lzvn_decoder_state));
	state->src = src;
	state->src_end = src + src_size;
	state->dst = dst;
	state->dst_begin = dst;
	state->dst_end = dst + dst_size;
	state->dst_current = dst;
	return state;
}

lzvn_encoder_state *lzvn_encoder_state_init(uint8_t * dst, size_t dst_size, const uint8_t * src, size_t src_size) {
	// Max input size check (limit to offsets on uint32_t).
	if (src_size > LZVN_ENCODE_MAX_SRC_SIZE) {
	src_size = LZVN_ENCODE_MAX_SRC_SIZE;
	}

	lzvn_encoder_state *state = malloc(sizeof(lzvn_encoder_state));
  	memset(&state, 0, sizeof(state));

	void *__restrict scratch_buffer = malloc(lzfse_encode_scratch_size() + 1);

	state->src = src;
	state->src_begin = 0;
	state->src_end = (lzvn_offset)src_size;
	state->src_literal = 0;
	state->src_current = 0;
	state->dst = dst;
	state->dst_begin = dst;
	state->dst_end = (unsigned char *)dst + dst_size - 8; // reserve 8 bytes for end-of-stream
	state->table = scratch_buffer;

	return state;
}
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
func EncodeLZVNBuffer(srcBuf []byte) []byte {
	decBuf := make([]byte, len(srcBuf)*10)

	state := C.lzvn_encoder_state_init(
		(*C.uint8_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&decBuf)).Data)),
		(C.size_t)(len(decBuf)),
		(*C.uint8_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&srcBuf)).Data)),
		(C.size_t)(len(srcBuf)),
	)
	defer C.free(unsafe.Pointer(state.table))
	defer C.free(unsafe.Pointer(state))

	C.lzvn_encode(state)

	dstSize := (C.size_t)(*state.dst - *state.dst_begin)

	return decBuf[:dstSize]
}

// DecodeLZVNBuffer function as declared in lzvn_decode_base.c:47
func DecodeLZVNBuffer(encBuf []byte, uncompressedSize uint64) []byte {
	decBuf := make([]byte, uncompressedSize)

	state := C.lzvn_decoder_state_init(
		(*C.uint8_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&decBuf)).Data)),
		(C.size_t)(len(decBuf)),
		(*C.uint8_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&encBuf)).Data)),
		(C.size_t)(len(encBuf)),
	)
	defer C.free(unsafe.Pointer(state))

	C.lzvn_decode(state)

	return decBuf
}

func testDecodeBuffer(t *testing.T, encBuf, wantBuf []byte) {
	t.Run("README", func(t *testing.T) {
		if got := DecodeBuffer(encBuf); !bytes.Contains(got, wantBuf) {
			ioutil.WriteFile("fail.out", got, 0755)
			t.Errorf("DecodeBuffer() = %v, want %v", got, wantBuf)
		}
	})
}

func testDecodeLZVNBuffer(t *testing.T, encBuf, wantBuf []byte) {
	t.Run("test/lzvn_enc.bin", func(t *testing.T) {
		if got := DecodeLZVNBuffer(encBuf, 68608); !bytes.Contains(got, wantBuf) {
			if err := ioutil.WriteFile("fail.out", got, 0755); err != nil {
				t.Errorf("failed to write fail.out: %v", err)
			}
			t.Errorf("DecodeLZVNBuffer() = %v, want %v", got, wantBuf)
		}
	})
}

func testEncodeLZVNBuffer(t *testing.T, srcBuf, wantBuf []byte) {
	t.Run("test/lzvn_dec.bin", func(t *testing.T) {
		if got := EncodeLZVNBuffer(srcBuf); !bytes.Contains(got, wantBuf) {
			if err := ioutil.WriteFile("fail.out", got, 0755); err != nil {
				t.Errorf("failed to write fail.out: %v", err)
			}
			t.Errorf("EncodeLZVNBuffer() = %v, want %v", got, wantBuf)
		}
	})
}
