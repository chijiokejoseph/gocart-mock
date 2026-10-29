package auth

import (
	"bufio"
	"errors"
	"fmt"
	"log"

	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/wallets"
)

// NewSigner creates a new Signer instance, prompting the user to either manually enter
// a new password or generate one automatically. The function interacts with the user
// via the provided bufio.Reader, and uses the given username and email for the Signer.
// If a Signer instance is not provided, a new one is created. The password is hashed
// and securely saved using the Signer. The function handles user input errors and
// system write errors gracefully, returning the created Signer or an error.
//
// Parameters:
//
//   - reader   - *bufio.Reader: Reader for user input.
//   - username - string: Username for the new Signer.
//   - email    - string: Email for the new Signer.
//   - signer   - *Signer: Optional existing Signer instance.
//
// Returns:
//
//   - *Signer: The created or updated Signer instance.
//   - error:   Error encountered during password entry, generation, or saving.
func NewSigner(reader *bufio.Reader, username, email string, signer *Signer) (*Signer, error) {
	menu := []string{"Enter New Password (manual)", "Generate New Password (automatic)"}
	var (
		password, hash secret
	)
	for {
		// Prompt the user to choose how to generate the new password
		choice, err := inputs.InputOption(
			reader,
			"How do you wish to generate your new password?: ",
			menu,
			inputs.NewInputConfig().WithTitle("Password Generation Menu"),
		)
		if err != nil {
			switch {
			// If the error is a system input error, return nil
			case errors.Is(err, inputs.ErrSysInput):
				return nil, err
			// For other errors, prompt again
			default:
				continue
			}
		}
		switch choice {
		case 1:
			// User chose to manually enter a password
			password, hash, err = enterPassword(reader)
			if err != nil {
				switch {
				// If the error is a system input error, return the error
				case errors.Is(err, inputs.ErrSysInput):
					return nil, err
				// For other errors, prompt again
				default:
					continue
				}
			}
		default:
			// User chose to automatically generate a password
			password, hash = generatePassword()
		}
		// Exit the loop after successful password entry/generation
		break
	}

	// Display the new password to the user and remind them to store it securely
	fmt.Printf("\nYour new password is: %s\nTake note of it and store it in a secure place.", password)

	// If no signer was provided, create a new one
	if signer == nil {
		signer = newSigner(username, email)
	}

	// Attempt to save the password hash using the signer
	err := signer.saveSecret(hash)
	if err != nil {
		switch true {
		// Handle system write errors
		case errors.Is(err, ErrSysWrite):
			fmt.Println("\nSystem write error. Failed to save password.")
			return nil, err
		// Handle maximum account errors
		case errors.Is(err, ErrMaxAccount):
			fmt.Println("\nMaximum number of accounts reached")
			return nil, err
		// Handle any unknown errors by logging and terminating
		default:
			log.Fatalln("Unknown case reached")
		}
	}
	// Return the created or updated signer
	return signer, nil
}

// GetSigner retrieves and authenticates a Signer for the specified account.
// It looks up account details, prompts the user for their password, validates
// the password against the stored hash, and initializes the Signer with the
// account's balance.
//
// Parameters:
//
//   - reader  - *bufio.Reader: Reader for user input (used to prompt for password).
//   - account - string: The account identifier to look up.
//
// Returns:
//
//   - *Signer: The authenticated Signer instance for the account.
//   - error:   An error if lookup fails, password input is interrupted, or validation fails.
func GetSigner(reader *bufio.Reader, account string) (*Signer, error) {
	// Lookup retrieves username, email, password hash, and balance for the account.
	username, email, hashTxt, balance, err := Lookup(account, true)
	if err != nil {
		return nil, err
	}

	// Prompt the user to enter their password.
	password, err := inputs.InputStr(reader, "Enter your password: ")
	if err != nil {
		// Return immediately if input is interrupted or user quits.
		return nil, err
	}

	// Validate the entered password against the stored hash.
	err = validatePasswordByHash(password, hashTxt)
	if err != nil {
		return nil, err
	}

	// Create a new Signer instance with the retrieved username and email.
	signer := newSigner(username, email)

	// Fund the signer's wallet with the account's balance.
	signer.Wallet().Fund(wallets.AirDrop(balance.InexactFloat64()))

	// Return the authenticated signer.
	return signer, nil
}
