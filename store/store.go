package store

import (
	"fmt"
	"slices"
	"strings"

	"github.com/chijiokejoseph/gocart-mock/conv"
	"github.com/chijiokejoseph/gocart-mock/err"
)

type MaxPrint struct {
	name  int
	price int
	stock int
}

func (sp *MaxPrint) Name() int {
	return sp.name
}

func (sp *MaxPrint) Price() int {
	return sp.price
}

func (sp *MaxPrint) Stock() int {
	return sp.stock
}

type Store struct {
	products []Product
}

func NewStore(products []Product) *Store {
	slices.SortFunc(products, func(a, b Product) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return &Store{products: products}
}

func (store *Store) Inventory() []*Product {
	prods := make([]*Product, len(store.products))
	for i := range store.products {
		prods[i] = &store.products[i]
	}
	return prods
}

func (store *Store) Names() []string {
	names := make([]string, len(store.products))
	for i, prod := range store.products {
		names[i] = prod.Name()
	}
	return names
}

func (store *Store) Get(name string) (*Product, error) {
	for _, prod := range store.Inventory() {
		if prod.Name() == name {
			return prod, nil
		}
	}
	return nil, ErrNoItemFound
}

func (store *Store) MaxPrint() *MaxPrint {
	return &MaxPrint{
		name:  store.MaxNameLen(),
		price: store.MaxPriceLen(),
		stock: store.MaxStockLen(),
	}
}

func (store *Store) MaxNameLen() int {
	maxLen := 0
	for _, prod := range store.products {
		if len(prod.Name()) > maxLen {
			maxLen = len(prod.Name())
		}
	}
	return maxLen
}

func (store *Store) MaxPriceLen() int {
	maxLen := 0
	for _, prod := range store.products {
		priceStr := err.Must(conv.Dec2StrFixed(prod.Price()))
		if len(priceStr) > maxLen {
			maxLen = len(priceStr)
		}
	}
	return maxLen
}

func (store *Store) MaxStockLen() int {
	maxLen := 0
	for _, prod := range store.products {
		if len(fmt.Sprintf("%d", prod.Stock())) > maxLen {
			maxLen = len(fmt.Sprintf("%d", prod.Stock()))
		}
	}
	return maxLen
}
