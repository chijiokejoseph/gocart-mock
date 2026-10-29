package purchase

import (
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/store"
)

func view(signer *auth.Signer, sp *store.MaxPrint) {
	orders := signer.Cart().Orders()
	if len(orders) == 0 {
		fmt.Println("\nThere are no items in your cart")
		return
	}
	fmt.Println("\nThe items in your cart are as follows:")
	items := store.FromOrders(orders)
	store.Showcase(items, sp)
	fmt.Println()
}
