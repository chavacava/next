package encoders

import (
	"errors"
	"io"

	"github.com/chavacava/next/internal/bitstream"
	"github.com/chavacava/next/internal/table"
	"github.com/chavacava/next/internal/types"
)

type GrowingIndex struct {
	IndexBased
	grows []types.Position
}

func NewGrowingIndex(nl table.NextList) GrowingIndex {
	result := GrowingIndex{
		NewIndexBased(nl),
		nl.Grows,
	}

	return result
}

func (ed GrowingIndex) Encode(to byte, bs *bitstream.BitStream, pos types.Position) error {
	idx, idxSize, err := ed.indexOf(to)
	if err != nil {
		return err
	}

	dbc := ed.dynamicBitCount(pos)
	idxSize = min(idxSize, dbc)

	bs.Append(bitstream.NewFromByte(byte(idx), idxSize))

	return nil
}

func (ed GrowingIndex) Decode(w io.Writer, bs bitstream.BitStream, pos types.Position) error {
	return errors.New("not yet implemented")
}

func (ed GrowingIndex) dynamicBitCount(pos types.Position) byte {
	last := byte(0)
	for i, p := range ed.grows {
		if pos < p {
			break
		}
		last = byte(i + 1)
	}

	return last + 1
}

func min(a, b byte) byte {
	if a < b {
		return a
	}

	return b
}
