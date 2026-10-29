package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/chijiokejoseph/gocart-mock/paths"
)

// TestLogFile tests that the LogFile variable is not nil.
//
// Params:
//   - t: *testing.T - the testing context provided by the Go testing framework.
func TestLogFile(t *testing.T) {
	// Ensure resources are cleaned up after the test completes.
	defer teardown()

	// Retrieve the LogFile variable to check its value.
	file := LogFile

	// If LogFile is nil, report an error in the test.
	if file == nil {
		t.Error("LogFile is nil")
	}
}

// teardown cleans up resources used during tests or application execution.
//
// Parameters:
//   - none
//
// Returns:
//   - none
func teardown() {
	// Retrieve the global log file to be closed.
	file := LogFile

	// Attempt to close the log file.
	// If closing fails, log the error and terminate the program.
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	// Attempt to remove all files and directories at the application's data path.
	// If removal fails, log the error and terminate the program.
	if err := os.RemoveAll(paths.DataPath()); err != nil {
		log.Fatal(err)
	}
}

// TestPassword verifies the password generation and validation logic.
//
// This test performs the following steps:
//   - Generates a new password and prints it.
//   - Checks that the generated password has the correct length.
//   - Iterates over a set of sample passwords, each designed to fail a specific validation rule.
//   - For each sample, checks that the error returned by fromPasswordString matches the expected error.
//
// Parameters:
//   - t: *testing.T - the testing context provided by the Go testing framework.
//
// Returns:
//   - none
func TestPassword(t *testing.T) {
	// Ensure resources are cleaned up after the test completes.
	defer teardown()

	// Generate a new password and print it for debugging.
	pass := newPassword()
	fmt.Printf("password = %s\n", pass)

	// Check that the generated password has the expected length.
	if len(pass) != NUM {
		t.Error("Password length is wrong")
	}

	// Define sample passwords and the expected validation errors for each.
	sample := []string{
		"1237bw%oo0lb1z",                    // Too short
		"abcdebw%ooolbzaZ-",                 // No number
		"abcdefghijk78l%-9%",                // No capital letter
		"ABCDEFJ88_H.XMII9%#77",             // No lowercase letter
		"ABCDEFGHIJKLMNOPQURST7788abyced99", // No special character
		"4abcide3sf88l44.sb.x",              // Valid password
	}
	expected := []error{
		ErrLessChars,
		ErrNoDigit,
		ErrNoUppercase,
		ErrNoLowercase,
		ErrNoSymbol,
		nil,
	}

	// Iterate over each sample password and check the validation error.
	for i, aCase := range sample {
		_, _, err := fromPasswordString(aCase)
		expectedErr := expected[i]
		if expectedErr != nil {
			// If the error does not match the expected error, report a test failure.
			if !errors.Is(err, expectedErr) {
				t.Errorf("Expected error %s but got %s for password %s\n", expectedErr, err, aCase)
			}
		}
	}
}

// TestNewHash verifies that a new password can be generated and encoded as a hash.
//
// This test performs the following steps:
//   - Generates a new password using newPassword().
//   - Encodes the password into a hash using newEncodedHash().
//   - Prints the resulting hash for inspection.
//
// Parameters:
//   - t: *testing.T - the testing context provided by the Go testing framework.
//
// Returns:
//   - none
func TestNewHash(t *testing.T) {
	// Ensure resources are cleaned up after the test completes.
	defer teardown()

	// Generate a new password.
	pass := newPassword()

	// Encode the password into a hash.
	hash := newEncodedHash(pass)

	// Print the resulting hash for debugging or inspection.
	fmt.Printf("hash = %s\n", hash)
}

// TestValidatePasswordByHash verifies that a generated password can be validated against its hash.
//
// This test performs the following steps:
//   - Registers a cleanup function to run after the test completes.
//   - Generates a new password and its corresponding hash.
//   - Validates that the password matches the generated hash using validatePasswordByHash.
//   - Reports an error if validation fails.
//
// Parameters:
//   - t: *testing.T - the testing context provided by the Go testing framework.
//
// Returns:
//   - none
func TestValidatePasswordByHash(t *testing.T) {
	// Register the teardown function to clean up resources after the test.
	t.Cleanup(teardown)

	// Generate a new password and its corresponding hash.
	password, hash := generatePassword()

	// Validate that the password matches the generated hash.
	err := validatePasswordByHash(string(password), string(hash))
	if err != nil {
		// Report an error if validation fails.
		t.Errorf("%s", err)
	}
}
