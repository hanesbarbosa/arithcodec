package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// setStdFile creates a temporary file to mock the standard input (stdin)
// or output (stdout) from the OS.
func setStdFile(input []byte) (*os.File, error) {
	stdfile, err := ioutil.TempFile("", "")

	if err != nil {
		return nil, fmt.Errorf("error creating temp file")
	}

	// Write bytes to stdfile.
	_, err = stdfile.Write(input)
	if err != nil {
		return nil, fmt.Errorf("error writing into temp file")
	}

	// Rewind reading point for file.
	_, err = stdfile.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("error rewinding temp file")
	}

	// Remove temporary file.
	defer os.Remove(stdfile.Name())

	return stdfile, nil
}

// TestWriteSymbol tests if one symbol is written to stdout.
func TestWriteSymbol(t *testing.T) {
	// Mock stdin with a temp file.
	stdout, err := setStdFile([]byte(""))

	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Writer defined to stdout.
	writer := bufio.NewWriter(stdout)

	// Write to stdout.
	symbol := 65
	err = writeSymbol(&symbol, writer)
	if err != nil {
		t.Error(err.Error())
	}

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

	if int(stdout_buffer[0]) != symbol {
		t.Errorf("expected stdout %b but received %b", symbol, stdout_buffer[0])
	}
}
