package purchase

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/store"
)

func research(reader *bufio.Reader, signer *auth.Signer, products []*store.Product, sp *store.MaxPrint) (bool, error) {
	menu := []string{
		"Search Again",
		"Add to cart",
		"Exit Search Menu",
	}
	confirmationPrompt := fmt.Sprintf("What do you wish to do next? Enter an option [1 - %d]: ", len(menu))
	for {
		choice, err := inputs.InputOption(reader, confirmationPrompt, menu, inputs.Config("Search Menu "))
		if err != nil {
			return false, err
		}
		switch choice {
		case 1:
			return true, nil
		case 2:
			err := add(reader, signer, products, sp)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return false, err
				default:
					continue
				}
			}
		default:
			return false, nil
		}
	}
}

func search(reader *bufio.Reader, signer *auth.Signer, shop *store.Store, sp *store.MaxPrint) error {
	for {
		prompt := "\nWelcome to the store. Please enter the name or brand of item you are seeking: "
		query, err := inputs.InputStr(reader, prompt)
		if err != nil {
			return err
		}
		products, err := store.Find(query, shop)
		switch true {
		case errors.Is(err, store.ErrQueryInvalid):
			continue
		case errors.Is(err, store.ErrNoItemFound):
			continue
		}

		items := store.FromProds(products)
		fmt.Println("\nThe following items matched your query")
		store.Showcase(items, sp)
		again, err := research(reader, signer, products, sp)
		if err != nil {
			switch {
			case errors.Is(err, inputs.ErrSysInput):
				return err
			default:
				continue
			}
		}
		if again {
			continue
		}
		return nil
	}
}
