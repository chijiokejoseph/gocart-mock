package auth

import (
	"github.com/chijiokejoseph/gocart-mock/store"
	"github.com/chijiokejoseph/gocart-mock/wallets"
)


// User is an interface that represents a user of the program
// and provides access to their information
type User interface {
	Username() string
	Email() string
	Wallet() *wallets.Wallet
	Cart() *store.Cart
	Verified() bool
	Signed() bool
	Deleted() bool
}
