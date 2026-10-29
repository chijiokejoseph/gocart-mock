package purchase

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/store"
)

func manageCart(reader *bufio.Reader, signer *auth.Signer, shop *store.Store, sp *store.MaxPrint) error {
	menu := []string{
		"Add items to your cart",
		"View Cart",
		"Remove items from your cart",
		"Clear cart",
		"Exit Manage Cart Menu",
	}
	prompt := fmt.Sprintf("Select an option [1 - %d]: ", len(menu))
	for {
		choice, err := inputs.InputOption(reader, prompt, menu, inputs.Config("Manage Cart Menu"))
		if err != nil {
			return err
		}
		switch choice {
		case 1:
			err := search(reader, signer, shop, sp)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
		case 2:
			view(signer, sp)
		case 3:
			err := remove(reader, signer, sp)
			if err != nil {
				return err
			}
		case 4:
			view(signer, sp)
			choice, err := inputs.InputBool(reader, "Do you wish to clear all the items in your cart?: ", "Clear Cart")
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
			if choice {
				signer.Cart().Clear()
			}
		default:
			return nil
		}
	}
}
