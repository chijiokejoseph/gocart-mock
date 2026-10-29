package auth

import (
	"bufio"
	"strings"
	"testing"
)

// TestEnterPassword tests the enterPassword and validatePasswordByHash functions.
// 
// Params:
//   - t: *testing.T - the testing context provided by the Go testing framework.
func TestEnterPassword(t *testing.T) {
	// Ensure teardown is called after the test completes.
	t.Cleanup(func() {
		teardown()
	})

	// Define sample passwords for testing.
	samples := []string{
		"abdefiose8932L_9wls",
	}

	// Create a reader simulating user input for the password.
	reader := bufio.NewReader(strings.NewReader(samples[0] + "\n"))

	// Call enterPassword to get the password and its hash.
	password, hash, err := enterPassword(reader)
	if err != nil {
		t.Errorf("%s", err)
	}

	// Validate the password against its hash.
	err = validatePasswordByHash(string(password), string(hash))
	if err != nil {
		t.Errorf("%s", err)
	}
}
