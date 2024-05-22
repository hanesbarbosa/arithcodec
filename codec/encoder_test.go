package main

import (
	"bufio"
	"os"
	"testing"
)

// TestSetBit tests the insertion of unique bits 0 or 1.
// Also tests if a completed buffer is written to the stdout.
func TestSetBit(t *testing.T) {
	var bits_inserted, buffer, expected int

	bits_inserted = 0
	buffer = 0
	expected = 0

	// Mock stdin with a temp file.
	stdout, err := setStdFile([]byte(""))

	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Writer defined to stdout.
	writer := bufio.NewWriter(stdout)

	// Insert byte = 11011001 bit by bit, therefore the sequence of entrance is reversed.
	// Inserting bits in the order 10011011 should result in the number 11011001.
	// Bits are inserted from left to right.
	// Since bits_inserted = 0 it should not write to stdout until bits_inserted = 8.
	setBit(1, &bits_inserted, &buffer, writer)
	expected = 0b10000000
	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}

	// 01 (2)
	setBit(0, &bits_inserted, &buffer, writer)
	expected = 0b01000000
	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}
	// 001 (3)
	setBit(0, &bits_inserted, &buffer, writer)
	expected = 0b00100000
	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}
	// 1001 (4)
	setBit(1, &bits_inserted, &buffer, writer)
	expected = 0b10010000
	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}
	// 11001 (5)
	setBit(1, &bits_inserted, &buffer, writer)
	expected = 0b11001000
	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}
	// 011001 (6)
	setBit(0, &bits_inserted, &buffer, writer)
	expected = 0b01100100
	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}
	// 1011001 (7)
	setBit(1, &bits_inserted, &buffer, writer)
	expected = 0b10110010
	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}

	// 11011001 (8)
	// Now bits_inserted should totally fill the buffer and write to standard output.
	setBit(1, &bits_inserted, &buffer, writer)
	expected = 0b11011001

	// Rewind reading point for stdout file after writing 8 bits.
	_, err = stdout.Seek(0, 0)
	if err != nil {
		t.Errorf("error rewinding stdout temp file")
	}

	// Checking value written at stdout buffer.
	stdout_buffer := make([]byte, 1)
	_, err = stdout.Read(stdout_buffer)
	if err != nil {
		t.Errorf("error reading stdout temp file")
	}

	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}

	if int(stdout_buffer[0]) != expected {
		t.Errorf("expected stdout %b but received %b", expected, stdout_buffer[0])
	}

	// Bits on the buffer should not be cleared before the next character.
	setBit(0, &bits_inserted, &buffer, writer)
	expected = 0b01101100
	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}
}

// TestSetOppositeBits tests the insertion of opposite bits.
// The setOppositeBits does not exist in the original paper.
func TestSetOppositeBits(t *testing.T) {
	var bits_inserted, buffer, opposite_bits, expected int

	bits_inserted = 0
	buffer = 0
	opposite_bits = 0
	expected = 0

	// Mock stdin with a temp file.
	stdout, err := setStdFile([]byte(""))

	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Writer defined to stdout.
	writer := bufio.NewWriter(stdout)

	// Tests if correct value is inserted.
	opposite_bits = 2
	setBit(1, &bits_inserted, &buffer, writer)
	setOppositeBits(1, &opposite_bits, &bits_inserted, &buffer, writer)
	expected = 0b00100000

	if buffer != expected {
		t.Errorf("expected value %b received %b", expected, buffer)
	}

	// We add 2 to test if the opposite_bits variable was correctly decreased in the previous function.
	opposite_bits += 2
	setOppositeBits(0, &opposite_bits, &bits_inserted, &buffer, writer)
	expected = 0b11001000

	if buffer != expected {
		t.Errorf("expected value %b received %b", expected, buffer)
	}

	// Tests if buffer full writes to the stdout.
	// opposite_bits = 2 and should be used in the next round.
	opposite_bits += 3
	setOppositeBits(1, &opposite_bits, &bits_inserted, &buffer, writer)
	expected = 0b00011001

	// Rewind reading point for stdout file after writing 8 bits.
	_, err = stdout.Seek(0, 0)
	if err != nil {
		t.Errorf("error rewinding stdout temp file")
	}

	// Checking value written at stdout buffer.
	stdout_buffer := make([]byte, 1)
	_, err = stdout.Read(stdout_buffer)
	if err != nil {
		t.Errorf("error reading stdout temp file")
	}

	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}

	if int(stdout_buffer[0]) != expected {
		t.Errorf("expected stdout %b but received %b", expected, stdout_buffer[0])
	}
}

