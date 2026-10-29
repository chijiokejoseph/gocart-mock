package profile

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
)

var passwordErr = errors.New("failed to set password")

func setPassword(reader *bufio.Reader, signer *auth.Signer) error {
	_, err := auth.NewSigner(reader, signer.Username(), signer.Email(), signer)
	if err != nil {
		switch {
		case errors.Is(err, inputs.ErrSysInput):
			return err
		default:
			return fmt.Errorf("%w: %s", passwordErr, err)
		}
	}
	fmt.Println("\nPassword change successful")
	return nil
}
