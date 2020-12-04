package encoders

import (
	"fmt"

	"github.com/chavacava/next/internal/bitstream"
	"github.com/chavacava/next/internal/huffman"
	"github.com/chavacava/next/internal/table"
)

type HuffmanBased struct {
	dictionary map[byte]bitstream.BitStream
	tree       huffman.Tree
}

func NewHuffmanBased(tree huffman.Tree) HuffmanBased {
	return HuffmanBased{
		dictionary: tree.Dictionary(),
		tree:       tree,
	}
}

func NewHuffmanBasedFromNextList(nl table.NextList) HuffmanBased {
	frequencies := make([]huffman.SymbolFreq, len(nl.List))
	for i, n := range nl.List {
		frequencies[i] = huffman.SymbolFreq{Symbol: n.S, Count: uint(n.Count)}
	}

	tree := huffman.NewTree(frequencies)

	result := HuffmanBased{
		dictionary: tree.Dictionary(),
		tree:       tree,
	}

	return result
}

func (ed HuffmanBased) RecordData() bitstream.BitStream {
	return ed.tree.AsBitstream()
}

func (ed HuffmanBased) Encode(to byte, bs *bitstream.BitStream) error {
	code, exists := ed.dictionary[to]
	if !exists {
		return fmt.Errorf("unknown symbol %v in dictionary", to)
	}

	bs.Append(code)

	return nil
}

func (ed HuffmanBased) Decode(bs *bitstream.BitStream) (byte, error) {
	s := ed.tree.Interpret(bs)
	return s, nil
}
