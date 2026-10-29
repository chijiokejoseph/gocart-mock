package store

import (
	"fmt"
	"slices"
	"testing"
)

func TestShop(t *testing.T) {
	shop := BuildStore("../data")
	query := "golden penny spaghetti"
	prods, err := Find(query, shop)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	firstProd := prods[0]
	idx := slices.Index(shop.products, *firstProd)
	fmt.Printf("\nItem %s", ShowItem(firstProd, shop.MaxPrint()))
	storeProd, err := shop.Get(firstProd.Name())
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	idx2 := slices.Index(shop.products, *storeProd)
	if idx != idx2 {
		t.Errorf("Product does not match")
		return
	}
	cart := NewCart()
	cart.AddOrder(firstProd)
	fmt.Printf("\nItem at %d is %s", idx, ShowItem(firstProd, shop.MaxPrint()))
	fmt.Printf("\nItem %s", ShowItem(storeProd, shop.MaxPrint()))
	if firstProd.Name() != storeProd.Name() {
		t.Errorf("Product does not match")
	}
	if firstProd.Stock() != storeProd.Stock() {
		t.Errorf("Product does not match")
	}
}
