package bitstream

import (
	"errors"
	"reflect"
	"testing"

	"github.com/chavacava/next/internal/types"
)

func TestIsEqual(t *testing.T) {
	tt := map[string]struct {
		this  BitStream
		other BitStream
		want  bool
	}{
		"emptys": {
			this:  New(),
			other: New(),
			want:  true,
		},
		"empty vs not empty": {
			this:  New(),
			other: NewFromBits([]Bit{false}),
			want:  false,
		},
		"not empty vs empty": {
			this:  NewFromBits([]Bit{true}),
			other: NewFromBits([]Bit{}),
			want:  false,
		},
		"not exactly the same": {
			this:  NewFromBits([]Bit{false, false, false, false, true, false, false, true}),
			other: NewFromBits([]Bit{false, false, false, false, true, true, false, true}),
			want:  false,
		},
		"long the same": {
			this:  NewFromBits([]Bit{true, false, false, false, true, false, false, true, false, false, false, true, false, false, true, false, false, false, true, false, false, true}),
			other: NewFromBits([]Bit{true, false, false, false, true, false, false, true, false, false, false, true, false, false, true, false, false, false, true, false, false, true}),
			want:  true,
		},
	}

	for name, tc := range tt {
		t.Run(name,
			func(t *testing.T) {
				got := tc.this.IsEqual(tc.other)
				if tc.want != got {
					t.Fatalf("for\n\t%v and\n\t%v\nexpected %v, got %v", tc.this, tc.other, tc.want, got)
				}
			},
		)
	}
}

func TestNewFromBits(t *testing.T) {
	tt := map[string]struct {
		bits []Bit
		want BitStream
	}{
		"empty": {
			bits: []Bit{},
			want: BitStream{bits: []Bit{}, idx: 0},
		},
		"one false bit": {
			bits: []Bit{false},
			want: BitStream{bits: []Bit{false}, idx: 0},
		},
		"one true bit": {
			bits: []Bit{true},
			want: BitStream{bits: []Bit{true}, idx: 0},
		},
		"many bits": {
			bits: []Bit{true, false, true, false, true, false, true, false, true, false, false},
			want: BitStream{bits: []Bit{true, false, true, false, true, false, true, false, true, false, false}, idx: 0},
		},
	}

	for name, tc := range tt {
		t.Run(name,
			func(t *testing.T) {
				got := NewFromBits(tc.bits)
				if !tc.want.IsEqual(got) {
					t.Fatalf("expected %v, got %v", tc.want, got)
				}
			},
		)
	}
}

func TestAppend(t *testing.T) {
	tt := map[string]struct {
		this  BitStream
		other BitStream
		want  BitStream
	}{
		"empty": {
			this:  BitStream{bits: []Bit{}, idx: 0},
			other: BitStream{bits: []Bit{}, idx: 0},
			want:  BitStream{bits: []Bit{}, idx: 0},
		},
		"empty + non empty": {
			this:  BitStream{bits: []Bit{}, idx: 0},
			other: BitStream{bits: []Bit{true}, idx: 1},
			want:  BitStream{bits: []Bit{true}, idx: 0},
		},
		"non empty + non empty": {
			this:  BitStream{bits: []Bit{true, true, true}, idx: 0},
			other: BitStream{bits: []Bit{false, false}, idx: 1},
			want:  BitStream{bits: []Bit{true, true, true, false, false}, idx: 0},
		},
	}

	for name, tc := range tt {
		t.Run(name,
			func(t *testing.T) {
				tc.this.Append(tc.other)
				if !tc.want.IsEqual(tc.this) {
					t.Fatalf("expected\n\t%v\ngot\n\t%v", tc.want, tc.this)
				}
			},
		)
	}
}

func TestToBytes(t *testing.T) {
	tt := map[string]struct {
		bs   BitStream
		want []byte
	}{
		"emptys": {
			bs:   NewFromBits([]Bit{}),
			want: []byte{},
		},
		"0 one bit long": {
			bs:   NewFromBits([]Bit{false}),
			want: []byte{0},
		},
		"1 one bit long": {
			bs:   NewFromBits([]Bit{true}),
			want: []byte{128},
		},
		"seven bits long": {
			bs:   NewFromBits([]Bit{true, true, true, true, true, true, false}),
			want: []byte{252},
		},
		"eight bits long": {
			bs:   NewFromBits([]Bit{false, true, true, true, true, true, true, true}),
			want: []byte{127},
		},
		"nine bits long": {
			bs: NewFromBits([]Bit{false, true, true, true, true, true, true, true,
				false}),
			want: []byte{127, 0},
		},
		"ten bits long": {
			bs: NewFromBits([]Bit{false, true, true, true, true, true, true, true,
				true, true}),
			want: []byte{127, 192},
		},
		"nineteen bits long": {
			bs: NewFromBits([]Bit{false, false, false, false, false, false, false, true,
				false, false, false, false, false, false, true, false,
				false, true, true}),
			want: []byte{1, 2, 96},
		},
	}

	for name, tc := range tt {
		t.Run(name,
			func(t *testing.T) {
				got := tc.bs.Bytes()
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("for\n\t%v \nexpected\n\t%v\ngot\n\t%v", tc.bs, tc.want, got)
				}
			},
		)
	}
}

