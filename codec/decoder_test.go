package main

import (
	"bufio"
	"os"
	"testing"
)

func TestLeastSignificantBit(t *testing.T) {
	// Mock stdin.
	stdin, err := setStdFile([]byte("e"))

	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Reader from stdin.
	reader := bufio.NewReader(stdin)

	// Read one byte from mocked input.
	// Expecting bits from character 'e' 101 = 01100101.
	buffer := 0
	bits_read := BitsCount
	expected_bits := []int{0, 1, 1, 0, 0, 1, 0, 1}

	var bit int
	for i := len(expected_bits) - 1; i >= 0; i-- {
		bit = leastSignificantBit(&buffer, &bits_read, reader)
		if bit != expected_bits[i] {
			t.Errorf("\nexpecting bit %d but got %d", expected_bits[i], bit)
		}
	}

	// Expecting EOF as -1.
	bit = leastSignificantBit(&buffer, &bits_read, reader)
	if bit != -1 {
		t.Errorf("\nexpecting -1 but got %d", bit)
	}
}

// TestInitialCodeValue tests if the code value is properly calculated for the first 16 bits.
func TestInitialCodeValue(t *testing.T) {
	// Mock stdin.
	stdin, err := setStdFile([]byte{0b11011001, 0b01000000})

	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Reader from stdin.
	reader := bufio.NewReader(stdin)

	// Read one byte from mocked input.
	// Expecting bits from character 'e' 101 = 01100101.
	buffer := 0
	bits_read := BitsCount

	// Calculate initial code value.
	// Stdin = [11011001, 01000000].
	// Code value = 39682.
	value := initialCodeValue(&buffer, &bits_read, reader)
	expected_value := 39682

	if value != expected_value {
		t.Errorf("expected code value %d but got %d", expected_value, value)
	}
}

// TestDecodeSymbol tests if symbols are properly decoded for the 'eaii!\n' input.
func TestDecodeSymbol(t *testing.T) {
	// Mock stdin = [217, 64, 157, 251, 7, 68, 14, 67].
	stdin, err := setStdFile([]byte{217, 64, 157, 251, 7, 68, 14, 67})

	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Reader from stdin.
	reader := bufio.NewReader(stdin)

	// Initiate model.
	model := NewModel(Adaptive, reader)
	model.Initialize()

	// Initiate variables.
	var buffer, bits_read, value, low, high, symbol, character int

	buffer = 0
	bits_read = BitsCount
	value = initialCodeValue(&buffer, &bits_read, reader)
	low = 0
	high = TopValue

	// Decode symbols.
	// Expected stdout = [101, 97, 105, 105, 33, 10].
	expected_characters := []int{101, 97, 105, 105, 33, 10}

	for i := 0; i < len(expected_characters); i++ {
		symbol = decodeSymbol(&value, &low, &high, &bits_read, &buffer, &model, reader)
		character = model.CharByIndex(symbol)
		if character != expected_characters[i] {
			t.Errorf("expected symbol %d but got %d", expected_characters[i], character)
		}
		model.UpdateModel(symbol)
	}
}

// TestDecode tests the decoding of a whole stream at once.
// We test the Decode function by using a phrase that is known to create an error
// when a boolean variable to stop the decoding is not used.
// The error makes the getBit function from decoding bring excessive symbols
// (more than provided) from stdin.
func TestDecode(t *testing.T) {
	var stdin, stdout *os.File
	var err error

	// Mock stdin with a temp file having the encoded version of the text "this is my thing".
	stdin, err = setStdFile([]byte{49, 80, 115, 185, 190, 99, 156, 184, 157, 215, 158, 186, 123, 187, 93, 88, 191, 43})
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
	Decode(reader, writer, &model)

	// Rewind reading point for file.
	_, err = stdout.Seek(0, 0)
	if err != nil {
		t.Errorf("error rewinding temp file")
	}

	// Checks if the stdout has the encoded 'this is my thing'.
	expected_buffer := []byte{116, 104, 105, 115, 32, 105, 115, 32, 109, 121, 32, 116, 104, 105, 110, 103, 10}

	// Checking value written at stdout buffer (17 bytes).
	// We create 1 more byte to potentially receive a wrong symbol.
	stdout_buffer := make([]byte, 18)
	_, err = stdout.Read(stdout_buffer)
	if err != nil {
		t.Errorf("error reading stdout temp file")
	}

	i := 0
	for i = 0; i < len(expected_buffer); i++ {
		if stdout_buffer[i] != byte(expected_buffer[i]) {
			t.Errorf("expected stdout %b at position %d but found %b", expected_buffer[i], i, stdout_buffer[i])
		}
	}

	// Checking if Decode passed the end of the given input stream.
	if stdout_buffer[i] != 0 {
		t.Errorf("more decoded symbols than necessary")
	}
}
