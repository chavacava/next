package encoders

import (
	"github.com/chavacava/next/internal/bitstream"
)

type Encoder interface {
	Encode(to byte, bs *bitstream.BitStream) error
	RecordData() bitstream.BitStream
}
type Decoder interface {
	Decode(bs *bitstream.BitStream) (byte, error)
}
