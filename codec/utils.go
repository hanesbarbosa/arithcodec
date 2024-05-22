package main

import (
	"bufio"
)

const (
	BitsCount        = 8
	BitsCodeValue    = 16
	TopValue         = (1 << BitsCodeValue) - 1
	FirstQuarter     = (TopValue / 4) + 1
	Half             = FirstQuarter * 2
	ThirdQuarter     = FirstQuarter * 3
	MaximumFrequency = ((1 << BitsCodeValue) / 4) - 1 //16383
)

// writeSymbol flushes bits to the standard output.
func writeSymbol(buf *int, out *bufio.Writer) error {
	err := out.WriteByte(byte(*buf))
	if err != nil {
		return err
	}
	// Flush writer's buffer to underlying stdout.
	// This can be used as a mechanism to accumulate bytes before outputing.
	// TODO: output when writer's buffer is full, being more efficient for I/O.
	err = out.Flush()
	if err != nil {
		return err
	}
	return nil
}

func subtractOffsetMiddle(l, h int) (int, int) {
	return l - FirstQuarter, h - FirstQuarter
}

func subtractOffsetTop(l, h int) (int, int) {
	return l - Half, h - Half
}

func scaleUpRange(l, h int) (int, int) {
	return (2 * l), (2*h + 1)
}

func setCodeRange(l, h int) int {
	return (h - l) + 1
}

func setHigh(cr, l int, m *Model, i int) int {
	return l + ((cr * m.CummulativeFrequency(i-1)) / m.CummulativeFrequency(0)) - 1
}

func setLow(cr, l int, m *Model, i int) int {
	return l + ((cr * m.CummulativeFrequency(i)) / m.CummulativeFrequency(0))
}
