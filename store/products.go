package store

import (
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/conv"
	e "github.com/chijiokejoseph/gocart-mock/err"
	"github.com/shopspring/decimal"
)

type Product struct {
	name  string
	price decimal.Decimal
	stock int
}

func (prod *Product) Name() string {
	return prod.name
}

func (prod *Product) Price() decimal.Decimal {
	return prod.price
}

func (prod *Product) Stock() int {
	return prod.stock
}

func (prod *Product) CalcPrice() decimal.Decimal {
	return prod.price.Mul(conv.Int2Dec(prod.stock))
}

func (prod *Product) AddItem() {
	prod.stock++
}

func (prod *Product) RemoveItem() {
	prod.stock--
}

func (prod *Product) String() string {
	return fmt.Sprintf("%s\tprice = %sNGN\tstock = %d", prod.name, e.Must(conv.Dec2StrFixed(prod.price)), prod.stock)
}

func NewProduct(name string, price decimal.Decimal, stock int) Product {
	return Product{name: name, price: price, stock: stock}
}
