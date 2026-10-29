package auth

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/inputs"
)

// enterPassword prompts the user to enter a new password, validates it, and returns the password and its hash.
// The function uses the provided bufio.Reader to read input from the user. It will repeatedly prompt the user
// until a valid password is entered or a system/quit error occurs.
//
// Parameters:
//   - reader *bufio.Reader - The input reader used to capture user input from the console.
//
// Returns:
//   - password secret - The entered password as a secret type.
//   - hash     secret - The hashed version of the entered password as a secret type.
//   - err      error  - An error if input fails or the user chooses to quit; otherwise, nil.
func enterPassword(reader *bufio.Reader) (password secret, hash secret, err error) {
	for {
		// Prompt the user to enter a new password.
		passwordTxt, err := inputs.InputStr(reader, "Enter your new password: ")
		if err != nil {
			// If the error is a system input error or the user chose to quit, return immediately.
			switch true {
			case errors.Is(err, inputs.ErrSysInput) || errors.Is(err, inputs.ErrQuit):
				return "", "", err
			}
		}

		// Attempt to parse and hash the entered password.
		password, hash, err := fromPasswordString(passwordTxt)
		if err != nil {
			// If password validation or hashing fails, inform the user and prompt again.
			fmt.Printf("\n%s\n", err)
			fmt.Println("\nPlease try again.")
			continue
		}

		// If everything is successful, return the password and its hash.
		return password, hash, nil
	}
}
