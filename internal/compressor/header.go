package compressor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/chavacava/next/internal/types"
)

/*

Position	Size 	What 		 		Example/Comment
0        	9    	magic      			\211 N E X T \r \n \032 \n
9			1		version number		major version
10			2		data offset			120 (relative to the start of the file)
12			1		root byte			65
13			8		original length 	25487852
21			1		trans recods count	number of transition records in this file
22			1		chksum				addition (overflowed) of previous bytes
xx			x		trans records


#  Transitions Record
Position	Size 	What 		 	Example/Comment
0			1		record type		1
1			1		from			65
x 			~ 		record data

# Record data by type

## Constant (type #0)
Position	Size 	What 		 	Example/Comment
0			1		to

## Huffman Tree (type #1)
Position	Size 	What 		 	Example/Comment
0			~		bs of tree

*/

// magic \211 N E X T \r \n \032 \n
var magic = []byte{137, 78, 69, 88, 84, 13, 10, 26, 10}

const versionNumber = uint8(0)

type length uint64
type offset uint16
type recordType uint8
type frequency uint64

// Record types[]byte

// Next List
const recordTypeNextList = recordType(0)

func WriteHeader(w io.Writer, rootSymbol byte, inputSize types.Size, symbolCount byte) {
	const headerSize = 22
	const offset = offset(headerSize + 1)

	var header = []interface{}{
		magic,
		versionNumber,
		offset,
		rootSymbol,
		inputSize,
		symbolCount,
	}

	buf := new(bytes.Buffer)
	for _, v := range header {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			panic(fmt.Sprintf("failed to write header: %v", err))
		}
	}
	if buf.Len() != headerSize {
		panic(fmt.Sprintf("bad header size %v, expected %v\nheader:%+v", buf.Len(), headerSize, buf.Bytes()))
	}

	checksum := checksum(buf.Bytes())
	err := binary.Write(buf, binary.LittleEndian, checksum)
	if err != nil {
		panic(fmt.Sprintf("failed to write header's checksum: %v", err))
	}

	w.Write(buf.Bytes())
}

func ReadHeader(r io.Reader) (rootSymbol uint8, inputSize types.Size, symbolCount byte, err error) {
	mgc := make([]byte, len(magic))
	l, err := r.Read(mgc)
	if err != nil {
		return 0, 0, 0, err
	}
	if l != len(magic) {
		return 0, 0, 0, fmt.Errorf("expected to read %d bytes of magic file header, got %d", len(magic), l)
	}
	if !reflect.DeepEqual(magic, mgc) {
		return 0, 0, 0, fmt.Errorf("expected to magic file header to be\n\t%v\ngot\n\t%v", magic, mgc)
	}

	//version number
	vn := make([]byte, 1)
	_, err = r.Read(vn)
	if err != nil {
		return 0, 0, 0, err
	}

	//records offset
	ro := make([]byte, 2)
	_, err = r.Read(ro)
	if err != nil {
		return 0, 0, 0, err
	}

	//root symbol
	rs := make([]byte, 1)
	_, err = r.Read(rs)
	if err != nil {
		return 0, 0, 0, err
	}
	rootSymbol = rs[0]

	//input length
	il := make([]byte, 8)
	_, err = r.Read(il)
	if err != nil {
		return 0, 0, 0, err
	}

	inputSize = types.Size(binary.LittleEndian.Uint64(il))

	//symbol count
	sc := make([]byte, 1)
	_, err = r.Read(sc)
	if err != nil {
		return 0, 0, 0, err
	}

	symbolCount = sc[0]

	//root symbol
	cs := make([]byte, 1)
	_, err = r.Read(cs)
	if err != nil {
		return 0, 0, 0, err
	}

	// TODO check sum

	return rootSymbol, inputSize, symbolCount, nil
}

func checksum(bs []byte) byte {
	var result byte
	for _, b := range bs {
		result += b
	}

	return result
}
