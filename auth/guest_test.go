package auth_test

import (
	"testing"

	"github.com/chijiokejoseph/gocart-mock/auth"
)

// TestNewGuest tests the initialization of a new guest user using the auth.NewGuest() function.
// It verifies that all fields are set to their expected zero values.
// 
// Params:
//   - t: *testing.T - the testing context provided by the Go testing framework.
func TestNewGuest(t *testing.T) {
	// Create a new guest user
	guest := auth.NewGuest()

	// Check that the guest's cart is nil
	if guest.Cart() != nil {
		t.Errorf("NewGuest() does not have a cart of nil pointer")
	}

	// Check that the guest's wallet is nil
	if guest.Wallet() != nil {
		t.Errorf("NewGuest() does not have a wallet of nil pointer")
	}

	// Check that the guest's username is an empty string
	if guest.Username() != "" {
		t.Errorf("NewGuest() should have a Username of empty string")
	}

	// Check that the guest's email is an empty string
	if guest.Email() != "" {
		t.Errorf("NewGuest() should have an Email of empty string")
	}

	// Check that the guest's deleted flag is false
	if guest.Deleted() != false {
		t.Errorf("NewGuest() should have a Deleted value of false")
	}

	// Check that the guest's verified flag is false
	if guest.Verified() != false {
		t.Errorf("NewGuest() should have a Verified value of false")
	}
}
