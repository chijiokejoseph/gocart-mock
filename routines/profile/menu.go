package profile

import (
	"bufio"
	"errors"
	"fmt"
	"slices"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/inputs"
)

func ManageAccount(reader *bufio.Reader, signer *auth.Signer) error {
	menu := []string{
		"Change your username",
		"Change your email",
		"Change your password",
		"View Account Balance/Details",
		"Reset Balance",
		"Delete Account",
		"Logout",
		"Exit Manage Account Menu",
	}
	prompt := fmt.Sprintf("Select an option [1 - %d]: ", len(menu))
	for {
		choice, err := inputs.InputOption(reader, prompt, menu, inputs.Config("Manage Account"))
		if err != nil {
			return err
		}
		unauthorizedCases := []int{len(menu) - 1, len(menu)}
		if slices.Contains(unauthorizedCases, choice) {
			switch choice {
			case len(menu) - 1:
				err := signer.Logout()
				if err != nil {
					return err
				}
				fmt.Printf("\nLogged user '%s' out successfully\n", signer.Username())
			case len(menu):
			}
			return nil
		}
		err = authorize(reader, signer)
		if err != nil {
			switch {
			case errors.Is(err, inputs.ErrSysInput):
				return err
			default:
				continue
			}
		}
		switch choice {
		case 1:
			err := setUsername(reader, signer)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
		case 2:
			err := setEmail(reader, signer)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
		case 3:
			err := setPassword(reader, signer)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
		case 4:
			view(signer)
		case 5:
			err := clearBalance(reader, signer)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
		case 6:
			err := deleteAccount(reader, signer)
			if err != nil {
				switch {
				case errors.Is(err, inputs.ErrSysInput):
					return err
				default:
					continue
				}
			}
			return nil
		}
	}
}
