package lzfse

import (
	"io/ioutil"
	"testing"
)

func TestDecodeBuffer(t *testing.T) {
	wantBuf, err := ioutil.ReadFile("test/dec.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/dec.bin': %v", err)
	}
	encBuff, err := ioutil.ReadFile("test/enc.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/enc.bin': %v", err)
	}
	testDecodeBuffer(t, encBuff, wantBuf)
}
