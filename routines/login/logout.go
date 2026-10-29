package login

import (
	"github.com/chijiokejoseph/gocart-mock/auth"
	"fmt"
	"time"
)

func Logout(user auth.User) {
	signer, err := auth.SignerFromUser(user)
	if err != nil {
		time.Sleep(700 * time.Millisecond)
		fmt.Println("\nLogging guest user out...")
	} else {
		var err error
		for i := 0; i < 5; i++ {
			err = signer.Save()
			if err != nil {
				fmt.Printf("\nFailed to log user out. Retrying... (%d/5) tries\n", i+1)
				time.Sleep(700 * time.Millisecond)
				continue
			}
			break
		}
		if err != nil {
			fmt.Println("Couldn't log user out")
		} else {
			fmt.Printf("\nLogged user '%s' out successfully.", signer.Username())
		}
		return
	}
}
