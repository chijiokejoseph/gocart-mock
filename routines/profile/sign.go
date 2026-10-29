package profile

import (
	"bufio"
	"github.com/chijiokejoseph/gocart-mock/auth"
	"fmt"
)

func authorize(reader *bufio.Reader, signer *auth.Signer) error {
	fmt.Println("\nTo proceed you'll need to verify your identity")
	_, err := auth.GetSigner(reader, signer.Username())
	if err != nil {
		return err
	}
	return nil
}
