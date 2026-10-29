package conv

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/shopspring/decimal"
)


// Float2Dec converts a float64 number to a decimal.Decimal type.
//
// This function takes a float64 value as input and returns its equivalent
// representation as a decimal.Decimal, which provides arbitrary-precision
// decimal arithmetic.
//
// Parameters:
//   - num: The float64 number to be converted.
//
// Returns:
//   - decimal.Decimal: The decimal representation of the input float64 number.
func Float2Dec(num float64) decimal.Decimal {
	// Convert the float64 number to decimal.Decimal using NewFromFloat.
	return decimal.NewFromFloat(num)
}

// Int2Dec converts an int to a decimal.Decimal type.
//
// This function takes an int value as input and returns its equivalent
// representation as a decimal.Decimal, which provides arbitrary-precision
// decimal arithmetic.
//
// Parameters:
//   - num: The int number to be converted.
//
// Returns:
//   - decimal.Decimal: The decimal representation of the input int number.
func Int2Dec(num int) decimal.Decimal {
	return decimal.NewFromInt(int64(num))
}

// Str2Dec converts a string representation of a number to a decimal.Decimal type.
//
// This function takes a string as input and returns its equivalent representation
// as a decimal.Decimal, which provides arbitrary-precision decimal arithmetic.
//
// Parameters:
//   - num: The string number to be converted.
//
// Returns:
//   - decimal.Decimal: The decimal representation of the input string number.
//   - error: An error if the conversion is not successful.
func Str2Dec(num string) (decimal.Decimal, error) {
	dec, err := decimal.NewFromString(num)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return dec, nil
}

// Options is a struct that holds the options for the Dec2Str function.
type Options struct {
	Precision int
}

// WithPrecision sets the precision for the Options instance.
//
// This method updates the Precision field of the Options instance
// with the provided precision value and returns the updated instance.
//
// Parameters:
//   - precision: The number of decimal places to be used for formatting.
//
// Returns:
//   - Options: The updated Options instance with the specified precision set.

func (co Options) WithPrecision(precision int) Options {
	co.Precision = precision
	return co
}

// NewOptions returns a new Options instance with default values.
//
// The default precision is set to 0.
func NewOptions() Options {
	return Options{}
}

// Dec2StrFixed is a variant of Dec2Str that takes a decimal.Decimal number
// and returns its string representation with the default precision of 2.
//
// The default options are specified as follows:
//   - Precision: 2
//
// Parameters:
//   - num: The decimal.Decimal number to be converted to a string.
//
// Returns:
//   - string: The string representation of the input decimal.Decimal number.
//   - error: An error if the conversion is not successful.
func Dec2StrFixed(num decimal.Decimal) (string, error) {
	return Dec2Str(num, NewOptions().WithPrecision(2))
}


// Dec2Str converts a decimal.Decimal number to a formatted string representation
// with thousands separators and a fixed number of decimal places as specified in opts.Precision.
// It returns the formatted string and an error if any occurs during formatting.
//
// Parameters:
//   - num: The decimal.Decimal number to be formatted.
//   - opts: Options struct containing formatting options (e.g., Precision).
//
// Returns:
//   - string: The formatted decimal number as a string with thousands separators.
//   - error: An error if formatting fails, otherwise nil.
func Dec2Str(num decimal.Decimal, opts Options) (string, error) {
	// Convert the decimal to a string with fixed precision
	decimalStr := num.StringFixed(int32(opts.Precision))
	runes := []rune(decimalStr)
	var (
		intPart     string // Holds the integer part of the number
		decimalPart string // Holds the decimal part of the number
	)

	// Split the string into integer and decimal parts
	if slices.Contains(runes, '.') {
		idx := slices.Index(runes, '.')
		intPart = decimalStr[:idx]
		decimalPart = decimalStr[idx+1:]
	} else {
		intPart = decimalStr
	}

	// Reverse the integer part to facilitate grouping by thousands
	intRunes := []rune(intPart)
	slices.Reverse(intRunes)
	intPart = string(intRunes)

	sepDigits := make([]string, 0) // Holds the groups of digits separated by commas
	rem := len(intRunes) % 3       // Remainder for digits not forming a complete group
	n := math.Floor(float64(len(intRunes)) / 3) // Number of complete groups of 3 digits

	// Group digits into thousands
	for i := 0; i < int(n); i++ {
		if n > 0 {
			start := i * 3
			end := (i + 1) * 3
			sepDigits = append(sepDigits, intPart[start:end])
			continue
		}
	}

	// Add any remaining digits that don't form a complete group
	if rem != 0 {
		start := len(intRunes) - rem
		sepDigits = append(sepDigits, intPart[start:])
	}

	// Join the groups with commas
	intPart = strings.Join(sepDigits, ",")

	// Reverse the integer part back to its original order
	intRunes = []rune(intPart)
	slices.Reverse(intRunes)
	intPart = string(intRunes)

	// Combine the integer and decimal parts
	decimalStr = fmt.Sprintf("%s.%s", intPart, decimalPart)

	// Remove the trailing dot if there is no decimal part
	if decimalPart == "" {
		decimalStr = strings.TrimSuffix(decimalStr, ".")
	}

	return decimalStr, nil
}
