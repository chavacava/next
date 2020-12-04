package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/chavacava/next/internal/compressor"
	"github.com/chavacava/next/internal/table"
)

func main() {
	doCompress := flag.Bool("c", false, "compress the input")
	doExpand := flag.Bool("e", false, "expand the input")
	input := flag.String("i", "", "input file name (defaults to stdin)")
	output := flag.String("o", "", "output file name (defaults to stdout)")
	flag.Parse()

	var err error
	reader := os.Stdin
	if *input != "" {
		reader, err = os.Open(*input)
		if err != nil {
			panic(err.Error())
		}
		defer reader.Close()
	}

	writer := os.Stdout
	if *output != "" {
		writer, err = os.Create(*output)
		if err != nil {
			panic(err.Error())
		}
	}

	if (*doCompress) && (*doExpand) {
		panic("can not do both compress and expand")
	}
	if !(*doCompress || *doExpand) {
		panic("you should ask for compressing or expanding the input")
	}

	switch {
	case *doCompress:
		t := table.New(reader)

		cx := compressor.NewCompressor(t)

		encoded := new(bytes.Buffer)
		_, err := reader.Seek(0, 0)
		if err != nil {
			panic(err.Error())
		}
		err = cx.Compress(reader, encoded)
		if err != nil {
			panic(err.Error())
		}

		_, err = writer.Write(encoded.Bytes())
		if err != nil {
			panic(err.Error())
		}
		writer.Close()

		fmt.Printf("original %d bytes\n", t.InputSize)
		fmt.Printf("encoded %d bytes\n", len(encoded.Bytes()))
		fmt.Printf("ratio %v %%\n", (1.0-float32(len(encoded.Bytes()))/float32(t.InputSize))*100)
	case *doExpand:
		dx := compressor.NewDecompressor()
		err := dx.Decompress(reader, writer)
		if err != nil {
			panic(err.Error())
		}
	}
}
