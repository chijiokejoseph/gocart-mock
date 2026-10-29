package auth

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/chijiokejoseph/gocart-mock/store"
	"github.com/chijiokejoseph/gocart-mock/wallets"
)

const (
	// saveFMT is the format string for saving a signer to the log file
	saveFMT = "%s,%s,%s,%f\n"
)

// ErrMaxAccount depicts that the number of local accounts stored
// in the program has been exceeded.
var ErrMaxAccount = fmt.Errorf("maximum number of accounts reached")

// writeSave writes the provided lines to the given file after processing.
// It removes empty lines, sorts the lines, truncates the file, and writes the sorted lines.
// Returns an error if any file operation fails.
//
// Parameters:
//
//   - lines - a slice of strings to be written to the file
//   - file  - a pointer to an os.File where the lines will be written
//
// Returns:
//
//   - error - TruncateErr if truncation or seeking fails, SysWriteErr if writing fails, or nil on success
func writeSave(lines []string, file *os.File) error {
	// Iterate over lines and remove empty lines
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			lines = append(lines[i:], lines[i+1:]...)
		}
	}
	// Sort the lines using a custom comparison function
	slices.SortFunc(lines, func(i, j string) int {
		return strings.Compare(i, j)
	})
	// Truncate the file to remove existing content
	err := file.Truncate(0)
	if err != nil {
		return ErrTruncate
	}
	// Seek to the beginning of the file before writing
	_, err = file.Seek(0, 0)
	if err != nil {
		return ErrTruncate
	}
	// Write the sorted lines to the file, joined by newline characters
	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		return ErrSysWrite
	}
	return nil
}

// ErrSysWrite is an error that occurs when a system write operation fails
var ErrSysWrite = fmt.Errorf("system write error. Failed to save password hash")

// ErrTruncate is an error that occurs when a truncate operation fails
var ErrTruncate = fmt.Errorf("truncate error. Failed to truncate accounts log file")

// Signer is a struct that represents a user account
// Fields:
//
//	- username - the username of the user
//	- email    - the email of the user
//	- wallet   - the wallet of the user
//	- cart     - the cart of the user
//	- signed   - whether the user is signed in or not
//	- deleted  - whether the user has been deleted or not
type Signer struct {
	username string
	email    string
	wallet   *wallets.Wallet
	cart     *store.Cart
	signed   bool
	deleted  bool
}

// Username returns the username associated with the Signer.
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - string - the username of the Signer
func (signer *Signer) Username() string {
	return signer.username
}

func (signer *Signer) Email() string {
	return signer.email
}

// Wallet returns the Wallet associated with the Signer.
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - *wallets.Wallet - a pointer to the Signer's Wallet instance
func (signer *Signer) Wallet() *wallets.Wallet {
	return signer.wallet
}

// Cart returns the Cart associated with the Signer.
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - *store.Cart - a pointer to the Signer's Cart instance
func (signer *Signer) Cart() *store.Cart {
	return signer.cart
}

// Verified indicates whether the Signer has been verified.
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - bool - true if the Signer is verified; otherwise, false
func (signer *Signer) Verified() bool {
	return true
}

// Signed returns whether the Signer has completed the signing process.
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - bool - true if the Signer is signed; otherwise, false
func (signer *Signer) Signed() bool {
	return signer.signed
}

// Logout logs the Signer out by saving the current state and marking the Signer as unsigned.
//
// This method first attempts to persist the Signer's state by calling Save().
// If saving fails, it returns the encountered error. Otherwise, it sets the
// signed field to false to indicate the Signer is no longer signed in.
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - error - an error if saving the Signer's state fails; otherwise, nil.
func (signer *Signer) Logout() error {
	// Save the current state of the Signer.
	err := signer.Save()
	if err != nil {
		// Return any error encountered during saving.
		return err
	}
	// Mark the Signer as not signed in.
	signer.signed = false
	return nil
}

// Deleted returns whether the Signer has been marked as deleted.
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - bool - true if the Signer is marked as deleted; otherwise, false.
func (signer *Signer) Deleted() bool {
	return signer.deleted
}

// SetUsername sets the username for the Signer.
//
// Parameters:
//
//   - username - a string representing the new username to assign to the Signer.
//
// Returns:
//
//   - none
func (signer *Signer) SetUsername(username string) {
	signer.username = username
}

// SetEmail sets the email address for the Signer.
//
// Parameters:
//
//   - email - a string representing the new email address to assign to the Signer.
//
// Returns:
//
//   - none
func (signer *Signer) SetEmail(email string) {
	signer.email = email
}

// Delete removes the Signer's record from the log file and marks the Signer as deleted.
//
// This method refreshes the log file, scans through each line, and excludes the entry
// matching the Signer's username and email. The updated list is then saved back to the file.
// After successful deletion, the Signer's cart and wallet are set to nil, and the deleted and
// signed flags are updated accordingly.
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - error - an error if writing the updated log file fails; otherwise, nil.
func (signer *Signer) Delete() error {
	// Refresh the log file to ensure it's up to date.
	Refresh()
	file := LogFile

	// Initialize a scanner to read the log file line by line.
	scanner := bufio.NewScanner(file)
	defer func() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	lines := make([]string, 0)

	// Iterate through each line in the log file.
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		username, email, _, _ := parts[0], parts[1], parts[2], parts[3]

		// Skip the line if it matches the Signer's username and email.
		if username == signer.username && email == signer.email {
			continue
		}
		// Otherwise, keep the line.
		lines = append(lines, line)
	}

	// Write the updated lines back to the log file.
	err := writeSave(lines, file)
	if err != nil {
		return err
	}

	// Clear the Signer's cart and wallet, and update status flags.
	signer.cart = nil
	signer.wallet = nil
	signer.deleted = true
	signer.signed = true

	return nil
}

