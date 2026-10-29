package conv_test

import (
	"testing"

	"github.com/chijiokejoseph/gocart-mock/conv"
	"github.com/shopspring/decimal"
)

type TestCase struct {
	sample   decimal.Decimal
	expected string
}

// TestDec2Str tests the Dec2Str function for correct decimal to string conversion with formatting.
//
// Params:
//   - t: *testing.T - the testing context
//
// Returns:
//   - none
func TestDec2Str(t *testing.T) {
	// Define test cases and verify Dec2Str output against expected results.
	testcases := []TestCase{
		{decimal.NewFromFloat(123000.456), "123,000.46"},
		{decimal.NewFromFloat(123456.789), "123,456.79"},
		{decimal.NewFromFloat(1240000.8), "1,240,000.80"},
		{decimal.NewFromFloat(420_000_000), "420,000,000.00"},
	}

	// Iterate over each test case, convert the decimal to string, and compare with the expected result.
	for _, testcase := range testcases {
		actual, err := conv.Dec2Str(testcase.sample, conv.NewOptions().WithPrecision(2))
		if err != nil {
			// Report error if conversion fails.
			t.Errorf("\n%s\n", err)
		}
		if actual != testcase.expected {
			// Report mismatch if the actual output does not match the expected string.
			t.Errorf("\nExpected %s, got %s\n", testcase.expected, actual)
		}
	}
}

// TestDec2Str2 tests the Dec2Str function for correct formatting of decimal numbers
// with zero decimal precision (no fractional part).
//
// Params:
//   - t: *testing.T
//     The testing context provided by the Go testing framework.
//
// Returns:
//   - none
func TestDec2Str2(t *testing.T) {
	// Define test cases with sample decimal values and their expected string representations.
	testcases := []TestCase{
		{decimal.NewFromFloat(123000.4), "123,000"},
		{decimal.NewFromFloat(123456.2), "123,456"},
		{decimal.NewFromFloat(1240000.0), "1,240,000"},
		{decimal.NewFromFloat(420_000_000), "420,000,000"},
	}

	// Iterate over each test case, convert the decimal to a string with zero precision,
	// and compare the result with the expected output.
	for _, testcase := range testcases {
		actual, err := conv.Dec2Str(testcase.sample, conv.NewOptions().WithPrecision(0))
		if err != nil {
			// Report an error if the conversion fails.
			t.Errorf("\n%s\n", err)
		}
		if actual != testcase.expected {
			// Report a mismatch if the actual output does not match the expected string.
			t.Errorf("\nExpected %s, got %s\n", testcase.expected, actual)
		}
	}
}
