package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/chijiokejoseph/gocart-mock/auth"
	"github.com/chijiokejoseph/gocart-mock/paths"
	"github.com/chijiokejoseph/gocart-mock/store"
)

func main() {
	// TODO 1: test run the entire project to test for normal feature behavior and 
	// spot potential discrepancies
	defer func() {
		if err := auth.LogFile.Close(); err != nil {
			log.Fatalln(err)
		}
	}()


	reader := bufio.NewReader(os.Stdin)
	shop := store.BuildStore(paths.DataPath())
	sp := shop.MaxPrint()

	// TODO 2: imporove console UI with emojis for a better UX
	fmt.Println("Welcome to the E-commerce-app")
	mainloop(reader, shop, sp)
}
