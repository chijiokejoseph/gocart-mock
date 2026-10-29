package purchase

import (
	"bufio"
	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/store"
	"fmt"
)

func add(reader *bufio.Reader, signer *auth.Signer, products []*store.Product, sp *store.MaxPrint) error {
	for {
		fmt.Println()
		items := store.FromProds(products)
		store.Showcase(items, sp)
		fmt.Println()
		item, err := inputs.InputNum(reader, "Enter the Search No of the item that you wish to add to cart: ")
		if err != nil {
			return err
		}
		prod := products[item-1]
		signer.Cart().AddOrder(prod)
		fmt.Printf("\n%s\n", prod.String())
		choice, err := inputs.InputBool(reader, "Do you wish to continue adding items?: ", "Add More Items")
		if err != nil {
			return err
		}
		if choice {
			continue
		}
		return nil
	}
}
