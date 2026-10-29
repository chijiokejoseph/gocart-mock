package login

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
)

// signUp is a routine function that executes an interactive
// process to get a user's details for an account sign-up.
//
// if the function experiences an error in reading user input
// `inputs.ErrSysInput` or the user cancels the operation
// `inputs.ErrQuit`, the function returns
func signUp(reader *bufio.Reader) (*auth.Signer, error) {
	fmt.Println("\n🔊 Welcome to the E-commerce Sign Up Page")
	fmt.Println("\n🔊 Please enter your details to create an account.")
	var (
		username string
		email    string
	)
	for {
		value, err := inputs.InputStr(reader, "ℹ️ Enter your username: ")
		username = value
		if err != nil {
			// whether the error is an inputs.ErrSysInput or an inputs.ErrQuit return it
			return nil, err
		}
		matchUsername, _, _, _, lookupErr := auth.Lookup(username, false)
		if lookupErr == nil {
			if matchUsername == username {
				fmt.Printf("\n⚠️ Username '%s' has already been taken.\n", username)
				continue
			}
		}
		email, err = inputs.InputStr(reader, "ℹ️ Enter your email: ")
		if err != nil {
			// whether the error is an inputs.ErrSysInput or an inputs.ErrQuit return it
			return nil, err
		}
		_, matchEmail, _, _, lookupErr := auth.Lookup(email, false)
		if lookupErr == nil {
			if matchEmail == email {
				fmt.Printf("\n⚠️ Email '%s' has already been used.\n", email)
				continue
			}
		}
		break
	}
	return auth.NewSigner(reader, username, email, nil)
}

// signIn is a routine function that executes an interactive program
// to get the details required to sign in from a user.

// if the function experiences an error reading in user input
// `inputs.ErrSysInput` or if the user cancels the operation
// `inputs.ErrQuit`, the function returns.
func signIn(reader *bufio.Reader) (*auth.Signer, error) {
	fmt.Println("\n🔊 Welcome to the E-commerce Sign In Page")
	for {
		fmt.Println("\n🔊 Please enter your details to sign in.")
		account, err := inputs.InputStr(reader, "ℹ️ Enter your username/email: ")
		if err != nil {
			return nil, err
		}
		signer, err := auth.GetSigner(reader, account)
		if err != nil {
			switch {
			case errors.Is(err, inputs.ErrSysInput) || errors.Is(err, inputs.ErrQuit):
				return nil, err
			default:
				continue
			}
		}
		return signer, nil
	}
}

func Login(reader *bufio.Reader) (*auth.Signer, error) {
	menu := []string{"Sign In", "Sign Up", "Return"}
	prompt := fmt.Sprintf("ℹ️ Do you wish to sign in or sign up? Enter [1 - %d]: ", len(menu))
	for {
		// prompts the user to sign in, sign up or return to main menu
		choice, err := inputs.InputOption(reader, prompt, menu, inputs.NewInputConfig().WithTitle("Sign In / Sign Up"))
		if err != nil {
			return nil, err
		}
		switch choice {
		case 1:
			signer, err := signIn(reader)
			// `signIn` only returns user cancelation error
			if err != nil {
				return nil, err
			}
			return signer, nil
		case 2:
			signer, err := signUp(reader)
			// `signUp` only returns user cancelation error
			if err != nil {
				return nil, err
			}
			return signer, nil
		}
	}
}
