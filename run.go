package main

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/routines"
	"github.com/chijiokejoseph/gocart-mock/routines/profile"
	"github.com/chijiokejoseph/gocart-mock/routines/purchase"
	"github.com/chijiokejoseph/gocart-mock/store"
)

var ErrLoggedOut = errors.New("user has been logged out")

func run(reader *bufio.Reader, signer *auth.Signer, shop *store.Store, sp *store.MaxPrint) error {

	mainMenu := []string{
		"Fund your Wallet",
		"Purchase",
		"Manage Account",
		"Exit",
	}
	for {
		// display `run` main menu and receive user input
		choice, err := inputs.InputOption(reader, fmt.Sprintf("Select an option [1 - %d]: ", len(mainMenu)), mainMenu, inputs.NewInputConfig().WithTitle("Main Menu"))

		// handle user input errors
		if err != nil {
			switch {
			// if error is systemic or error due to operation cancelation from the user, return
			case errors.Is(err, inputs.ErrSysInput) || errors.Is(err, inputs.ErrQuit):
				return err
			// if error is user induced, but not user quit, reprompt the user
			default:
				continue
			}
		}

		if choice == len(mainMenu) {
			// if user selects to exit
			fmt.Println("\n🎉🎉 Thanks for shopping with us today. Hope you got all the good stuffs. Check us again next time. 🎊🎊")
			return nil
		}

		switch choice {
		case 1:
			err = routines.Fund(reader, signer)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
		case 2:
			err := purchase.Mainloop(reader, signer, shop, sp)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
		case 3:
			err = profile.ManageAccount(reader, signer)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
			if signer.Deleted() || !signer.Signed() {
				return ErrLoggedOut
			}
		}
	}
}
