package store

import (
	"fmt"

	"github.com/chijiokejoseph/gocart-mock/conv"
	e "github.com/chijiokejoseph/gocart-mock/err"
	"github.com/shopspring/decimal"
)

type Order struct {
	name  string
	price decimal.Decimal
	stock int
}

func (order *Order) Name() string {
	return order.name
}

func (order *Order) Price() decimal.Decimal {
	return order.price
}

func (order *Order) Stock() int {
	return order.stock
}

func (order *Order) CalcPrice() decimal.Decimal {
	return order.price.Mul(conv.Int2Dec(order.stock))
}

func (order *Order) AddItem() {
	order.stock++
}

func (order *Order) RemoveItem() {
	order.stock--
}

func (order *Order) String() string {
	return fmt.Sprintf("%s\tprice = %sNGN\tstock = %d", order.name, e.Must(conv.Dec2StrFixed(order.price)), order.stock)
}

func NewOrder(name string, price decimal.Decimal, stock int) Order {
	return Order{name: name, price: price, stock: stock}
}
