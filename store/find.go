package store

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	e "github.com/chijiokejoseph/gocart-mock/err"
)

var ErrQueryInvalid = errors.New("invalid search query")

var ErrNoItemFound = errors.New("no item found")

func matchReg(regArray []*regexp.Regexp, itemName string) bool {
	for _, reg := range regArray {
		if !reg.MatchString(itemName) {
			return false
		}
	}
	return true
}

func Find(query string, shop *Store) ([]*Product, error) {
	products := make([]*Product, 0)
	if query == "" {
		return products, ErrQueryInvalid
	}
	words := strings.Fields(query)
	regex := make([]*regexp.Regexp, len(words))
	for i, word := range words {
		reg, err := regexp.Compile(fmt.Sprintf(`(?i)%s`, word))
		if err != nil {
			fmt.Printf("\n❌ Query entered by the user has invalid characters. Please avoid the use of symbols or punctuations in your query and try again.\n")
			return products, fmt.Errorf("%w", err)
		}
		regex[i] = reg
	}
	for _, name := range shop.Names() {
		if !matchReg(regex, name) {
			continue
		}
		item := e.Must(shop.Get(name))
		products = append(products, item)
	}
	if len(products) == 0 {
		fmt.Printf("\n⚠️ No item found for user's query '%s'\n", query)
		return products, ErrNoItemFound
	}
	return products, nil
}
