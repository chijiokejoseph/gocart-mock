package main

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/auth"
	e "github.com/chijiokejoseph/gocart-mock/err"
	"github.com/chijiokejoseph/gocart-mock/inputs"
	"github.com/chijiokejoseph/gocart-mock/routines"
	"github.com/chijiokejoseph/gocart-mock/routines/login"
	"github.com/chijiokejoseph/gocart-mock/store"
)

func mainloop(reader *bufio.Reader, shop *store.Store, sp *store.MaxPrint) {
	// initialize user as a guest by defaul
	var (
		user auth.User = auth.NewGuest()
	)

	menu := []string{
		"Login / Register",
		"Continue as Guest",
		"Exit",
	}

	for {

		choice, err := inputs.InputOption(reader, fmt.Sprintf("Select an operation (1 - %d): ", len(menu)), menu, inputs.NewInputConfig().WithTitle("Login Menu"))
		if err != nil {
			switch {
			// if user cancels login, return to the main function
			case errors.Is(err, inputs.ErrSysInput) || errors.Is(err, inputs.ErrQuit):
				return
			default:
				// on any other error simply re-attempt login
				continue
			}
		}

		switch choice {
		case 1:
			// attempt to login user
			user, err = login.Login(reader)

			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return
				case errors.Is(err, inputs.ErrQuit):
					continue
				default:
					continue
				}
			}

			// logout user automatically once function execution is complete
			defer func() {
				login.Logout(user)
				fmt.Println("\nℹ️ The program has exited.")
			}()
		case 2:
		case 3:
			return
		}

		// if user is a guest
		if !user.Signed() {
			// convert user to guest proper
			guest := e.Must(auth.GuestFromUser(user))

			// allow guests to search cart only
			err := routines.Preview(reader, guest, shop, sp)

			// return if there is a system level failure if not continue
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return
				default:
					continue
				}
			}
			continue
		}

		// convert user to signer
		signer := e.Must(auth.SignerFromUser(user))

		// carry out whatever operations user may wish to perform
		err = run(reader, signer, shop, sp)
		if err != nil {
			switch {
			// if there is a system level error return to main menu
			case errors.Is(err, inputs.ErrSysInput):
				return
			default:
				continue
			}
		}

		// if there's no error, the user has exited the program, hence break
		break
	}
}
