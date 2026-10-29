package purchase

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/conv"
	e "github.com/chijiokejoseph/gocart-mock/err"

	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/store"
)

func confirm(reader *bufio.Reader, signer *auth.Signer, sp *store.MaxPrint) (bool, error) {
	fmt.Println("\nBefore payment please confirm the items in your cart")
	fmt.Println("\nCart:")
	items := store.FromOrders(signer.Cart().Orders())
	store.Showcase(items, sp)

	for {
		choice, err := inputs.InputBool(reader, "Do you wish to continue to payment?: ", "Confirm Cart")
		if err != nil {
			if errors.Is(err, inputs.ErrSysInput) || errors.Is(err, inputs.ErrQuit) {
				return false, nil
			}
			continue
		}
		return choice, nil
	}
}

func pay(reader *bufio.Reader, signer *auth.Signer) error {
	c := signer.Cart()
	w := signer.Wallet()
	fmt.Printf("\nYou are about to be charged the following fee for the items in your cart: NGN %s\n", e.Must(conv.Dec2Str(c.Total(), conv.NewOptions().WithPrecision(2))))
	choice, err := inputs.InputBool(reader, "Yes/No: ", "Pay")
	if err != nil {
		return err
	}
	if !choice {
		return nil
	}
	err = w.Charge(c.Total())
	if err != nil {
		return err
	}
	c.ClearPostPay()
	fmt.Println("\nPayment successful. Your cart has been cleared and the items have been delivered")
	return nil
}

func checkout(reader *bufio.Reader, signer *auth.Signer, sp *store.MaxPrint) error {
	confirmed, err := confirm(reader, signer, sp)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("\nYou have cancelled this payment")
		return nil
	}
	return pay(reader, signer)
}
