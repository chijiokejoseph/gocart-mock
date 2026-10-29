package routines

import (
	"bufio"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/conv"
	e "github.com/chijiokejoseph/gocart-mock/err"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/wallets"
)

func Fund(reader *bufio.Reader, signer *auth.Signer) error {
	wallet := signer.Wallet()
	menu := []string{"Yes", "No"}
	airdropSet := []wallets.AirDrop{wallets.ONE, wallets.TWO, wallets.FIVE, wallets.TEN}
	airdropOptions := []string{
		wallets.ONE.String(),
		wallets.TWO.String(),
		wallets.FIVE.String(),
		wallets.TEN.String(),
	}
	consentPrompt := fmt.Sprintf("Do you wish to fund your account again? Enter [1 - %d]: ", len(menu))
	amountPrompt := fmt.Sprintf("By how much do you wish to fund account? Enter [1 - %d]: ", len(airdropOptions))
	for {
		choice, err := inputs.InputOption(reader, amountPrompt, airdropOptions, inputs.NewInputConfig().WithTitle("Airdrop Menu"))
		if err != nil {
			// whether the error is an inputs.ErrSysInput or an inputs.ErrQuit return it
			return err
		}
		airdrop := airdropSet[choice-1]
		wallet.Fund(airdrop)
		airdropDec := conv.Float2Dec(float64(airdrop))
		airdropStr := e.Must(conv.Dec2Str(airdropDec, conv.NewOptions().WithPrecision(0)))
		fmt.Printf("\nYou have added NGN %s to your account\n", airdropStr)

		choice, err = inputs.InputOption(reader, consentPrompt, menu, inputs.NewInputConfig().WithTitle("Fund Wallet"))
		if err != nil {
			// whether the error is an inputs.ErrSysInput or inputs.ErrQuit return it
			return err
		}
		switch choice {
		// user enters 'yes'
		case 1:
		// user enters 'no'
		case 2:
			err := signer.Save()
			if err != nil {
				fmt.Println("\nCouldn't update user's account balance")
				return err
			}
			fmt.Printf("\nYour wallet now has the following balance: %s\n", e.Must(conv.Dec2StrFixed(wallet.Balance())))
			return nil
		// user enters invalid option number
		default:
			continue
		}
	}
}
