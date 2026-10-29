package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/chijiokejoseph/gocart-mock/paths"
)

// accountLog is the path to the account log file
var accountLog = fmt.Sprintf("%s/accounts.log", paths.DataPath())

// accessLog opens the account log file for reading and writing.
//
// Params:
//   - none
//
// Returns:
//   - *os.File: a pointer to the opened log file
func accessLog() *os.File {
	// Attempt to open (or create if it doesn't exist) the account log file with read/write permissions
	file, err := os.OpenFile(accountLog, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		// Log a fatal error and terminate the program if the file cannot be opened
		log.Fatalln(err)
	}
	// Return the opened file pointer
	return file
}

// LogFile is a global variable that holds the log file
var LogFile = accessLog()

// Refresh closes the current log file and reopens it to ensure the file handle is up-to-date.
//
// Params:
//   - none
//
// Returns:
//   - none
func Refresh() {
	// Attempt to close the current log file.
	if err := LogFile.Close(); err != nil {
		// Log a fatal error and terminate the program if closing fails.
		log.Fatal(err)
	}
	// Reopen the log file and update the global LogFile variable.
	LogFile = accessLog()
}

// generatePassword generates a new secure password and its encoded hash.
//
// Params:
//   - none
//
// Returns:
//   - password secret: the newly generated password as a secret type
//   - encodedHash secret: the encoded hash of the password, suitable for storage or verification
func generatePassword() (password secret, encodedHash secret) {
	// Generate a new secure password that meets the required criteria.
	password = newPassword()
	// Generate an encoded hash (salt + hash, base64-encoded) for the password.
	encodedHash = newEncodedHash(password)
	// Return both the password and its encoded hash.
	return
}

// fromPasswordString validates a raw password string and returns the password and its encoded hash.
//
// Params:
//   - rawPassword string: the plain-text password input to be validated and encoded
//
// Returns:
//   - password secret: the validated password as a secret type
//   - encodedHash secret: the encoded hash (salt + hash, base64-encoded) of the password
//   - err error: an error if validation fails; otherwise, nil
func fromPasswordString(rawPassword string) (password secret, encodedHash secret, err error) {
	// Validate the raw password string according to password policy.
	validPass, err := validatePassword(rawPassword)
	if err != nil {
		// Return an error if the password is invalid.
		return "", "", err
	}
	// Generate an encoded hash (salt + hash, base64-encoded) for the validated password.
	encodedHash = newEncodedHash(validPass)
	// Return the validated password and its encoded hash.
	return validPass, encodedHash, nil
}

// ErrLessChars depicts that the password is not long enough
var ErrLessChars = errors.New("password must be at least 16 characters long")

// ErrNoLowercase depicts that the password does not contain a lowercase letter
var ErrNoLowercase = errors.New("password must contain at least one lowercase letter")

// ErrNoUppercase depicts that the password does not contain an uppercase letter
var ErrNoUppercase = errors.New("password must contain at least one uppercase letter")

// ErrNoDigit depicts that the password does not contain a number
var ErrNoDigit = errors.New("password must contain at least one number")

// ErrNoSymbol depicts that the password does not contain a special character
var ErrNoSymbol = errors.New("password must contain at least one special character")

// validatePassword checks whether a password string meets the required security criteria.
//
// The password must:
//   - Be at least NUM characters long
//   - Contain at least one lowercase letter
//   - Contain at least one uppercase letter
//   - Contain at least one digit
//   - Contain at least one special character
//
// Params:
//   - password string: the plain-text password to validate
//
// Returns:
//   - secret: the validated password as a secret type (if valid)
//   - error: an error describing the first failed requirement, or nil if valid
func validatePassword(password string) (secret, error) {
	// Check for minimum length requirement.
	if len(password) < NUM {
		return "", ErrLessChars
	}

	var (
		countSmall int // Number of lowercase letters
		countCaps  int // Number of uppercase letters
		countNum   int // Number of digits
		countSpec  int // Number of special characters
	)

	// Iterate over each character and count types.
	for _, r := range password {
		if slices.Contains([]rune(SMALL), r) {
			countSmall++
		} else if slices.Contains([]rune(CAPS), r) {
			countCaps++
		} else if slices.Contains([]rune(NUMBERS), r) {
			countNum++
		} else if slices.Contains([]rune(SPECIAL), r) {
			countSpec++
		}
	}

	// Check for at least one lowercase letter.
	if countSmall == 0 {
		return "", ErrNoLowercase
	}
	// Check for at least one uppercase letter.
	if countCaps == 0 {
		return "", ErrNoUppercase
	}
	// Check for at least one digit.
	if countNum == 0 {
		return "", ErrNoDigit
	}
	// Check for at least one special character.
	if countSpec == 0 {
		return "", ErrNoSymbol
	}

	// All requirements met; return the password as a secret.
	return secret(password), nil
}

var ErrBadPassword = errors.New("invalid password entered")
var ErrCorruptHash = errors.New("user account hash is corrupted. Unable to sign in")

// validatePasswordByHash verifies that a plain-text password matches a stored encoded hash.
//
// This function decodes the stored hash to extract the salt and hash bytes, then re-computes
// the hash for the provided password using the extracted salt. It compares the computed hash
// with the stored hash in a constant-time manner to prevent timing attacks.
//
// Params:
//   - password string: the plain-text password to validate
//   - storedHash string: the encoded hash (salt + hash, base64-encoded) to validate against
//
// Returns:
//   - error: ErrCorruptHash if the stored hash is invalid or corrupted,
//     ErrBadPassword if the password does not match the hash,
//     or nil if the password is valid
func validatePasswordByHash(password string, storedHash string) error {
	// Decode the stored hash to extract the salt and hash bytes.
	storedSaltBytes, storedHashBytes, err := newDecodedHash(secret(storedHash))
	if err != nil {
		// Print an error message and return ErrCorruptHash if decoding fails.
		fmt.Printf("\n%s\n", ErrCorruptHash)
		return ErrCorruptHash
	}

	// Compute the hash for the provided password using the extracted salt.
	computedHash := reprHash(password, storedSaltBytes)
	computedHashBytes := []byte(computedHash)

	// Check if the computed hash length matches the stored hash length.
	if len(computedHashBytes) != len(storedHashBytes) {
		fmt.Println("\nWrong password entered")
		return ErrBadPassword
	}

	// Compare the computed hash and stored hash byte-by-byte.
	for i := range computedHashBytes {
		if computedHashBytes[i] != storedHashBytes[i] {
			fmt.Println("\nWrong password entered")
			return ErrBadPassword
		}
	}

	// Password is valid; return nil.
	return nil
}
