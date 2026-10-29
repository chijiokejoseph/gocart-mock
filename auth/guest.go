package auth

import (
	"errors"

	"github.com/chijiokejoseph/gocart-mock/store"
	"github.com/chijiokejoseph/gocart-mock/wallets"
)

var ErrUserStillSigned = errors.New("cannot convert a Signed User to Guest except the User is signed out.")

// Guest represents a guest user
// that is not signed in
//
// Fields:
//   - username: The username of the guest user
//   - email: The email of the guest user
//   - wallet: The wallet of the guest user
//   - cart: The cart of the guest user
//   - signed: Whether the guest user is signed in
//   - deleted: Whether the guest user is deleted
type Guest struct {
	username string
	email    string
	wallet   *wallets.Wallet
	cart     *store.Cart
	signed   bool
	deleted  bool
}

// Username returns the username of the guest user
//
// Parameters: None
//
// Returns:
//   - string: The username of the guest user
func (guest *Guest) Username() string {
	return guest.username
}

// Email returns the email of the guest user.
//
// Parameters: None
//
// Returns:
//   - string: The email of the guest user

func (guest *Guest) Email() string {
	return guest.email
}

// Wallet returns the wallet of the guest user.
//
// Parameters: None
//
// Returns:
//   - *wallets.Wallet: Always returns nil as guest users do not have a wallet

func (guest *Guest) Wallet() *wallets.Wallet {
	return nil
}

// Cart returns the cart of the guest user.
//
// Parameters: None
//
// Returns:
//   - *store.Cart: Always returns nil as guest users do not have a cart
func (guest *Guest) Cart() *store.Cart {
	return nil
}

// Verified returns whether the guest user has been verified.
//
// Parameters: None
//
// Returns:
//   - bool: Always returns false as guest users do not have a verification status
func (guest *Guest) Verified() bool {
	return false
}

// Signed returns whether the guest user is signed in.
//
// Parameters: None
//
// Returns:
//   - bool: Always returns false as guest users are never signed in
func (guest *Guest) Signed() bool {
	return false
}

// Deleted returns whether the guest user is deleted.
//
// Parameters: None
//
// Returns:
//   - bool: Always returns false as guest users are never deleted
func (guest *Guest) Deleted() bool {
	return false
}

// NewGuest creates a new Guest user.
//
// Parameters: None
//
// Returns:
//   - *Guest: A pointer to a newly created Guest instance with default values.
func NewGuest() *Guest {
	return &Guest{
		username: "Guest",
		email:    "Guest",
		wallet:   nil,
		cart:     nil,
	}
}

// ErrSignerPurgeFailed is an error that occurs when trying to convert a signer to a guest
var ErrSignerPurgeFailed = errors.New("cannot convert undeleted signer to guest: signer data deletion failed")

// GuestFromSigner takes a signer and returns a Guest.
//
// If the signer is undeleted, it returns an error as the Guest type is
// only used for deleted signers. The error is `ErrSignerPurgeFailed`.
//
// Parameters:
//   - signer *Signer: The signer to convert.
//
// Returns:
//   - *Guest: A pointer to a Guest with the same username and email as the signer.
//   - error: An error if the signer is undeleted, otherwise nil.
func GuestFromSigner(signer *Signer) (*Guest, error) {
	if signer.deleted {
		return nil, ErrSignerPurgeFailed
	}
	return &Guest{
		username: signer.username,
		email:    signer.email,
		wallet:   nil,
		cart:     nil,
	}, nil
}

func GuestFromUser(user User) (*Guest, error) {
	if user.Signed() {
		return nil, ErrUserStillSigned
	}

	return &Guest{
		username: user.Username(),
		email:    user.Email(),
	}, nil
}
