package huffman

import (
	"fmt"
	"testing"

	"github.com/chavacava/next/internal/bitstream"
)

func TestNewTree(t *testing.T) {
	tt := map[string]struct {
		fs   []SymbolFreq
		want string
	}{
		"1 element": {
			fs: []SymbolFreq{
				SymbolFreq{byte(65), 1},
			},
			want: "[ S: 65 ]",
		},
		"3 elements": {
			fs: []SymbolFreq{
				SymbolFreq{byte(65), 1},
				SymbolFreq{byte(66), 2},
				SymbolFreq{byte(67), 2},
			},
			want: "{ l: [ S: 67 ] r: { l: [ S: 65 ] r: [ S: 66 ] } }",
		},
	}
	for _, tc := range tt {

		got := NewTree(tc.fs).root.String()

		if got != tc.want {
			t.Fatalf("expected\n\t%v\ngot\n\t%v", tc.want, got)
		}
	}
}

func TestDictionary(t *testing.T) {
	tt := map[string]struct {
		fs   []SymbolFreq
		want string
	}{
		"3 elements": {
			fs: []SymbolFreq{
				SymbolFreq{byte(65), 1},
				SymbolFreq{byte(66), 2},
				SymbolFreq{byte(67), 3},
			},
			want: "map[65:{[false false] 0} 66:{[false true] 0} 67:{[true] 0}]",
		},
	}
	for _, tc := range tt {
		tree := NewTree(tc.fs)
		got := tree.Dictionary()
		if fmt.Sprintf("%v", got) != tc.want {
			t.Fatalf("expected dictionary to be\n\t%v\ngot\n\t%v", tc.want, got)
		}
	}
}

func TestInterpret(t *testing.T) {
	fs := []SymbolFreq{
		SymbolFreq{byte(65), 1},
		SymbolFreq{byte(66), 2},
		SymbolFreq{byte(67), 3},
	}
	tree := NewTree(fs)

	tt := map[string]struct {
		bs         bitstream.BitStream
		wantSymbol byte
	}{
		"symbol 67": {
			bs:         bitstream.NewFromBits([]bitstream.Bit{true, true, true, true, true, true}),
			wantSymbol: byte(67),
		},
		"symbol 67 again": {
			bs:         bitstream.NewFromBits([]bitstream.Bit{true, true, false, true}),
			wantSymbol: byte(67),
			//"map[65:[false false] 66:[false true] 67:[true]]",
		},
		"symbol 66": {
			bs:         bitstream.NewFromBits([]bitstream.Bit{false, true, true, true, true, true}),
			wantSymbol: byte(66),
			//"map[65:[false false] 66:[false true] 67:[true]]",
		},
		"symbol 65": {
			bs:         bitstream.NewFromBits([]bitstream.Bit{false, false, true, true}),
			wantSymbol: byte(65),
			//"map[65:[false false] 66:[false true] 67:[true]]",
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			gotSymbol := tree.Interpret(tc.bs)

			if gotSymbol != tc.wantSymbol {
				t.Fatalf("expected symbol %v, got %v", tc.wantSymbol, gotSymbol)
			}
		})
	}
}

func TestAsBitstream(t *testing.T) {
	fs := []SymbolFreq{
		SymbolFreq{byte(65), 1},
		SymbolFreq{byte(66), 2},
		SymbolFreq{byte(67), 3},
	}

	tree := NewTree(fs)
	want := bitstream.NewFromBits([]bitstream.Bit{false, false, true})
	want.Append(bitstream.NewFromFullByte(65))
	want.Append(bitstream.NewFromBits([]bitstream.Bit{true}))
	want.Append(bitstream.NewFromFullByte(66))
	want.Append(bitstream.NewFromBits([]bitstream.Bit{true}))
	want.Append(bitstream.NewFromFullByte(67))

	got := tree.AsBitstream()

	if !want.IsEqual(got) {
		t.Fatalf("expected\n\t%v\ngot\n\t%v", want, got)
	}
}

func TestNewTreeFromBS(t *testing.T) {
	fs := []SymbolFreq{
		SymbolFreq{byte(65), 1},
		SymbolFreq{byte(66), 2},
		SymbolFreq{byte(67), 3},
	}

	want := NewTree(fs)
	bs := bitstream.NewFromBits([]bitstream.Bit{false, false, true})
	bs.Append(bitstream.NewFromFullByte(65))
	bs.Append(bitstream.NewFromBits([]bitstream.Bit{true}))
	bs.Append(bitstream.NewFromFullByte(66))
	bs.Append(bitstream.NewFromBits([]bitstream.Bit{true}))
	bs.Append(bitstream.NewFromFullByte(67))

	got := NewTreeFromBS(bs)

	if want.String() != got.String() {
		t.Fatalf("expected\n\t%v\ngot\n\t%v", want, got)
	}
}
