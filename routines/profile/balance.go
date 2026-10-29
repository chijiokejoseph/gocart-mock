package profile

import (
	"bufio"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/conv"
	e "github.com/chijiokejoseph/gocart-mock/err"
	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/wallets"
)

func clearBalance(reader *bufio.Reader, signer *auth.Signer) error {
	fmt.Println("\nThis action will withdraw all money from your app account. Please confirm.")
	choice, err := inputs.InputBool(reader, "Yes / No: ", "Confirm")
	if err != nil {
		return err
	}
	if !choice {
		return nil
	}
	oldBalance := signer.Wallet().Balance()
	signer.Wallet().Clear()
	if err := signer.Save(); err != nil {
		signer.Wallet().Fund(wallets.AirDrop(oldBalance.InexactFloat64()))
	}
	fmt.Printf("\nYour wallet has been cleared successfully. Your new balance is %s\n", e.Must(conv.Dec2StrFixed(signer.Wallet().Balance())))
	return nil
}
