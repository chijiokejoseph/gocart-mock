package routines

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/store"
)

func Preview(reader *bufio.Reader, guest *auth.Guest, shop *store.Store, sp *store.MaxPrint) error {
	for {
		// Prompt the user to enter search query for a shoppint item
		prompt := "\nWelcome to the store. Please enter the name or brand of item you are seeking: "
		query, err := inputs.InputStr(reader, prompt)
		if err != nil {
			return err
		}

		// find matching products based on user search query
		products, err := store.Find(query, shop)
		switch true {
		// If error occurs with user query, reprompt the user
		case errors.Is(err, store.ErrQueryInvalid):
			continue
		case errors.Is(err, store.ErrNoItemFound):
			continue
		}

		items := store.FromProds(products)
		fmt.Println("\nThe following items matched your query")
		store.Showcase(items, sp)
		again, err := inputs.InputBool(reader, "Do you wish to search again?: ", "Search Store Again")
		if err != nil {
			switch {
			case errors.Is(err, inputs.ErrSysInput) || errors.Is(err, inputs.ErrQuit):
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
