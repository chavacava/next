# NEXT

Next is a prototype implementation of the _next compression algorithm_.
The compressor is a proof-of-concept of the algorithm and is not intended to be used in any way.

## The Next Compression Algorithm

The idea behind this compression algorithm is to encode file content as a [Finite State Transducer](https://en.wikipedia.org/wiki/Finite-state_transducer) plus the input sequence that, feeded to the transducer, produces the file content.

Therefore, file content is encoded as _transducer+input_. For example if we want to encode the following content:

```
Simplicity is prerequisite for reliability
```

We can represent the transducer with a transition table like the following

| State | Next |
|-------|-------| 
| o |'r'|
| y |' '|
| ' ' |'i' 'p' 'f' 'r'|
| t |'y' 'e'|
| b |'i'|
| i |'m' 'c' 't' 's' 'a' 'l'|
| p |'l' 'r'|
| u |'i'|
| f |'o'|
| S |'i'|
| s |' ' 'i'|
| c |'i'|
| r |'e' ' '|
| e |'r' 'q' ' ' 'l'|
| q |'u'|
| a |'b'|
| m |'p'|
| l |'i'|

Then, knowing that the _Start symbol_ is `S`, we can re-generate the content with the following input where each number is the index of the next symbol in the list of the current one.

```
0 0 0 0 0 1 0 2 0 0 0 3 0 1 1 0 0 0 1 0 0 3 0 2 1 2 2 0 0 1 3 0 3 0 4 0 0 5 0 2 0 
```

That is:

```
S i m p l i c i t y _ i s _ p r e r e q u i s i t e _ f o r _ r e l i a b i l i t y
 0 0 0 0 0 1 0 2 0 0 0 3 0 1 1 0 0 0 1 0 0 3 0 2 1 2 2 0 0 1 3 0 3 0 4 0 0 5 0 2 0
 ```

Of course to reach some kind of compression the input need to be encoded efficiently.
For example we can try to reduce the storage size of the input by representing it with shortest bit sequences, like this the above input could be encoded te following sequence of 76 bits (10 bytes):

```
0000000000010001000000011001100000100001100010110100011101100100000101000100
```

For:
```
S i    m p l i    c i    t y _  i    s _  p r e  r e  q u i    s i    t e  _  f o r _  r e  l i    a b i    l i    t y
 0 0000 0 0 0 0001 0 0010 0 0 00 0011 0 01 1 0 00 0 01 0 0 0011 0 0010 1 10 10 0 0 1 11 0 11 0 0100 0 0 0101 0 0010 0
```

Notice that the number of bits required to encode transitions depends on the number of transitions of each symbol. While symbol `s` has two transitions (to ' ' and `i`) thus it only requires one bit to encode them, symbol `i` has six transitions thus it requires four bits encoding.

The actual implementation uses a Huffman tree to encode the transitions of each symbol, a more efficient compression that the one described above. 

# How to...

## ...Build

Go to the `cmd` directory and run

```
go build -o next .
```

that will create `next` executable.

## ...Play with `next`

```
Usage of next:
  -c    compress the input
  -e    expand the input
  -i string
        input file name (defaults to stdin)
  -o string
        output file name (defaults to stdout)
```