func TestNewFromBytes(t *testing.T) {
	tt := map[string]struct {
		bs   []byte
		want BitStream
	}{
		"zero bytes": {
			bs:   []byte{},
			want: NewFromBits([]Bit{}),
		},
		"one byte (0)": {
			bs:   []byte{0},
			want: NewFromBits([]Bit{false, false, false, false, false, false, false, false}),
		},
		"one byte (25)": {
			bs:   []byte{25},
			want: NewFromBits([]Bit{false, false, false, true, true, false, false, true}),
		},
		"two bytes (25,1)": {
			bs: []byte{25, 1},
			want: NewFromBits([]Bit{
				false, false, false, true, true, false, false, true,
				false, false, false, false, false, false, false, true,
			}),
		},
		"ten bytes (1,2,3,4,5,6,7,8,9,10)": {
			bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want: NewFromBits([]Bit{
				false, false, false, false, false, false, false, true,
				false, false, false, false, false, false, true, false,
				false, false, false, false, false, false, true, true,
				false, false, false, false, false, true, false, false,
				false, false, false, false, false, true, false, true,
				false, false, false, false, false, true, true, false,
				false, false, false, false, false, true, true, true,
				false, false, false, false, true, false, false, false,
				false, false, false, false, true, false, false, true,
				false, false, false, false, true, false, true, false,
			}),
		},
	}

	for name, tc := range tt {
		t.Run(name,
			func(t *testing.T) {
				got := newFromBytes(tc.bs)
				if !tc.want.IsEqual(got) {
					t.Fatalf("expected\n\t%v\ngot\n\t%v", tc.want, got)
				}
			},
		)
	}
}

func TestNewFromByte(t *testing.T) {
	tt := map[string]struct {
		b    byte
		size byte
		want BitStream
	}{
		"0 size 1": {
			b:    0,
			size: 1,
			want: NewFromBits([]Bit{false}),
		},
		"1 size 2": {
			b:    1,
			size: 2,
			want: NewFromBits([]Bit{false, true}),
		},
		"9 size 8": {
			b:    9,
			size: 8,
			want: NewFromBits([]Bit{false, false, false, false, true, false, false, true}),
		},
	}

	for name, tc := range tt {
		t.Run(name,
			func(t *testing.T) {
				got := NewFromByte(tc.b, tc.size)
				if !tc.want.IsEqual(got) {
					t.Fatalf("expected %v, got %v", tc.want, got)
				}
			},
		)
	}
}

func TestRead(t *testing.T) {
	tt := map[string]struct {
		bs      BitStream
		idx     types.Position
		want    bool
		wantErr error
	}{
		"read on empty bs": {
			bs:      NewFromBits([]Bit{}),
			idx:     0,
			want:    false,
			wantErr: errors.New(""),
		},
		"read first bit": {
			bs:      NewFromBits([]Bit{true}),
			idx:     0,
			want:    true,
			wantErr: nil,
		},
		"read third bit": {
			bs:      NewFromBits([]Bit{true, true, false, true, true}),
			idx:     2,
			want:    false,
			wantErr: nil,
		},
	}

	for name, tc := range tt {
		t.Run(name,
			func(t *testing.T) {
				tc.bs.idx = tc.idx
				got, err := tc.bs.Read()

				if tc.wantErr != nil && err == nil {
					t.Fatalf("error expected reading at index %v in %v, got %v", tc.idx, tc.bs, got)
				}

				if tc.wantErr == nil && err != nil {
					t.Fatalf("unexpected error %v", err)
				}

				if tc.wantErr != nil {
					return
				}

				if tc.want != got {
					t.Fatalf("expected %v, got %v", tc.want, got)
				}

				if tc.bs.idx != tc.idx+1 {
					t.Fatalf("expected index to increment by one to %v, got %v", tc.idx+1, tc.bs.idx)
				}
			},
		)
	}
}

func TestReadByte(t *testing.T) {
	tt := map[string]struct {
		bs      BitStream
		idx     types.Position
		want    byte
		wantErr error
	}{
		"read on empty bs": {
			bs:      NewFromBits([]Bit{}),
			idx:     0,
			want:    0,
			wantErr: errors.New(""),
		},
		"read on too short bs": {
			bs:      NewFromBits([]Bit{true, true, true, true, true, true, true}),
			idx:     0,
			want:    0,
			wantErr: errors.New(""),
		},
		"read far too long": {
			bs:      NewFromBits([]Bit{true, true, false, false, false, false, false, false, true, false}),
			idx:     2,
			want:    2,
			wantErr: nil,
		},
		"read 1 from the begining": {
			bs:      NewFromBits([]Bit{false, false, false, false, false, false, false, true, false, true, true, true, true}),
			idx:     0,
			want:    1,
			wantErr: nil,
		},
		"read 3 from pos 8": {
			bs:      NewFromBits([]Bit{false, false, false, false, false, false, false, true, false, false, false, false, false, false, true, true}),
			idx:     8,
			want:    3,
			wantErr: nil,
		},
	}

	for name, tc := range tt {
		t.Run(name,
			func(t *testing.T) {
				tc.bs.idx = tc.idx
				got, err := tc.bs.ReadByte()

				if tc.wantErr != nil && err == nil {
					t.Fatalf("error expected reading at index %v in %v, got %v", tc.idx, tc.bs, got)
				}

				if tc.wantErr == nil && err != nil {
					t.Fatalf("unexpected error %v", err)
				}

				if tc.wantErr != nil {
					return
				}

				if tc.want != got {
					t.Fatalf("expected %v, got %v", tc.want, got)
				}

				if tc.bs.idx != tc.idx+byteSize {
					t.Fatalf("expected index to increment by one to %v, got %v", tc.idx+1, tc.bs.idx)
				}
			},
		)
	}
}
