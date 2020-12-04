package encoders

import (
	"github.com/chavacava/next/internal/bitstream"
	"github.com/chavacava/next/internal/table"
)

type Constant struct {
	s byte
}

func NewConstant(s byte) Constant {
	return Constant{s}
}

func NewConstantFromNextList(nl table.NextList) Constant {
	return Constant{nl.List[0].S}
}

func (ed Constant) RecordData() bitstream.BitStream {
	return bitstream.NewFromByte(ed.s, 8)
}

func (ed Constant) Encode(to byte, bs *bitstream.BitStream) error {
	return nil
}

func (ed Constant) Decode(bs *bitstream.BitStream) (byte, error) {
	return ed.s, nil
}
