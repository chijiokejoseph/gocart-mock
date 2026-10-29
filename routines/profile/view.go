package profile

import (
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/conv"
	e "github.com/chijiokejoseph/gocart-mock/err"
)

func view(signer *auth.Signer) {
	balance := signer.Wallet().Balance()
	fmt.Println("\nYour username is: ", signer.Username())
	fmt.Println("\nYour email is: ", signer.Email())
	fmt.Println("\nYour balance is NGN ", e.Must(conv.Dec2StrFixed(balance)))
}
