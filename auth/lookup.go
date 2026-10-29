package auth

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/chijiokejoseph/gocart-mock/conv"
	e "github.com/chijiokejoseph/gocart-mock/err"
)

// ErrNoMatch is an error that indicates that a username or email was not found in the log file.
var ErrNoMatch = errors.New("username/email not found")

// Lookup searches for an account in the log file by username or email.
// It scans each line of the log file, splitting the line into its components:
// username, email, hash, and wallet balance. If a match is found for the given
// account (either username or email), it returns the corresponding values.
// If no match is found, it optionally prints a message and returns an error.
//
// Parameters:
//   - account - the username or email to search for
//   - print   - if true, prints a message when no match is found
//
// Returns:
//   - username - the matched username
//   - email    - the matched email
//   - hash     - the matched hash
//   - balance  - the matched wallet balance as a decimal.Decimal
//   - err      - error if no match is found, otherwise nil
func Lookup(account string, print bool) (username string, email string, hash string, balance decimal.Decimal, err error) {
	// Refresh reloads the log file before scanning.
	Refresh()

	// Create a scanner to read the log file line by line.
	scanner := bufio.NewScanner(LogFile)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		// Read the current line from the log file.
		line := scanner.Text()
		// Split the line into its components: username, email, hash, wallet balance.
		parts := strings.Split(line, ",")
		username, email, hash, walletStr := parts[0], parts[1], parts[2], parts[3]
		// Check if the username or email matches the provided account.
		if username == account || email == account {
			// Convert the wallet balance string to decimal and return all values.
			return username, email, hash, e.Must(conv.Str2Dec(walletStr)), nil
		}
	}
	// If no match is found, prepare a message.
	msg := fmt.Sprintf("no username or email matches the value %s", account)
	if print {
		// Optionally print the message if 'print' is true.
		fmt.Printf("\n%s\n", msg)
	}
	// Return zero values and a 'no match' error.
	return "", "", "", decimal.Decimal{}, ErrNoMatch
}
