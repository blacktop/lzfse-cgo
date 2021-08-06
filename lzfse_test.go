package lzfse

import (
	"io/ioutil"
	"testing"
)

func TestDecodeBuffer(t *testing.T) {
	wantBuf, err := ioutil.ReadFile("README.md")
	if err != nil {
		t.Errorf("failed to read test file 'README.md': %v", err)
	}
	encBuff, err := ioutil.ReadFile("test/enc.bin")
	if err != nil {
		t.Errorf("failed to read test file 'test/enc.bin': %v", err)
	}
	testDecodeBuffer(t, encBuff, wantBuf)
}
