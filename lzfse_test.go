package lzfse

import (
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
