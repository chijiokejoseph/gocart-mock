package profile

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
)

var usernameErr = errors.New("username set failed")

func setUsername(reader *bufio.Reader, signer *auth.Signer) error {
	oldUsername := signer.Username()
	fmt.Println("\nThis is your existing username: ", oldUsername)
	username, err := inputs.InputStr(reader, "Enter your new username: ")
	if err != nil {
		switch {
		case errors.Is(err, inputs.ErrSysInput):
			return err
		default:
			return fmt.Errorf("%w. %w", usernameErr, err)
		}
	}
	signer.SetUsername(username)
	err = signer.Save()
	if err != nil {
		signer.SetUsername(oldUsername)
		return fmt.Errorf("%w. %w", usernameErr, err)
	}
	fmt.Println("\nUsername change successful")
	return nil
}
