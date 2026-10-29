package store_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/chijiokejoseph/gocart-mock/store"
)

func setup() *store.Store {
	shop := store.BuildStore("../data")
	return shop
}

func TestFind(t *testing.T) {
	query := "Spaghetti golden penny"
	shop := setup()
	prods, err := store.Find(query, shop)
	if err != nil {
		if errors.Is(err, store.ErrQueryInvalid) {
			t.Errorf("%s", err)
			return
		}
		return
	}
	items := store.FromProds(prods)
	store.Showcase(items, shop.MaxPrint())
}

func TestAddToStore(t *testing.T) {
	query := "golden penny spaghetti"
	shop := setup()
	prods, err := store.Find(query, shop)
	if err != nil {
		if errors.Is(err, store.ErrQueryInvalid) {
			t.Errorf("%s", err)
			return
		}
		if errors.Is(err, store.ErrNoItemFound) {
			fmt.Println("No item found matching query")
			return
		}
	}
	items := store.FromProds(prods)
	store.Showcase(items, shop.MaxPrint())
	firstProd := prods[1]
	cart := store.NewCart()
	cart.AddOrder(firstProd)
	fmt.Printf("\nItem %s", store.ShowItem(firstProd, shop.MaxPrint()))
	storeProd, err := shop.Get(firstProd.Name())
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	fmt.Printf("\nStore Item %s", store.ShowItem(storeProd, shop.MaxPrint()))
	if storeProd.Name() != firstProd.Name() {
		t.Errorf("Product does not match")
	}
	//prod := store.NewProduct("test", decimal.NewFromInt(22), 2)
	//set := []store.Product{prod, prod, prod, prod}
	//setItem1 := &set[0]
	//setItem1.RemoveItem()
	//fmt.Printf("\nProduct : %s", store.ShowItem(setItem1, shop.MaxPrint()))
	//fmt.Printf("\nProduct - %s", store.ShowItem(&set[0], shop.MaxPrint()))
	if storeProd.Stock() != firstProd.Stock() {
		t.Errorf("Product does not match")
	}
}
