package purchase

import (
	"bufio"
	"github.com/chijiokejoseph/gocart-mock/auth"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/store"
)

func remove(reader *bufio.Reader, signer *auth.Signer, sp *store.MaxPrint) error {
	cart := signer.Cart()
	for {
		orders := cart.Orders()
		items := store.FromOrders(orders)
		store.Showcase(items, sp)
		fmt.Println()
		prompt := fmt.Sprintf("Enter the Display No of the item you wish to remove from your cart [1 - %d]: ", len(items))
		itemNo, err := inputs.InputNum(reader, prompt)
		if err != nil {
			return err
		}
		cart.RemoveOrder(orders[itemNo-1])
		fmt.Println()
		choice, err := inputs.InputBool(reader, "Do you wish to remove more items from your cart?: ", "Remove More Items")
		if err != nil {
			return err
		}
		if choice {
			continue
		} else {
			break
		}
	}
	return nil
}
