package profile

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
)

var emailErr = errors.New("set email error")

func setEmail(reader *bufio.Reader, signer *auth.Signer) error {
	oldEmail := signer.Email()
	fmt.Println("\nThis is your existing email: ", oldEmail)
	email, err := inputs.InputStr(reader, "Enter your new email: ")
	if err != nil {
		switch {
		case errors.Is(err, inputs.ErrSysInput):
			return err
		default:
			return fmt.Errorf("%w. %w", emailErr, err)
		}
	}
	signer.SetEmail(email)
	if err := signer.Save(); err != nil {
		signer.SetEmail(oldEmail)
		return fmt.Errorf("%w. %w", emailErr, err)
	}
	fmt.Println("\nEmail change successful")
	return nil
}
