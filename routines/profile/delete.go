package profile

import (
	"bufio"
	"errors"
	"fmt"
	"log"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
)

var deleteErr = errors.New("failed to delete profile")

func deleteAccount(reader *bufio.Reader, signer *auth.Signer) error {
	choice, err := inputs.InputBool(reader, "This action is irreversible. Please confirm that you wish to delete this account?: ", "Delete Account")
	username := signer.Username()
	if err != nil {
		switch {
		case errors.Is(err, inputs.ErrSysInput):
			return err
		default:
			return fmt.Errorf("%w. %w", deleteErr, err)
		}
	}
	if !choice {
		return fmt.Errorf("%w. %w", deleteErr, err)
	}
	err = signer.Delete()
	if err != nil {
		fmt.Println("\nFailed to delete account")
		log.Println(err)
		return fmt.Errorf("%w. %w", deleteErr, err)
	}
	fmt.Printf("\nUser '%s' acount deleted successfully\n", username)
	return nil
}
