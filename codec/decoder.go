package main

import (
	"bufio"
	"io"
)

var doneDecoding bool

// Decode decodes an encoded input stream (i.e., standard input) as
// its original source on the output stream (i.e., standard output).
func Decode(input *bufio.Reader, output *bufio.Writer, model *Model) {
	var symbol, buffer, low, high, codeValue int

	buffer = 0
	low = 0
	high = TopValue

	doneDecoding = false
	bitsRead := BitsCount
	codeValue = initialCodeValue(&buffer, &bitsRead, input)

	for {
		symbol = decodeSymbol(&codeValue, &low, &high, &bitsRead, &buffer, model, input)
		if symbol == EOFChar || doneDecoding {
			break
		}
		character := model.CharByIndex(symbol)
		writeSymbol(&character, output)
		model.UpdateModel(symbol)
	}
}

// decodeSymbol decodes unique symbols.
func decodeSymbol(codeValue, low, high *int, bitsRead *int, buffer *int, model *Model, input *bufio.Reader) int {
	var symbol int
	codeRange := setCodeRange(*low, *high)
	accumulated := ((*codeValue-*low+1)*model.CummulativeFrequency(0) - 1) / codeRange

	for symbol = 1; model.CummulativeFrequency(symbol) > accumulated; symbol++ {
	}

	*high = setHigh(codeRange, *low, model, symbol)
	*low = setLow(codeRange, *low, model, symbol)

	for {
		if *high < Half {
		} else if *low >= Half {
			*codeValue -= Half
			*low, *high = subtractOffsetTop(*low, *high)
		} else if *low >= FirstQuarter && *high < ThirdQuarter {
			*codeValue -= FirstQuarter
			*low, *high = subtractOffsetMiddle(*low, *high)
		} else {
			break
		}
		*low, *high = scaleUpRange(*low, *high)
		*codeValue = (2 * *codeValue) + leastSignificantBit(buffer, bitsRead, input)
	}

	return symbol
}

// initialCodeValue sets the initial code value used as a base to decode.
func initialCodeValue(buffer *int, bitsRead *int, input *bufio.Reader) int {
	codeValue := 0
	for i := 1; i <= BitsCodeValue; i++ {
		codeValue = (2 * codeValue) + leastSignificantBit(buffer, bitsRead, input)
	}
	return codeValue
}

// getBit returns the LSB (Least Significant Bit) or right-most bit.
func leastSignificantBit(buffer *int, bitsRead *int, input *bufio.Reader) int {
	var character byte
	var bit int
	var err error

	if *bitsRead == BitsCount {
		character, err = input.ReadByte()
		if err == io.EOF {
			doneDecoding = true
			return -1
		}
		*buffer = int(character)
		*bitsRead = 0
	}

	bit = *buffer & 1
	*buffer >>= 1
	*bitsRead += 1

	return bit
}