// TestOutputBit tests if the function is moving properly the bits into the output buffer.
func TestOutputBit(t *testing.T) {
	var opposite_bits, bits_inserted, buffer, expected int
	opposite_bits = 0
	bits_inserted = 0
	buffer = 0
	expected = 0

	// Mock stdin with a temp file.
	stdout, err := setStdFile([]byte(""))

	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Writer defined to stdout.
	writer := bufio.NewWriter(stdout)

	// Tests the inclusion of a simple value having 0s zeroes and 1s ones.
	// Insert the value 10110111 (8) in the buffer.
	outputBit(1, &opposite_bits, &bits_inserted, &buffer, writer)
	outputBit(1, &opposite_bits, &bits_inserted, &buffer, writer)
	outputBit(1, &opposite_bits, &bits_inserted, &buffer, writer)
	opposite_bits += 2
	outputBit(0, &opposite_bits, &bits_inserted, &buffer, writer)
	opposite_bits += 1
	outputBit(0, &opposite_bits, &bits_inserted, &buffer, writer)
	expected = 0b10110111

	// Rewind reading point for stdout file after writing 8 bits.
	_, err = stdout.Seek(0, 0)
	if err != nil {
		t.Errorf("error rewinding stdout temp file")
	}

	// Checking value written at stdout buffer.
	stdout_buffer := make([]byte, 1)
	_, err = stdout.Read(stdout_buffer)
	if err != nil {
		t.Errorf("error reading stdout temp file")
	}

	if buffer != expected || opposite_bits != 0 {
		t.Errorf("expected value %b and 0 opposite bits but received %b and %d", expected, buffer, opposite_bits)
	}

	if int(stdout_buffer[0]) != expected {
		t.Errorf("expected stdout %b but received %b", expected, stdout_buffer[0])
	}

	// TODO: refactor.
	// Mock stdin with a temp file.
	stdout, err = setStdFile([]byte(""))

	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Writer defined to stdout.
	writer = bufio.NewWriter(stdout)

	// Insert the value 11110110 (8) in the buffer.
	opposite_bits += 2
	outputBit(0, &opposite_bits, &bits_inserted, &buffer, writer)
	opposite_bits += 4
	outputBit(0, &opposite_bits, &bits_inserted, &buffer, writer)
	expected = 0b11110110

	// Rewind reading point for stdout file after writing 8 bits.
	_, err = stdout.Seek(0, 0)
	if err != nil {
		t.Errorf("error rewinding stdout temp file")
	}

	// Checking value written at stdout buffer.
	stdout_buffer = make([]byte, 1)
	_, err = stdout.Read(stdout_buffer)
	if err != nil {
		t.Errorf("error reading stdout temp file")
	}

	if buffer != expected {
		t.Errorf("expected value %b and 0 opposite bits but received %b and %d", expected, buffer, opposite_bits)
	}

	if int(stdout_buffer[0]) != expected {
		t.Errorf("expected stdout %b but received %b", expected, stdout_buffer[0])
	}
}

// TestEncodeSymbol tests the encoding of a small text symbol by symbol.
func TestEncodeSymbol(t *testing.T) {
	var opposite_bits, bits_inserted, buffer, expected, low, high, symbol int

	opposite_bits = 0
	bits_inserted = 0
	buffer = 0
	expected = 0
	low = 0
	high = TopValue

	// Mock stdin with a temp file.
	stdout, err := setStdFile([]byte(""))

	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Writer defined to stdout.
	writer := bufio.NewWriter(stdout)

	model := NewModel(Adaptive, nil)
	model.Initialize()

	// Tests the encoding of 'eaii!' ending with a new line character.
	characters := []int{101, 97, 105, 105, 33, 10}
	expected_characters := []int{0b10110010, 0b10000001, 0b00111010, 0b11100111, 0b11111110, 0b00100000}

	for i := 0; i < len(expected_characters); i++ {
		symbol = model.IndexByChar(characters[i])
		encodeSymbol(symbol, &low, &high, &opposite_bits, &bits_inserted, &buffer, &model, writer)
		model.UpdateModel(symbol)
		if buffer != expected_characters[i] {
			t.Errorf("expected value %b but received %b", expected_characters[i], buffer)
		}
	}

	// Character received: 'EOF' (-1) translated to 257.
	// When EOF is received the character is not translated via char_to_index (in function 'main').
	// Also the statistical model is not updated.
	expected = 0b00011100
	symbol = 257
	encodeSymbol(symbol, &low, &high, &opposite_bits, &bits_inserted, &buffer, &model, writer)

	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}

	// Ending encoding procedure.
	// Done encoding also occur in this function (flush buffer to memory or stdout).
	expected = 0b01000011
	finishEncoding(&low, &opposite_bits, &bits_inserted, &buffer, writer)

	if buffer != expected {
		t.Errorf("expected value %b but received %b", expected, buffer)
	}
}

// TestEncode tests the encoding of a whole stream at once.
// By testing the Encode function we check the processing through
// the reading loop of symbols from a standard input.
func TestEncode(t *testing.T) {
	var stdin, stdout *os.File
	var err error

	// Mock stdin with a temp file having the text "this is my thing".
	stdin, err = setStdFile([]byte{116, 104, 105, 115, 32, 105, 115, 32, 109, 121, 32, 116, 104, 105, 110, 103, 10})
	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Mock stdin with a temp file.
	stdout, err = setStdFile([]byte(""))
	if err != nil {
		t.Errorf("error mocking stdout")
	}

	// stding Reader and stdout writer.
	reader := bufio.NewReader(stdin)
	writer := bufio.NewWriter(stdout)

	// Populate statistical model.
	model := NewModel(Adaptive, reader)
	model.Initialize()

	// Encode stdin and output to stdout.
	Encode(reader, writer, &model)

	// Rewind reading point for file.
	_, err = stdout.Seek(0, 0)
	if err != nil {
		t.Errorf("error rewinding temp file")
	}

	// Checks if the stdout has the encoded 'this is my thing'.
	expected_buffer := []byte{49, 80, 115, 185, 190, 99, 156, 184, 157, 215, 158, 186, 123, 187, 93, 88, 191, 43}

	// Checking value written at stdout buffer (18 bytes).
	stdout_buffer := make([]byte, 18)
	_, err = stdout.Read(stdout_buffer)
	if err != nil {
		t.Errorf("error reading stdout temp file")
	}

	for i := 0; i < len(stdout_buffer); i++ {
		if stdout_buffer[i] != byte(expected_buffer[i]) {
			t.Errorf("expected stdout %b at position %d but found %b", expected_buffer[i], i, stdout_buffer[i])
		}
	}
}
