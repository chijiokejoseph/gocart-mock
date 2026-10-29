package inputs

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

// ErrParseNum is an error depicting that user input is not parseable to an integer.
var ErrParseNum = errors.New("error in InputNum parsing user input to an integer number")

// ErrParseFloat is an error depicting that user input is not parseable to a float
var ErrParseFloat = errors.New("error in InputFloat parsing user input to a float")

// InputNum prompts the user for input and parses it as an integer.
//
// Params:
//   - reader: *bufio.Reader -
//     The reader to read user input from.
//   - prompt: string -
//     The prompt message to display to the user.
//
// Returns:
//   - int:
//     The parsed integer value from user input.
//   - error:
//     nil if the user input is successfully parsed as an integer else error encountered
//
// An error is returned if:
//   - reading input fails (e.g., I/O error),
//   - the input cannot be parsed as an integer (ErrParseNum),
//   - or the user cancels the operation by entering "quit" or "exit" (see ErrQuit from InputStr).
//
// Notes:
//   - If the user enters "quit" or "exit", InputStr returns a ErrQuit error, which is propagated by this function.
//   - If the input is not a valid integer, a wrapped ErrParseNum is returned.
func InputNum(reader *bufio.Reader, prompt string) (int, error) {
	// Get user input as a string using InputStr.
	input, err := InputStr(reader, prompt)
	if err != nil {
		// Propagate error if reading input fails or user cancels (ErrQuit).
		return 0, err
	}
	var num int
	// Try to parse the input string as an integer.
	_, err = fmt.Sscanf(input, "%d", &num)
	if err != nil {
		// Print error message if input is not an integer.
		fmt.Printf("\nUser input '%s' is not an integer\n", input)
		// Return a wrapped error indicating parsing failure.
		return 0, fmt.Errorf("%w. Error: %s", ErrParseNum, err)
	}
	// Return the parsed integer and nil error.
	return num, nil
}

// InputDecimal prompts the user for input and parses it as a decimal.Decimal.
//
// Params:
//   - reader: *bufio.Reader -
//     The reader to read user input from.
//   - prompt: string -
//     The prompt message to display to the user.
//
// Returns:
//   - decimal.Decimal -
//     The parsed decimal value from user input.
//   - error -
//     nil if the user input is successfully parsed as a decimal else error encountered
//
// An error is returned if:
//   - reading input fails (e.g., I/O error),
//   - the input cannot be parsed as a float (ParseFloatErr),
//   - or the user cancels the operation by entering "quit" or "exit" (see ErrQuit from InputStr/InputFloat).
//
// Notes:
//   - This function builds on top of InputFloat and propagates all errors from it, including user cancellation.
//   - If the input is not a valid float, a wrapped ParseFloatErr is returned.
func InputDecimal(reader *bufio.Reader, prompt string) (decimal.Decimal, error) {
	// Get user input as a float64 using InputFloat.
	input, err := InputFloat(reader, prompt)
	if err != nil {
		// Propagate error if reading input fails or user cancels (ErrQuit).
		return decimal.Decimal{}, err
	}
	// Convert the float64 input to a decimal.Decimal.
	num := decimal.NewFromFloat(input)
	return num, nil
}

// InputFloat prompts the user for input and parses it as a float64.
//
// Params:
//   - reader: *bufio.Reader -
//     The reader to read user input from.
//   - prompt: string -
//     The prompt message to display to the user.
//
// Returns:
//   - float64 -
//     The parsed floating-point value from user input.
//   - error -
//     nil if the user input is successfully parsed as a float else error encountered
//
// An error is returned if:
//   - reading input fails (e.g., I/O error),
//   - the input cannot be parsed as a float (ParseFloatErr),
//   - or the user cancels the operation by entering "quit" or "exit" (see ErrQuit from InputStr).
//
// Notes:
//   - If the user enters "quit" or "exit", InputStr returns a ErrQuit error, which is propagated by this function.
//   - If the input is not a valid float, a wrapped ParseFloatErr is returned.
func InputFloat(reader *bufio.Reader, prompt string) (float64, error) {
	// Get user input as a string using InputStr.
	input, err := InputStr(reader, prompt)
	if err != nil {
		// Propagate error if reading input fails or user cancels (ErrQuit).
		return 0, err
	}
	var num float64
	// Try to parse the input string as a float64.
	_, err = fmt.Sscanf(input, "%f", &num)
	if err != nil {
		// Print error message if input is not a float.
		fmt.Printf("\nUser input '%s' is not a float\n", input)
		// Return a wrapped error indicating parsing failure.
		return 0, fmt.Errorf("%w. Error: %s", ErrParseFloat, err)
	}
	// Return the parsed float and nil error.
	return num, nil
}
