package compressor

import (
	"errors"
	"fmt"
	"io"

	"github.com/chavacava/next/internal/bitstream"
	"github.com/chavacava/next/internal/compressor/encoders"
	"github.com/chavacava/next/internal/table"
	"github.com/chavacava/next/internal/types"
)

// Compressor represents a data compressor
// use the constructor to create new instances
type Compressor struct {
	tt  table.TransitionsTable
	eds map[byte]encoders.Encoder
}

// NewCompressor yields a new compressor from the basis of the given
// transition table
func NewCompressor(tt table.TransitionsTable) Compressor {
	result := Compressor{
		tt:  tt,
		eds: make(map[byte]encoders.Encoder, len(tt.Transitions)),
	}
	// setup endcoders
	for s, nl := range tt.Transitions {
		result.eds[s] = encoderFactory(*nl)
	}

	return result
}

func encoderFactory(nl table.NextList) encoders.Encoder {
	s := len(nl.List)
	switch {
	case s == 1:
		return encoders.NewConstantFromNextList(nl)
	default:
		return encoders.NewHuffmanBasedFromNextList(nl)
	}
}

// Compress compresses the content from input and writes the result in the given writer
func (c Compressor) Compress(input io.Reader, w io.Writer) error {
	WriteHeader(w, c.tt.Root, c.tt.InputSize, byte(len(c.tt.Transitions)))
	var p = make([]byte, 1)
	_, err := input.Read(p)
	if err != nil {
		return errors.New("Unable to read the input, maybe is it empty?")
	}

	current := p[0]
	var compressedContent = bitstream.BitStream{}
	var pos types.Position = 0
	for {
		_, err := input.Read(p)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		next := p[0]

		encoder := c.eds[current]
		encoder.Encode(next, &compressedContent)

		current = next
		pos++
	}

	binaryRecords := bitstream.BitStream{}
	for from, e := range c.eds {
		recordHeader := bitstream.BitStream{}
		switch e.(type) {
		case encoders.Constant:
			recordHeader.Append(bitstream.NewFromFullByte(0)) // add constant for 0 constant record type
		case encoders.HuffmanBased:
			recordHeader.Append(bitstream.NewFromFullByte(1)) // add constant for 1 huffman tree record type
		default:
			panic(fmt.Sprintf("unknown encoder type %t", e))
		}

		recordHeader.Append(bitstream.NewFromFullByte(from))

		binaryRecords.Append(recordHeader)
		binaryRecords.Append(e.RecordData())
	}

	binaryRecords.Append(compressedContent)
	w.Write(binaryRecords.Bytes())

	return nil
}
