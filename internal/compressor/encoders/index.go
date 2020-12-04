package encoders

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/chavacava/next/internal/bitstream"
	"github.com/chavacava/next/internal/table"
	"github.com/chavacava/next/internal/types"
)

type IndexBased struct {
	next    []byte
	idxSize byte
}

func NewIndexBased(nl table.NextList) IndexBased {
	result := IndexBased{
		next:    make([]byte, len(nl.List)),
		idxSize: minBitsCount(len(nl.List)),
	}

	if result.idxSize > 8 {
		panic(fmt.Sprintf("%v elements requires %v bits\n", len(nl.List), result.idxSize))
	}

	for i, n := range nl.List {
		result.next[i] = n.S
	}

	return result
}

func (ed IndexBased) Encode(to byte, bs *bitstream.BitStream, pos types.Position) error {
	idx, idxSize, err := ed.indexOf(to)
	if err != nil {
		return err
	}
	bs.Append(bitstream.NewFromByte(byte(idx), idxSize))

	return nil
}

func (ed IndexBased) Decode(w io.Writer, bs bitstream.BitStream, pos types.Position) error {
	return errors.New("not yet implemented")
}

func (ed IndexBased) indexOf(to byte) (idx types.NextIndex, idxSize byte, err error) {
	for i, n := range ed.next {
		if n == to {
			if i > 255 {
				panic("index bigger than expected")
			}

			return types.NextIndex(i), ed.idxSize, nil
		}
	}

	return 0, 0, fmt.Errorf("%v not found in the next list", to)
}

// minBitsCount yields the minimum number of bits required to encode the given number n.
func minBitsCount(n int) byte {
	return byte(math.Log2(float64(n-1))) + 1
}
