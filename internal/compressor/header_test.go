package compressor

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/chavacava/next/internal/types"
)

func TestWriteHeader(t *testing.T) {
	tt := map[string]struct {
		root        byte
		size        types.Size
		symbolCount byte
		want        []byte
	}{
		"empty content": {
			root:        0,
			size:        0,
			symbolCount: 0,
			want:        []byte{137, 78, 69, 88, 84, 13, 10, 26, 10, 0, 23, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 26},
		},
		"1 byte content length": {
			root:        0,
			size:        1,
			symbolCount: 1,
			want:        []byte{137, 78, 69, 88, 84, 13, 10, 26, 10, 0, 23, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 28},
		},
		"1000 bytes content length root 255": {
			root:        255,
			size:        1000,
			symbolCount: 3,
			want:        []byte{137, 78, 69, 88, 84, 13, 10, 26, 10, 0, 23, 0, 255, 232, 3, 0, 0, 0, 0, 0, 0, 3, 7},
		},
	}

	for name, tc := range tt {
		t.Run(name,
			func(t *testing.T) {
				got := new(bytes.Buffer)
				WriteHeader(got, tc.root, tc.size, tc.symbolCount)
				if !reflect.DeepEqual(tc.want, got.Bytes()) {
					t.Fatalf("expected\n\t%v\ngot\n\t%v", tc.want, got.Bytes())
				}
			},
		)
	}

}
