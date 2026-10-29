package store_test

import (
	"testing"

	s "github.com/chijiokejoseph/gocart-mock/store"
)

func TestBuildStore(t *testing.T) {
	store := s.BuildStore("../data")
	t.Logf("len(store) = %d", len(store.Inventory()))
	sp := store.MaxPrint()
	for i, prod := range store.Inventory() {
		var item s.Item = prod
		t.Logf("\nItem %d: \t%s\n", i+1, s.ShowItem(item, sp))
		if i == 20 {
			break
		}
	}
}

func TestShowCase(t *testing.T) {
	store := s.BuildStore("../data")
	sp := store.MaxPrint()
	itemPtrs := s.FromProds(store.Inventory())
	s.Showcase(itemPtrs[:100], sp)

}
