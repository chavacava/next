package compressor

import (
	"fmt"
	"io"

	"github.com/chavacava/next/internal/bitstream"
	"github.com/chavacava/next/internal/compressor/encoders"
	"github.com/chavacava/next/internal/huffman"
	"github.com/chavacava/next/internal/types"
)

// Decompressor represents a decompressor for data compressed by the Compressor
// Use the constructor to create new instances
type Decompressor struct {
	decoders map[byte]encoders.Decoder
}

// NewDecompressor yields a new decompressor
func NewDecompressor() Decompressor {
	return Decompressor{}
}

// Decompress decompresses the data it reads from the given reader and writes the result in the given writer
func (d Decompressor) Decompress(r io.Reader, w io.Writer) error {
	rootSymbol, inputSize, symbolCount, err := ReadHeader(r)
	if err != nil {
		return fmt.Errorf("error while reading the file header: %v", err)
	}

	bs := bitstream.New()
	for {
		b := make([]byte, 1)
		_, err := r.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
		bs.Append(bitstream.NewFromFullByte(b[0]))
	}

	bsp := &bs
	// setup decoders
	decoders := make(map[byte]encoders.Decoder, symbolCount)
	for i := byte(0); i < symbolCount; i++ {
		recordType, err := bs.ReadByte()
		if err != nil {
			return err
		}
		from, err := bs.ReadByte()
		if err != nil {
			return err
		}

		switch recordType {
		case 0: // Constant
			to, err := bs.ReadByte()
			if err != nil {
				return err
			}
			decoders[from] = encoders.NewConstant(to)
		case 1: // huffman tree
			tree := huffman.NewTreeFromBS(bsp)
			decoders[from] = encoders.NewHuffmanBased(tree)
		default:
			panic(fmt.Sprintf("unknown record type %v decoding %v th symbol", recordType, i+1))
		}
	}

	current := rootSymbol
	w.Write([]byte{current})
	generatedSymbolCount := types.Size(1)
	for {
		if generatedSymbolCount == inputSize {
			break
		}

		decoder, exists := decoders[current]
		if !exists {
			return fmt.Errorf("no decoder for symbol %v (when generating symbol #%v)", current, generatedSymbolCount)
		}
		var next byte
		next, err = decoder.Decode(&bs)
		if err != nil {
			return err
		}
		w.Write([]byte{next})
		current = next
		generatedSymbolCount++
	}

	return nil
}
