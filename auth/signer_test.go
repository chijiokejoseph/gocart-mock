package auth

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/chijiokejoseph/gocart-mock/wallets"
)

// setup initializes test data for unit tests.
//
// This helper function creates a new Signer instance with a fixed username and email,
// generates a password and its hash, logs the setup status, and returns all three values.
//
// Parameters:
//   - t (*testing.T): The testing context used for logging and test control.
//
// Returns:
//   - *Signer: A pointer to the initialized Signer instance.
//   - secret: The generated password.
//   - secret: The hash of the generated password.
func setup(t *testing.T) (*Signer, secret, secret) {
	// Define a fixed username for the test Signer.
	username := "smithy7"
	// Define a fixed email for the test Signer.
	email := "smithy7@gmail.com"
	// Generate a new password and its hash for the test Signer.
	password, hash := generatePassword()
	// Log a message indicating that setup was successful.
	t.Log("setup successful")
	// Create and return a new Signer, along with the password and hash.
	return newSigner(username, email), password, hash
}

// TestCase represents a test case for the SaveSecretAndLookup function.
// Fields:
// 	- name: The name of the test case.
// 	- setup: A function that sets up the test case.
type TestCase struct {
	name  string
	setup func(t *testing.T) (*Signer, secret, secret)
}

// TestSaveSecretAndLookup verifies that a signer's secret (password hash) can be saved and
// subsequently retrieved and validated from the log file.
//
// This test performs the following steps:
//   - Sets up test data and ensures cleanup after the test.
//   - Saves a generated password hash for a test signer.
//   - Refreshes the log file and prints its contents for inspection.
//   - Looks up the saved hash in the log file and prints it.
//   - Validates that the original password matches the retrieved hash.
//
// Parameters:
//   - t (*testing.T): The testing context used for assertions and test control.
//
// Returns:
//   - None. Reports errors via the testing framework.
func TestSaveSecretAndLookup(t *testing.T) {
	// Ensure any necessary teardown is performed after the test.
	t.Cleanup(func() {
		teardown()
	})

	// Initialize the test case with a setup function.
	testcase := TestCase{
		name:  "",
		setup: setup,
	}

	// Run the setup to get a test signer, password, and hash.
	signer, password, hash := testcase.setup(t)
	fmt.Println("hash: ", hash)

	// Save the generated hash for the signer.
	err := signer.saveSecret(hash)
	if err != nil {
		// If a system write or truncate error occurs, log fatally and stop the test.
		if errors.Is(err, ErrSysWrite) || errors.Is(err, ErrTruncate) {
			log.Fatalln(err)
		}
	}

	// Refresh the log file to ensure it's up to date.
	Refresh()
	file := LogFile
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Print each line in the log file for debugging purposes.
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("line = ", line)
	}

	// Look up the hash for the signer in the log file.
	_, _, hashTxt, _, err := Lookup(signer.Username(), true)
	fmt.Println("hashTxt: ", hashTxt)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	// Validate that the original password matches the retrieved hash.
	err = validatePasswordByHash(string(password), hashTxt)
	if err != nil {
		t.Errorf("%s\n", err)
	}
}

// TestSignerFromUser verifies that a User interface value can be correctly
// converted back to a Signer using the SignerFromUser function.
//
// This test performs the following steps:
//   - Registers a cleanup function to run after the test.
//   - Sets up a test case and initializes a Signer.
//   - Converts the Signer to a User interface.
//   - Attempts to recover the Signer from the User using SignerFromUser.
//   - Asserts that the recovered Signer has the expected username.
//
// Parameters:
//   - t (*testing.T): The testing context used for assertions and test control.
//
// Returns:
//   - None. Reports errors via the testing framework.
func TestSignerFromUser(t *testing.T) {
	// Register a cleanup function to run after the test completes.
	t.Cleanup(func() {
		teardown()
	})

	// Initialize the test case with a setup function.
	testcase := TestCase{
		name:  "",
		setup: setup,
	}

	// Run the setup to get a test signer.
	signer, _, _ := testcase.setup(t)

	// Assign the signer to a variable of type User interface.
	var user User = signer

	// Attempt to recover the Signer from the User interface.
	altSigner, err := SignerFromUser(user)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	// Assert that the recovered Signer has the expected username.
	if altSigner.Username() != signer.Username() {
		t.Errorf("expected Signer Username: %s, Got Username: %s\n", signer.Username(), altSigner.Username())
	}
}

