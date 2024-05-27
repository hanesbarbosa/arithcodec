# arithcodec

Arithmetic Coding compression.

## Installation

Install the library using `go get`:

```bash
go get github.com/Symetrix-Corp/bramah-demo-arithmetic-coding/bramah-codec
```

## Usage

This project provides a command line interface (CLI) for encoding and decoding files using arithmetic coding.

To use the CLI:

1. **Compile the project:**
```bash
go build
```

2. **Run the CLI with the desired options:**

```bash
./bramah-codec -e -i <input_file> -o <output_file>
```
or
```bash
./bramah-codec -d -i <input_file> -o <output_file> -s <statistics_file>
```

**Options:**

* **-e | -d:** Encode or decode (mutually exclusive).
* **-f | -a:** Fixed or adaptive model (mutually exclusive).
* **-i <input_file>:** Input file to be encoded or decoded.
* **-o <output_file>:** Output file for the encoded or decoded data.
* **-s <statistics_file>:** Input file containing frequency statistics for the fixed model (only required for decoding).

**Example:**

To encode the file `input.txt` into `output.bin` using the adaptive model:

```bash
./bramah-codec -e -a -i input.txt -o output.bin
```

To decode the file `output.bin` using the fixed model with frequency statistics in `statistics.txt`:

```bash
./bramah-codec -d -f -i output.bin -o output_decoded.txt -s statistics.txt
```

## Further Development

* **File input/output:** Currently, the CLI reads from standard input and writes to standard output. It can be extended to read and write from/to files.
* **Fixed model statistics:**  The CLI currently can create fixed model statistics based on a given input file. Further development could focus on more efficient methods for creating and storing fixed model statistics.
* **Fixed model encoding:**  This could be implemented to enable encoding using fixed model frequencies.

## Documentation

This package contains various functions for performing arithmetic coding. Here's a breakdown of the most important ones:

* `NewModel(class int, input *bufio.Reader)`: Creates a new `Model` object based on either the fixed or adaptive model.
* `Initialize() error`: Initializes the `Model` with frequency statistics and sets up tables for symbol translation.
* `UpdateModel(symbol int)`: Updates the frequency statistics of the `Model` for adaptive models.
* `Encode(input *bufio.Reader, output *bufio.Writer, model *Model)`: Encodes the input stream using the provided `Model`.
* `Decode(input *bufio.Reader, output *bufio.Writer, model *Model)`: Decodes the input stream using the provided `Model`.

For more details, refer to the code comments.