// saveSecret updates or adds the Signer's secret hash in the log file.
//
// This method scans the log file for an entry matching the Signer's username and email.
// If found, it updates the entry with the new hash and current wallet balance.
// If not found, it appends a new entry for the Signer. The updated list is then saved
// back to the log file.
//
// Parameters:
//
//   - hash - the secret hash to be saved for the Signer.
//
// Returns:
//
//   - error - an error if writing the updated log file fails; otherwise, nil.
func (signer *Signer) saveSecret(hash secret) error {
	// Refresh the log file to ensure it's up to date.
	Refresh()
	file := LogFile

	// Initialize a scanner to read the log file line by line.
	scanner := bufio.NewScanner(file)
	defer func() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	lines := make([]string, 0)
	exists := false

	// Iterate through each line in the log file.
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		username, email, _, _ := parts[0], parts[1], parts[2], parts[3]

		// If the line matches the Signer's username and email, update the entry.
		if username == signer.username && email == signer.email {
			line = fmt.Sprintf(saveFMT, signer.username, signer.email, hash, signer.wallet.Balance().InexactFloat64())
			exists = true
		}
		// Add the (possibly updated) line to the list.
		lines = append(lines, line)
	}

	// If no existing entry was found, append a new entry for the Signer.
	if !exists {
		newLine := fmt.Sprintf(saveFMT, signer.username, signer.email, hash, signer.Wallet().Balance().InexactFloat64())
		lines = append(lines, newLine)
	}

	// Write the updated lines back to the log file.
	err := writeSave(lines, file)
	if err != nil {
		return err
	}
	return nil
}

// Save updates the Signer's record in the log file with the current wallet balance.
//
// This method scans the log file for entries matching the Signer's username or email.
// If a match is found, it updates the corresponding line with the latest username, email,
// hash, and wallet balance. The updated list of lines is then written back to the log file.
//
// Parameters:
//   - none
//
// Returns:
//   - error - an error if writing the updated log file fails; otherwise, nil.
func (signer *Signer) Save() error {
	// Refresh the log file to ensure it's up to date.
	Refresh()
	file := LogFile

	// Initialize a scanner to read the log file line by line.
	scanner := bufio.NewScanner(file)
	defer func() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	lines := make([]string, 0)

	// Iterate through each line in the log file.
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) == 0 {
			continue
		}
		username, email, hash, _ := parts[0], parts[1], parts[2], parts[3]

		// If the line matches the Signer's username or email, update the entry with the current balance.
		if username == signer.username || email == signer.email {
			line = fmt.Sprintf(saveFMT, signer.username, signer.email, hash, signer.wallet.Balance().InexactFloat64())
		}
		// Add the (possibly updated) line to the list.
		lines = append(lines, line)
	}

	// Write the updated lines back to the log file.
	err := writeSave(lines, file)
	if err != nil {
		return err
	}
	return nil
}

// newSigner creates and returns a new Signer instance with the provided username and email.
//
// The returned Signer is initialized with a new Wallet, a new Cart, and is marked as signed.
//
// Parameters:
//
//	username - a string representing the username for the new Signer.
//	email    - a string representing the email address for the new Signer.
//
// Returns:
//
//	*Signer - a pointer to the newly created Signer instance.
func newSigner(username, email string) *Signer {
	return &Signer{
		username: username,
		email:    email,
		wallet:   wallets.NewWallet(),
		cart:     store.NewCart(),
		signed:   true,
	}
}

// InvalidUserErr is an error that occurs when a user is not signed in
var InvalidUserErr = fmt.Errorf("invalid user. User is not signed in")

// SignerFromUser creates a new Signer instance from a given User.
//
// This function checks if the provided User is verified and not deleted.
// If the User is invalid (either not verified or deleted), it returns nil and an InvalidUserErr.
// Otherwise, it constructs and returns a new Signer populated with the User's username, email,
// wallet, and cart, and marks the Signer as signed.
//
// Parameters:
//
//	user - a User instance from which to create the Signer.
//
// Returns:
//
//	*Signer - a pointer to the newly created Signer if the User is valid; otherwise, nil.
//	error   - InvalidUserErr if the User is not verified or is deleted; otherwise, nil.
func SignerFromUser(user User) (*Signer, error) {
	// Check if the user is verified and not deleted.
	if !user.Verified() || user.Deleted() {
		return nil, InvalidUserErr
	}
	// Create and return a new Signer based on the User's information.
	return &Signer{
		username: user.Username(),
		email:    user.Email(),
		wallet:   user.Wallet(),
		cart:     user.Cart(),
		signed:   true,
	}, nil
}
