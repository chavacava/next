package bitstream

import (
	"fmt"
	"math"

	"github.com/chavacava/next/internal/types"
)

// Bit unit of information type
type Bit bool

// BitStream represents a stream of bits
// Use one of the available constructors to instantiate it
type BitStream struct {
	bits []Bit
	idx  types.Position
}

const byteSize = 8

// New yields a new, empty, BitStream
func New() BitStream {
	return BitStream{[]Bit{}, 0}
}

// NewFromFullByte yields a bitstream from the given byte, thus the bitstream length will 8
func NewFromFullByte(b byte) BitStream {
	return NewFromByte(b, byteSize)
}

// NewFromByte creates a bitStream of the given size < 8 from the given byte b
func NewFromByte(b byte, size byte) BitStream {
	result := New()
	bits := fmt.Sprintf("%08b", b)

	if math.Log2(float64(b)) > float64(size) {
		panic(fmt.Sprintf("canont represent %v in %d bits", b, size))
	}

	fitted := bits[len(bits)-int(size):]
	for _, b := range fitted {
		result.bits = append(result.bits, b == '1')
	}

	return result
}

// NewFromBits yields a bitstream containing the given bits
func NewFromBits(bits []Bit) BitStream {
	bs := New()
	newBits := make([]Bit, len(bits))
	copy(newBits, bits)
	bs.bits = newBits

	return bs
}

// newFromBytes yields a bitstream containing the bits of the given bytes
func newFromBytes(bytes []byte) BitStream {
	result := New()
	for _, b := range bytes {
		newBS := NewFromFullByte(b)
		result.Append(newBS)
	}

	return result
}

// Append appends the given bitstream to this one
func (bs *BitStream) Append(other BitStream) {
	bs.bits = append(bs.bits, other.bits...)
}

// Bytes yields the slice of bytes representation of this bitstream
// the last byte might need to be padded with 0s
func (bs BitStream) Bytes() []byte {
	result := make([]byte, 0)
	i := 0
	for ; i+byteSize <= len(bs.bits); i += byteSize {
		slice := bs.bits[i : i+byteSize]
		result = append(result, NewFromBits(slice).Byte())
	}

	rest := bs.bits[i:]
	if len(rest) != 0 {
		// pad
		for len(rest) < byteSize { // pad the byte
			rest = append(rest, Bit(false))
		}
		result = append(result, NewFromBits(rest).Byte())
	}

	return result
}

// Read yields the next bit of this bitstream
// Error will arise if there is no more bits to read
func (bs *BitStream) Read() (bool, error) {
	if bs.idx >= types.Position(len(bs.bits)) {
		return false, fmt.Errorf("reading position %v out of range %v", bs.idx, len(bs.bits)-1)
	}

	b := bs.bits[bs.idx]
	bs.idx++

	return bool(b), nil
}

// ReadByte yields the byte representation of the next 8 bits of this stream.
// Error will arise if there is not at least 8 bits to read
func (bs *BitStream) ReadByte() (byte, error) {
	if bs.idx+byteSize-1 >= types.Position(len(bs.bits)) {
		return 0, fmt.Errorf("reading positions %v to %v out of range %v", bs.idx, bs.idx+byteSize-1, len(bs.bits)-1)
	}

	next8bits := bs.bits[bs.idx : bs.idx+byteSize]
	result := NewFromBits(next8bits).Byte()

	bs.idx += byteSize

	return result, nil
}

// Byte yields the byte representation of this stream
// Error will arise if the length of the stream is bigger than 8
func (bs BitStream) Byte() byte {
	bitCount := len(bs.bits)
	if bitCount > 8 {
		panic(fmt.Sprintf("cannot convert into byte, bitstream too long %d", bitCount))
	}

	applicablePows := []byte{128, 64, 32, 16, 8, 4, 2, 1}
	offset := byteSize - bitCount
	result := byte(0)
	for i := bitCount - 1; i >= 0; i-- {
		b := bs.bits[i]
		if b {
			result += applicablePows[offset+i]
		}
	}

	return result
}

// IsEqual returns true if this stream contains the same bits and
// in the same order than the given other stream
func (bs BitStream) IsEqual(other BitStream) bool {
	if len(bs.bits) != len(other.bits) {
		return false
	}

	for i, b := range bs.bits {
		if b != other.bits[i] {
			return false
		}
	}

	return true
}
