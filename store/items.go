package store

import (
	"fmt"
	"slices"
	"strings"

	"github.com/chijiokejoseph/gocart-mock/conv"
	e "github.com/chijiokejoseph/gocart-mock/err"
	"github.com/shopspring/decimal"
)

type Item interface {
	Name() string
	Price() decimal.Decimal
	Stock() int
	String() string
	AddItem()
	RemoveItem()
}

func ShowItem(item Item, sp *MaxPrint) string {
	maxNameLen := sp.Name()
	maxPriceLen := sp.Price()
	maxStockLen := sp.Stock()
	return fmt.Sprintf("%-*s\tprice: NGN %-*s\tstock: %-*s", maxNameLen, item.Name(), maxPriceLen, e.Must(conv.Dec2StrFixed(item.Price())), maxStockLen, fmt.Sprintf("%d", item.Stock()))
}

func Showcase(items []Item, sp *MaxPrint) {
	slices.SortFunc(items, func(i, j Item) int {
		return strings.Compare(i.Name(), j.Name())
	})
	length := len(items)
	lengthStr := fmt.Sprintf("%d", length)
	digits := len(lengthStr)
	for i, item := range items {
		fmt.Printf("\n%-*s: \t%s\n", digits, fmt.Sprintf("%d", i+1), ShowItem(item, sp))
	}
}

func FromProds(prods []*Product) []Item {
	slices.SortFunc(prods, sort)
	items := make([]Item, len(prods))
	for i, prod := range prods {
		items[i] = Item(prod)
	}
	return items
}

func sort[I Item](i, j I) int {
	return strings.Compare(i.Name(), j.Name())
}

func FromOrders(orders []*Order) []Item {
	slices.SortFunc(orders, sort)
	items := make([]Item, len(orders))
	for i, order := range orders {
		items[i] = Item(order)
	}
	return items
}
