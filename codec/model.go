package main

import (
	"bufio"
	"fmt"
	"io"
)

const (
	CharsCount   = 256 // Maximum number of characters
	SymbolsCount = CharsCount + 1
	EOFChar      = CharsCount + 1

	Fixed = iota
	Adaptive
)

// Object model
type Model struct {
	Class int
	// Alphabet               [SymbolsCount + 1]int
	// CharsCount   = 256,
	// SymbolsCount = CharsCount + 1,
	indexForChar           [CharsCount]int
	charForIndex           [SymbolsCount + 1]int
	frequencies            [SymbolsCount + 1]int
	cummulativeFrequencies [SymbolsCount + 1]int
}

// New creates a new model of type class.
// TODO: create a package for model so "New" can be referenced as "Model.New(Model.Fixed)" instead of "NewModel()".
func NewModel(class int, input *bufio.Reader) Model {
	model := Model{Class: class}
	if model.Class == Fixed {
		// Update frequencies according to input statistics.
		for {
			c, err := input.ReadByte()
			if err == io.EOF {
				break
			}
			// frequencies[0] = 0
			model.frequencies[int(c)+1] += 1
		}
		// Set all non-initialized frequencies to 1.
		for i := 1; i < len(model.frequencies); i++ {
			if model.Frequency(i) == 0 {
				model.SetFrequency(i, 1)
			}
		}
	}

	return model
}

// CharByIndex defines a getter for the charForIndex array.
// TODO: error management.
func (m *Model) CharByIndex(i int) int {
	return m.charForIndex[i]
}

// IndexByChar defines a getter for the indexForChar array.
// TODO: error management.
func (m *Model) IndexByChar(c int) int {
	return m.indexForChar[c]
}

// CummulativeFrequency defines a getter for the cummulativeFrequencies array.
// TODO: error management.
func (m *Model) CummulativeFrequency(i int) int {
	return m.cummulativeFrequencies[i]
}

// Frequency defines a getter for the frequencies array.
// TODO: error management.
func (m *Model) Frequency(i int) int {
	return m.frequencies[i]
}

// SetCharByIndex defines a setter for the charForIndex array.
// TODO: error management.
func (m *Model) SetCharByIndex(i, c int) {
	m.charForIndex[i] = c
}

// SetIndexByChar defines a setter for the indexForChar array.
// TODO: error management.
func (m *Model) SetIndexByChar(c, i int) {
	m.indexForChar[c] = i
}

// SetCummulativeFrequency defines a setter for the cummulativeFrequencies array.
// TODO: error management.
func (m *Model) SetCummulativeFrequency(i, f int) {
	m.cummulativeFrequencies[i] = f
}

// SetFrequency defines a setter for the frequencies array.
// TODO: error management.
func (m *Model) SetFrequency(i, f int) {
	m.frequencies[i] = f
}

// Setup tables that translate between symbol indexes and characters.
func (m *Model) setupTables() {
	for c := 0; c < CharsCount; c++ {
		i := c + 1
		m.SetIndexByChar(c, i)
		m.SetCharByIndex(i, c)
	}
}

// Initialize populates the alphabet and frequencies of the model.
func (m *Model) Initialize() error {
	// Setup tables for translation.
	// Equal procedure for adaptive and fixed models.
	m.setupTables()
	// Setup initial frequency counts.
	// TODO: make fixed and adaptive models equal in the initialization.
	// TODO (cont.): a fixed model class will receive or count frequencies after defining the default.
	if m.Class == Fixed {
		err := m.setFixedFrequencies()
		if err != nil {
			return err
		}
	} else if m.Class == Adaptive {
		m.setAdaptiveFrequencies()
	} else {
		// TODO: class of model not defined. Error.
	}
	return nil
}

// setFixedFrequencies ...
func (m *Model) setFixedFrequencies() error {
	m.SetCummulativeFrequency(SymbolsCount, 0)
	// Check if any frequency, besides index = 0 is equal to 0.
	for i := 1; i < SymbolsCount; i++ {
		if m.Frequency(i) == 0 {
			return fmt.Errorf("frequency not defined for index %d", i)
		}
	}
	// Setup cummulative frequency counts.
	fc := 0
	for i := SymbolsCount; i > 0; i-- {
		fc = m.CummulativeFrequency(i) + m.Frequency(i)
		m.SetCummulativeFrequency(i-1, fc)
	}
	if m.CummulativeFrequency(0) > MaximumFrequency {
		return fmt.Errorf("cummulative frequency at index 0 exceeds the maximum allowed")
	}
	return nil
}

// setAdaptiveFrequencies ...
func (m *Model) setAdaptiveFrequencies() {
	for i := 0; i <= SymbolsCount; i++ {
		m.SetFrequency(i, 1)
		m.SetCummulativeFrequency(i, SymbolsCount-i)
	}
	m.SetFrequency(0, 0)
}

// UpdateModel updates the statistics of symbols according to a new input char.
// Only used for adaptive models.
// TODO: UpdateModel needs adaptation for when the fixed model is chosen.
func (m *Model) UpdateModel(symbol int) {
	if m.Class == Adaptive {
		var characterIndex, characterSymbol int

		if m.CummulativeFrequency(0) == MaximumFrequency {
			m.halveCounts()
		}

		// Find symbol's new index.
		index := m.symbolIndex(symbol)

		if index < symbol {
			characterIndex = m.CharByIndex(index)
			characterSymbol = m.CharByIndex(symbol)

			m.SetCharByIndex(index, characterSymbol)
			m.SetCharByIndex(symbol, characterIndex)

			m.SetIndexByChar(characterIndex, symbol)
			m.SetIndexByChar(characterSymbol, index)
		}

		f := m.Frequency(index)
		m.SetFrequency(index, f+1)

		for index > 0 {
			index--
			cf := m.CummulativeFrequency(index)
			m.SetCummulativeFrequency(index, cf+1)
		}
	}
}

// halveCounts halves all the counts keeping them non-zero.
func (m *Model) halveCounts() {
	acc := 0
	for i := SymbolsCount; i >= 0; i-- {
		f := (m.Frequency(i) + 1) / 2
		m.SetFrequency(i, f)
		m.SetCummulativeFrequency(i, acc)
		acc += m.Frequency(i)
	}
}

// symbolIndex finds symbol's new index.
func (m *Model) symbolIndex(s int) (i int) {
	i = s
	for m.Frequency(i) == m.Frequency(i-1) {
		i--
	}
	return
}
