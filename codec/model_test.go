package main

import (
	"bufio"
	"testing"
)

func TestNewModel(t *testing.T) {
	// Test if a fixed model is created with frequencies read from standard input.
	// Mock stdin.
	stdin, err := setStdFile([]byte("abcdefgabcdefabcdeabcdabcaba"))
	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Reader from stdin.
	reader := bufio.NewReader(stdin)

	// Create fixed model
	fixedModel := NewModel(Fixed, reader)
	expected := map[byte]int{'a': 7, 'b': 6, 'c': 5, 'd': 4, 'e': 3, 'f': 2, 'g': 1}

	for k, v := range expected {
		i := int(k) + 1
		if fixedModel.Frequency(i) != v {
			t.Errorf("character '%d' should have a frequency of %d", i, v)
		}
	}
}

func TestInitialize(t *testing.T) {
	// Mock stdin.
	stdin, err := setStdFile([]byte(""))
	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Reader from stdin.
	reader := bufio.NewReader(stdin)
	fixedModel := NewModel(Fixed, reader)

	// Test if the fixed model returns an error when the total accumulated frequency
	// is greater than the maximum frequency allowed (256 * 64 = 16384).
	for i := 1; i < len(fixedModel.frequencies); i++ {
		fixedModel.SetFrequency(i, 64)
	}
	err = fixedModel.Initialize()
	if err == nil {
		t.Errorf("the total accumulated frequency cannot exceed the maximum allowed for the fixed model")
	}

	// Test if a new fixed model initializes non-initialized frequencies with 1.
	for i := 1; i < len(fixedModel.frequencies); i++ {
		if fixedModel.Frequency(i) == 0 {
			t.Error("frequencies must be initialized to 1")
			break
		}
	}

	// Test if a fixed model initializes proper cummulative frequencies.
}

func TestUpdateModel(t *testing.T) {
	// Mock stdin.
	stdin, err := setStdFile([]byte(""))
	if err != nil {
		t.Errorf("error mocking stdin")
	}

	// Reader from stdin.
	reader := bufio.NewReader(stdin)
	fixedModel := NewModel(Fixed, reader)

	// Test if fixed model does not get changed when frequencies are updated by the encoding process.
	f := fixedModel.frequencies
	cf := fixedModel.cummulativeFrequencies

	fixedModel.UpdateModel(25)

	if f != fixedModel.frequencies || cf != fixedModel.cummulativeFrequencies {
		t.Error("frequencies for the fixed model should not be updated")
	}
}