// TestSigner verifies that funding a signer's wallet updates the balance as expected.
//
// This test performs the following steps:
//   - Registers a cleanup function to run after the test.
//   - Sets up a test case and initializes a Signer.
//   - Funds the signer's wallet with a predefined amount (wallets.TEN).
//   - Checks that the wallet balance matches the expected value after funding.
//   - Reports an error if the actual balance does not match the expected balance.
//
// Parameters:
//   - t (*testing.T): The testing context used for assertions and test control.
//
// Returns:
//   - None. Reports errors via the testing framework.
func TestSigner(t *testing.T) {
	// Register a cleanup function to run after the test completes.
	t.Cleanup(func() {
		teardown()
	})

	// Initialize the test case with a setup function.
	testcase := TestCase{
		name:  "",
		setup: setup,
	}

	// Run the setup to get a test signer.
	signer, _, _ := testcase.setup(t)

	// Get the signer's wallet.
	w := signer.Wallet()

	// Fund the wallet with a predefined amount (wallets.TEN).
	w.Fund(wallets.TEN)

	// Retrieve the actual balance from the signer's wallet.
	actual := signer.Wallet().Balance().IntPart()
	// Retrieve the expected balance from the funded wallet.
	expected := w.Balance().IntPart()

	// Assert that the actual balance matches the expected balance.
	if actual != expected {
		t.Errorf("Expected Wallet Balance of %d, Got Signer Wallet Balance of %d\n", expected, actual)
	}
}

// TestSignerDelete verifies that deleting a signer removes their record from the log file
// and that subsequent lookups for the deleted signer fail as expected.
//
// This test performs the following steps:
//   - Registers a cleanup function to run after the test.
//   - Creates two test signers, funds their wallets, and saves their password hashes.
//   - Deletes the first signer and checks that the deletion was successful.
//   - Refreshes the log file and logs its contents for inspection.
//   - Attempts to look up the deleted signer and asserts that the correct error is returned.
//
// Parameters:
//   - t (*testing.T): The testing context used for assertions and test control.
//
// Returns:
//   - None. Reports errors via the testing framework.
func TestSignerDelete(t *testing.T) {
	// Register a cleanup function to run after the test completes.
	t.Cleanup(func() {
		teardown()
	})

	// Create a new signer with username "smithy7".
	signer := newSigner("smithy7", "")
	// Generate a password and hash for the signer.
	_, hash := generatePassword()
	// Fund the signer's wallet with a predefined amount.
	signer.Wallet().Fund(wallets.TEN)
	// Save the generated hash for the signer.
	err := signer.saveSecret(hash)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	// Create an alternate signer with username "smithy8".
	altSigner := newSigner("smithy8", "")
	// Generate a password and hash for the alternate signer.
	_, altHash := generatePassword()
	// Fund the alternate signer's wallet.
	altSigner.Wallet().Fund(wallets.FIVE)
	// Save the generated hash for the alternate signer.
	err = altSigner.saveSecret(altHash)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	// Delete the first signer and check for errors.
	err = signer.Delete()
	if err != nil {
		t.Errorf("%s\n", err)
	}

	// Refresh the log file to ensure it's up to date.
	Refresh()
	// Scan and log each line in the log file for debugging purposes.
	scanner := bufio.NewScanner(LogFile)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
	for scanner.Scan() {
		t.Logf(scanner.Text())
	}

	// Attempt to look up the deleted signer; expect an ErrNoMatch error.
	_, _, _, _, err = Lookup(signer.Username(), true)
	actual := err
	expected := ErrNoMatch
	if !errors.Is(actual, expected) {
		t.Errorf("Expected error: %s, Got error: %s\n", expected, actual)
	}
}
