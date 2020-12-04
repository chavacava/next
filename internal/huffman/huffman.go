package huffman

import (
	"fmt"
	"sort"

	"github.com/chavacava/next/internal/bitstream"
)

type node interface {
	value() uint
	String() string
}

// Tree represents a Huffman tree
// To be instantiated using a constructor
type Tree struct {
	root node
}

// NewTree yields a tree corresponding to the given list of symbol frequencies
func NewTree(fs []SymbolFreq) Tree {
	// Sort frequencies
	sort.Sort(byFreq(fs))

	wrkList := []node{}
	for _, f := range fs {
		wrkList = append(wrkList, f)
	}

	for {
		if len(wrkList) < 2 {
			break
		}

		newNode := makeNewNode(wrkList[0], wrkList[1])

		wrkList = insertItem(wrkList[2:], newNode)
	}

	return Tree{wrkList[0]}
}

// NewTreeFromBS yields a tree from a bitstream encoding of a tree
func NewTreeFromBS(bs *bitstream.BitStream) Tree {
	root := newTreeFromBS(bs)
	return Tree{root: root}
}

func newTreeFromBS(bs *bitstream.BitStream) node {
	b, err := bs.Read()
	if err != nil {
		panic(err)
	}
	switch b {
	case intNodeMarker:
		newNode := intNode{}
		left := newTreeFromBS(bs)
		right := newTreeFromBS(bs)
		newNode.left = left
		newNode.right = right
		return newNode
	default: //case leafNodeMarker:
		newNode := SymbolFreq{}
		newNode.Symbol, err = bs.ReadByte()
		if err != nil {
			panic(err)
		}
		return newNode
	}
}

func (t Tree) String() string {
	return t.root.String()
}

// DictionaryType represents a look up table from bytes to its corresponding Huffman codes
type DictionaryType map[byte]bitstream.BitStream

// Dictionary returns the dictionary defined by this Huffman tree
func (t Tree) Dictionary() DictionaryType {
	result := DictionaryType{}
	t.buildDictionary(result, bitstream.BitStream{}, t.root)
	return result
}

// Interpret yilds a byte by interpreting the given bitstream on this Huffman tree
func (t Tree) Interpret(bs *bitstream.BitStream) byte {
	return t.walk(t.root, bs)
}

// AsBitstreams encodes this Huffman tree in a bitstream
func (t Tree) AsBitstream() bitstream.BitStream {
	result := bitstream.BitStream{}
	t.asBitstream(&result, t.root)
	return result
}

const intNodeMarker = false
const leafNodeMarker = true

func (t Tree) asBitstream(bs *bitstream.BitStream, n node) {
	switch nt := n.(type) {
	case intNode:
		bs.Append(bitstream.NewFromBits([]bitstream.Bit{intNodeMarker}))
		t.asBitstream(bs, nt.left)
		t.asBitstream(bs, nt.right)
	case SymbolFreq:
		bs.Append(bitstream.NewFromBits([]bitstream.Bit{leafNodeMarker}))
		newBS := bitstream.NewFromByte(nt.Symbol, 8)
		bs.Append(newBS)
	default:
		panic(fmt.Sprintf("unknown Huffman tree node type %t", nt))
	}
}

const rightFlag = true
const leftFlag = false

func (t Tree) walk(n node, bs *bitstream.BitStream) byte {
	switch nt := n.(type) {
	case intNode:
		var child node
		bit, err := bs.Read()
		if err != nil {
			panic(err)
		}
		switch bit {
		case rightFlag:
			child = nt.right
		case leftFlag:
			child = nt.left
		}
		return t.walk(child, bs)
	case SymbolFreq:
		return nt.Symbol
	default:
		panic(fmt.Sprintf("unknown Huffman tree node type %t", nt))
	}
}

func (t Tree) buildDictionary(d DictionaryType, bs bitstream.BitStream, n node) {
	switch nt := n.(type) {
	case intNode:
		if nt.left != nil {
			newBS := bitstream.New()
			newBS.Append(bs)
			newBS.Append(bitstream.NewFromBits([]bitstream.Bit{leftFlag}))
			t.buildDictionary(d, newBS, nt.left)
		}
		if nt.right != nil {
			newBS := bitstream.New()
			newBS.Append(bs)
			newBS.Append(bitstream.NewFromBits([]bitstream.Bit{rightFlag}))
			t.buildDictionary(d, newBS, nt.right)
		}
	case SymbolFreq:
		code := bitstream.New()
		code.Append(bs)
		d[nt.Symbol] = code
	default:
		panic(fmt.Sprintf("unknown Huffman tree node type %t", nt))
	}
}

func makeNewNode(l, r node) node {
	result := intNode{left: l, right: r, val: l.value() + r.value()}

	if in, ok := l.(intNode); ok {
		in.father = &result
	}

	if in, ok := r.(intNode); ok {
		in.father = &result
	}

	return result
}

func insertItem(wl []node, item node) []node {
	for i, it := range wl {
		if it.value() < item.value() {
			continue
		}

		return append(wl[:i], append([]node{item}, wl[i:]...)...)
	}

	return append(wl, item)
}

type intNode struct {
	father *intNode // TODO remove father
	left   node
	right  node
	val    uint
}

func (n intNode) String() string {
	return fmt.Sprintf("{ l: %v r: %v }", n.left.String(), n.right.String())
}

func (n intNode) value() uint {
	return n.val
}

// SymbolFreq represents a symbol and its frequency
type SymbolFreq struct {
	Symbol byte
	Count  uint
}

func (f SymbolFreq) String() string {
	return fmt.Sprintf("[ S: %v ]", f.Symbol)
}

func (f SymbolFreq) value() uint {
	return f.Count
}

type byFreq []SymbolFreq

func (bf byFreq) Len() int           { return len(bf) }
func (bf byFreq) Less(i, j int) bool { return bf[i].value() < bf[j].value() }
func (bf byFreq) Swap(i, j int)      { bf[i], bf[j] = bf[j], bf[i] }
