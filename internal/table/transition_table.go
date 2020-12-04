package table

import (
	"fmt"
	"io"
	"math"

	"github.com/chavacava/next/internal/types"
)

type next struct {
	S     byte
	Count types.SymbolCountType
}

func (n next) String() string {
	return fmt.Sprintf("[%d,%d]", n.S, n.Count)
}

type NextList struct {
	List  []*next
	Grows []types.Position
}

func newNextList() NextList {
	return NextList{List: []*next{}}
}

func (nl *NextList) add(s byte) types.NextIndex {
	for i, n := range nl.List {
		if n.S == s {
			n.Count++
			return types.NextIndex(i)
		}
	}

	if len(nl.List) == 256 {
		panic(fmt.Sprintf("if we add '%v' next list will be bigger than expected (256):\n%+v\n", s, nl.List))
	}

	new := next{S: s, Count: 1}
	nl.List = append(nl.List, &new)

	return types.NextIndex(len(nl.List) - 1)
}

func (nl *NextList) dynamicBitCount(pos types.Position) byte {
	last := byte(0)
	for i, p := range nl.Grows {
		if pos < p {
			break
		}
		last = byte(i + 1)
	}

	return last + 1
}

func (nl *NextList) String() string {
	var result string
	for _, nx := range nl.List {
		result += nx.String()
	}
	return result
}

// TransitionsTable of symbols and their transitions
// use the constructor New
type TransitionsTable struct {
	Root        byte
	InputSize   types.Size
	Transitions map[byte]*NextList
}

// New table from byte stream
func New(input io.ReadSeeker) TransitionsTable {
	table := TransitionsTable{Transitions: map[byte]*NextList{}}

	var p = make([]byte, 1)
	var pSize = types.Size(len(p))
	_, err := input.Read(p)
	if err != nil {
		return table
	}
	var inputSize = pSize

	table.Root = byte(p[0])

	previous := table.Root
	var pos types.Position = 0
	for {
		_, err := input.Read(p)
		if err != nil {
			if err == io.EOF {
				break
			}

			panic(err.Error())
		}

		current := byte(p[0])
		table.addNext(previous, current, pos)
		inputSize += pSize
		previous = current
		pos++
	}

	table.InputSize = inputSize

	return table
}

func (t *TransitionsTable) addNext(from byte, to byte, pos types.Position) {
	nexts, exists := t.Transitions[from]
	if !exists {
		nl := newNextList()
		nexts = &nl
		t.Transitions[from] = nexts
	}

	nexts.add(to)

	currentNecessaryBits := byte(len(nexts.Grows)) + 1
	if minBitsCount(len(nexts.List)-1) > currentNecessaryBits {
		nexts.Grows = append(nexts.Grows, pos)
	}
}

// minBitsCount yields the minimum number of bits required to encode the given number n.
func minBitsCount(n int) byte {
	return byte(math.Log2(float64(n))) + 1
}
