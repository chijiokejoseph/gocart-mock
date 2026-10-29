package purchase

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/store"
)

func Mainloop(reader *bufio.Reader, signer *auth.Signer, shop *store.Store, sp *store.MaxPrint) error {
	menu := []string{
		"Search",
		"Manage Cart",
		"Checkout",
		"Exit Purchase Menu",
	}
	prompt := fmt.Sprintf("Select an option [1 - %d]: ", len(menu))
	for {
		choice, err := inputs.InputOption(reader, prompt, menu, inputs.Config("Purchase Menu"))
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
			err := manageCart(reader, signer, shop, sp)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
		case 3:
			err := checkout(reader, signer, sp)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
		default:
			return nil
		}
	}
}
