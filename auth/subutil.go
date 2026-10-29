package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	e "github.com/chijiokejoseph/gocart-mock/err"
	"golang.org/x/crypto/argon2"
)

const (
	SMALL     = "abcdefghijiklmnopqrstuvwxyz"    // lowercase ascii
	CAPS      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"     // uppercase ascii
	NUMBERS   = "0123456789"                     // digits
	SPECIAL   = "#!@#$%^&*()_+"                  // special characters
	LETTERS   = SMALL + CAPS + NUMBERS + SPECIAL // all characters
	NUM       = 16                               // password length
	EncodeFmt = "%s:%s"                          // format for encoded salt and hash
)

type secret string

// newPassword generates a new valid password that meets the required criteria.
//
// Parameters:
//   - None
//
// Returns:
//   - secret: a valid password of type secret
//
// The function repeatedly generates random passwords using the `LETTERS` character set
// until one passes validation. It uses cryptographically secure random number generation
// and validates the password using `validatePassword`.
func newPassword() secret {
	// Allocate a byte slice to hold the password characters
	password := make([]byte, NUM)
	var validPass secret
	for {
		// Fill the password slice with random characters from LETTERS
		for i := 0; i < NUM; i++ {
			// Generate a secure random index into LETTERS
			n := e.Must(rand.Int(rand.Reader, big.NewInt(int64(len(LETTERS)))))
			password[i] = LETTERS[n.Int64()]
		}
		// Validate the generated password
		if password, err := validatePassword(string(password)); err == nil {
			// If valid, assign to validPass and exit loop
			validPass = password
			break
		}
		// If not valid, repeat the process
	}
	// Return the valid password
	return validPass
}

// newHash generates a cryptographic salt and hashes the given password using Argon2id.
//
// Parameters:
//   - password secret: the password to be hashed, provided as a secret type
//
// Returns:
//   - salt secret: the randomly generated salt used for hashing, as a secret type
//   - hash secret: the resulting Argon2id hash of the password and salt, as a secret type
//
// The function creates a 16-byte random salt, then uses Argon2id with specific parameters
// to hash the password. Both the salt and hash are returned as secret types.
func newHash(password secret) (salt secret, hash secret) {
	// Allocate a 16-byte slice for the salt
	saltStr := make([]byte, 16)
	// Fill the salt slice with cryptographically secure random bytes
	_ = e.Must(rand.Read(saltStr))
	// Hash the password using Argon2id with the generated salt and fixed parameters
	hashStr := argon2.IDKey([]byte(password), saltStr, 1, 64*1024, 4, 32)
	// Convert the salt and hash to the secret type
	salt = secret(saltStr)
	hash = secret(hashStr)
	// Return the salt and hash
	return
}

// reprHash hashes the given password with the provided salt using Argon2id.
//
// Parameters:
//   - password string: the plaintext password to be hashed
//   - salt []byte: the cryptographic salt to use in the hashing process
//
// Returns:
//   - secret: the resulting Argon2id hash, converted to the secret type
//
// The function uses Argon2id with fixed parameters (1 iteration, 64 MB memory, 4 threads, 32-byte output)
// to securely hash the password and salt.
func reprHash(password string, salt []byte) secret {
	// Hash the password using Argon2id with the provided salt and fixed parameters
	hashStr := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	// Convert the hash to the secret type and return
	return secret(hashStr)
}

// encodeSecret encodes the given secret value into a base64-encoded string.
//
// Parameters:
//   - key secret: the secret value to be encoded
//
// Returns:
//   - secret: the base64-encoded representation of the input secret
//
// The function converts the secret to a byte slice, encodes it using base64 encoding,
// and returns the result as a secret type.
func encodeSecret(key secret) secret {
	// Convert the secret to a byte slice and encode it as a base64 string
	hashEncoded := base64.StdEncoding.EncodeToString([]byte(key))
	// Return the encoded string as a secret type
	return secret(hashEncoded)
}

// decodeSecret decodes a base64-encoded secret string back into its original byte slice.
//
// Parameters:
//   - key secret: the base64-encoded secret value to be decoded
//
// Returns:
//   - []byte: the decoded byte slice representing the original secret
//   - error: an error if the decoding fails, otherwise nil
//
// The function converts the secret to a string, decodes it from base64 encoding,
// and returns the resulting bytes. If decoding fails, an error is returned.
func decodeSecret(key secret) ([]byte, error) {
	// Convert the secret to a string and decode it from base64
	return base64.StdEncoding.DecodeString(string(key))
}

// newEncodedHash generates a new encoded hash string for the given password.
//
// The function creates a random salt and hashes the password using newHash,
// then encodes both the salt and hash using base64 encoding. The encoded
// values are concatenated using a format string (EncodeFmt) and returned
// as a single secret value.
//
// Parameters:
//   - password secret: the password to be hashed and encoded
//
// Returns:
//   - secret: a base64-encoded string containing the salt and hash, formatted as specified by EncodeFmt
//
// The resulting encoded string can be safely stored or transmitted, and later
// decoded and verified using the corresponding decoding and verification functions.
func newEncodedHash(password secret) secret {
	// Generate a random salt and hash the password
	salt, hash := newHash(password)
	// Encode the salt and hash using base64 encoding
	encodedSalt := encodeSecret(salt)
	encodedHash := encodeSecret(hash)
	// Format and return the encoded salt and hash as a single secret string
	return secret(fmt.Sprintf(EncodeFmt, encodedSalt, encodedHash))
}

// newDecodedHash decodes an encoded password hash string into its salt and hash byte slices.
//
// The function expects the input to be a colon-separated string containing base64-encoded salt and hash values.
// It splits the input, decodes each part from base64, and returns the resulting byte slices.
//
// Parameters:
//   - key secret: the encoded key hash string, formatted as "encodedSalt:encodedHash"
//
// Returns:
//   - saltBytes []byte: the decoded salt as a byte slice
//   - hashBytes []byte: the decoded hash as a byte slice
//   - err error: an error if the input is invalid or decoding fails; otherwise, nil
//
// Example usage:
//
//	salt, hash, err := newDecodedHash(encodedKey)
//	if err != nil {
//	    // handle error
//	}
func newDecodedHash(key secret) (saltBytes []byte, hashBytes []byte, err error) {
	// Check if the key contains the expected separator
	if !strings.Contains(string(key), ":") {
		return nil, nil, fmt.Errorf("invalid encoded key hash")
	}
	// Split the key into salt and hash parts
	parts := strings.Split(string(key), ":")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid encoded key hash")
	}
	encodedSalt, encodedHash := parts[0], parts[1]
	// Decode the salt from base64
	saltBytes, err = decodeSecret(secret(encodedSalt))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid salt in encoded key hash")
	}
	// Decode the hash from base64
	hashBytes, err = decodeSecret(secret(encodedHash))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid hash in encoded key hash")
	}
	// Return the decoded salt and hash
	return saltBytes, hashBytes, nil
}
