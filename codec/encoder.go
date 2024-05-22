package main

import (
	"bufio"
	"io"
)

// Encode codes an input stream (i.e., standard input) as
// a code on the output stream (i.e., standard output).
// TODO: error management.
func Encode(input *bufio.Reader, output *bufio.Writer, model *Model) {
	var symbol, buffer, low, high, oppositeBits, bitsInserted int

	buffer = 0
	low = 0
	high = TopValue
	oppositeBits = 0
	bitsInserted = 0

	// Loop through characters coming from stdin.
	for {
		ch, err := input.ReadByte()
		if err == io.EOF {
			break
		}
		symbol = model.IndexByChar(int(ch))
		encodeSymbol(symbol, &low, &high, &oppositeBits, &bitsInserted, &buffer, model, output)
		model.UpdateModel(symbol)
	}
	encodeSymbol(EOFChar, &low, &high, &oppositeBits, &bitsInserted, &buffer, model, output)
	finishEncoding(&low, &oppositeBits, &bitsInserted, &buffer, output)
}

// encodeSymbol encodes unique symbols.
func encodeSymbol(symbol int, low *int, high *int, oppositeBits *int, bitsInserted *int, buffer *int, model *Model, output *bufio.Writer) {
	codeRange := setCodeRange(*low, *high)

	// Narrow the code region to that allotted to this symbol.
	*high = setHigh(codeRange, *low, model, symbol)
	*low = setLow(codeRange, *low, model, symbol)

	for {
		if *high < Half {
			outputBit(0, oppositeBits, bitsInserted, buffer, output)
		} else if *low >= Half {
			outputBit(1, oppositeBits, bitsInserted, buffer, output)
			*low, *high = subtractOffsetTop(*low, *high)
		} else if *low >= FirstQuarter && *high < ThirdQuarter {
			*oppositeBits += 1
			*low, *high = subtractOffsetMiddle(*low, *high)
		} else {
			break
		}
		*low, *high = scaleUpRange(*low, *high)
	}
}

// finishEncoding writes the leftover bits to the standard output.
func finishEncoding(low *int, oppositeBits *int, bitsInserted *int, buffer *int, output *bufio.Writer) {
	*oppositeBits += 1
	if *low < FirstQuarter {
		outputBit(0, oppositeBits, bitsInserted, buffer, output)
	} else {
		outputBit(1, oppositeBits, bitsInserted, buffer, output)
	}
	*buffer >>= *bitsInserted

	// TODO: Manage errors properly returning to the calling function.
	// Output to stdout.
	writeSymbol(buffer, output)
}

// outputBit prepares the next batch (buffer) of bits to be written into memory.
// This function has the same purpose of the bit_plus_follow function from the original paper.
// Pointers are used to avoid wrong values by copy, and therefore an unique value is consistent throughout the application.
func outputBit(bit int, oppositeBits *int, bitsInserted *int, buffer *int, output *bufio.Writer) {
	setBit(bit, bitsInserted, buffer, output)
	if *oppositeBits > 0 {
		setOppositeBits(bit, oppositeBits, bitsInserted, buffer, output)
	}
}

// setOppositeBits concatenates opposite bits on the buffer.
func setOppositeBits(bit int, oppositeBits *int, bitsInserted *int, buffer *int, output *bufio.Writer) {
	// We can XOR the bit received to get the opposite value.
	// 0 ^ 1 = 1, 1 ^ 1 = 0.
	inverseBit := bit ^ 1
	for *oppositeBits > 0 {
		setBit(inverseBit, bitsInserted, buffer, output)
		*oppositeBits -= 1
	}
}

// setBit writes bits to the buffer and 8-bits symbol
// to the standard output when completed.
func setBit(bit int, bitsInserted *int, buffer *int, output *bufio.Writer) {
	// Move bits to the right in order to insert a new one
	// in the left-most allowed position.
	*buffer >>= 1
	// Shift bit (0 or 1) from the right-most position to
	// the left-most allowed position.
	bit <<= BitsCount - 1
	// Execute an OR operation to insert the new bit into
	// the left-most allowed position, while preserving the
	// existent bits on the right.
	*buffer |= bit
	// Count how many bits were inserted.
	*bitsInserted += 1

	if *bitsInserted >= BitsCount {
		*bitsInserted = 0
		writeSymbol(buffer, output)
	}
}
