package lzfse

import (
	"C"
)
import (
	"reflect"
	"testing"
)

func TestDecodeBuffer(t *testing.T) {
	type args struct {
		srcBuffer []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeBuffer(tt.args.srcBuffer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}
