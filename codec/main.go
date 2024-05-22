// Package bramah-codec implements a library for Arithmetic Coding.
// It is based on the C implementation by Witten et al. in the paper
// "ARITHMETIC CODING FOR DATA COMPRESSION" from 1987.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func synopsis() string {
	return "Usage: bramah-codec [-e | -d] [-f | -a] [-i <input file>] [-s <statistics file>] [-o <output file>]"
}

func main() {
	// TODO:
	// 1 - Read input from file. Write output to file.
	// 2 - Create statistics from file for fixed model (experimental).
	// 3 - Study a way to encode fixed statistics (experimental).

	// Define CLI flags.

	// -e encode, -d decode
	e := flag.Bool("e", false, "Encode standard input")
	d := flag.Bool("d", false, "Decode standard input")

	// -f fixed, -a adaptive
	f := flag.Bool("f", false, "Fixed model")
	a := flag.Bool("a", false, "Adaptive model")

	// -i input file, -o output file
	inFileName := flag.String("i", "", "Input file")
	statsFileName := flag.String("s", "", "Input file for frequencies' samples")
	outFileName := flag.String("o", "", "Output file")

	flag.Parse()

	// Input source.
	var input, stats *bufio.Reader
	var inFile, statsFile *os.File
	var err error

	if *inFileName == "" {
		// User did not choose an input file.
		inFile = os.Stdin
	} else {
		// User did choose an input file.
		// TODO: check what is the best FileMode.
		inFile, err = os.OpenFile(*inFileName, os.O_RDONLY, 0755)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
	input = bufio.NewReader(inFile)

	// Output destination.
	var output *bufio.Writer
	var outFile *os.File

	if *outFileName == "" {
		// User did not choose an output file.
		outFile = os.Stdout
	} else {
		// User did choose an output file.
		// TODO: check what is the best FileMode.
		outFile, err = os.OpenFile(*outFileName, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
	output = bufio.NewWriter(outFile)

	// Choose model.
	var model Model

	switch {
	case *f:
		// Check if statistics are available for decoding
		// in the fixed model.
		if *d {
			statsFile, err = os.OpenFile(*statsFileName, os.O_RDONLY, 0755)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			stats = bufio.NewReader(statsFile)
			// Model is created based on the frequencies of the original file.
			// By doing this it is possible to decode with the fixed model that originated the code.
			model = NewModel(Fixed, stats)
			// Close stats file.
			if err := statsFile.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		} else {
			model = NewModel(Fixed, input)
		}
	case *a:
		model = NewModel(Adaptive, input)
	default:
		// TODO: printing function and complete description.
		fmt.Fprintln(os.Stderr, synopsis())
		os.Exit(1)
	}

	err = model.Initialize()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Choose algorithm.
	switch {
	case *e:
		// Rewind input file before encoding.
		_, err = inFile.Seek(0, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		// Encode.
		Encode(input, output, &model)
	case *d:
		// Decode.
		Decode(input, output, &model)
	default:
		// TODO: printing function and complete description.
		fmt.Fprintln(os.Stderr, synopsis())
		os.Exit(1)
	}

	// Close input file.
	if err := inFile.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// Close output file.
	if err := outFile.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
