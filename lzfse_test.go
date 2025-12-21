package lzfse

import (
	"bytes"
	"os"
	"testing"
)

func TestDecodeBuffer(t *testing.T) {
	wantBuf, err := os.ReadFile("test/dec.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/dec.bin': %v", err)
	}
	encBuff, err := os.ReadFile("test/enc.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/enc.bin': %v", err)
	}
	testDecodeBuffer(t, encBuff, wantBuf)
}

func TestDecodeBufferInto(t *testing.T) {
	wantBuf, err := os.ReadFile("test/dec.bin")
	if err != nil {
		t.Fatalf("failed to read test file 'test/dec.bin': %v", err)
	}
	encBuff, err := os.ReadFile("test/enc.bin")
	if err != nil {
		t.Fatalf("failed to read test file 'test/enc.bin': %v", err)
	}

	t.Run("DecodeBufferInto", func(t *testing.T) {
		dst := make([]byte, len(wantBuf))
		n := DecodeBufferInto(encBuff, dst)
		if n == 0 {
			t.Fatal("DecodeBufferInto returned 0")
		}
		if !bytes.Equal(dst[:n], wantBuf) {
			t.Errorf("DecodeBufferInto() result mismatch")
		}
	})

	t.Run("DecodeBufferWithScratch", func(t *testing.T) {
		dst := make([]byte, len(wantBuf))
		scratch := make([]byte, DecodeScratchSize())
		n := DecodeBufferWithScratch(encBuff, dst, scratch)
		if n == 0 {
			t.Fatal("DecodeBufferWithScratch returned 0")
		}
		if !bytes.Equal(dst[:n], wantBuf) {
			t.Errorf("DecodeBufferWithScratch() result mismatch")
		}
	})

	t.Run("empty_src", func(t *testing.T) {
		dst := make([]byte, 100)
		if n := DecodeBufferInto(nil, dst); n != 0 {
			t.Errorf("expected 0 for nil src, got %d", n)
		}
		if n := DecodeBufferInto([]byte{}, dst); n != 0 {
			t.Errorf("expected 0 for empty src, got %d", n)
		}
	})

	t.Run("empty_dst", func(t *testing.T) {
		if n := DecodeBufferInto(encBuff, nil); n != 0 {
			t.Errorf("expected 0 for nil dst, got %d", n)
		}
		if n := DecodeBufferInto(encBuff, []byte{}); n != 0 {
			t.Errorf("expected 0 for empty dst, got %d", n)
		}
	})
}

func TestEncodeBuffer(t *testing.T) {
	wantBuf, err := os.ReadFile("test/enc.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/enc.bin': %v", err)
	}
	decBuff, err := os.ReadFile("test/dec.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/dec.bin': %v", err)
	}
	testEncodeBuffer(t, decBuff, wantBuf)
}

func TestDecodeLZVNBuffer(t *testing.T) {
	wantBuf, err := os.ReadFile("test/lzvn_dec.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/lzvn_dec.bin': %v", err)
	}
	encBuff, err := os.ReadFile("test/lzvn_enc.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/lzvn_enc.bin': %v", err)
	}
	testDecodeLZVNBuffer(t, encBuff, wantBuf)
}

func TestEncodeLZVNBuffer(t *testing.T) {
	wantBuf, err := os.ReadFile("test/lzvn_enc.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/lzvn_enc.bin': %v", err)
	}
	srcBuff, err := os.ReadFile("test/lzvn_dec.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/lzvn_dec.bin': %v", err)
	}
	testEncodeLZVNBuffer(t, srcBuff, wantBuf)
}
