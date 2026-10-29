package store

import (
	"log"

	"github.com/shopspring/decimal"
)

type Cart struct {
	orders   []*Order
	products []*Product
	total    decimal.Decimal
}

func NewCart() *Cart {
	return &Cart{orders: []*Order{}}
}

func (c *Cart) Orders() []*Order {
	return c.orders
}

func (c *Cart) Total() decimal.Decimal {
	c.total = decimal.Zero
	for _, order := range c.orders {
		c.total = c.total.Add(order.CalcPrice())
	}
	return c.total
}

func (c *Cart) AddOrder(product *Product) {
	for _, order := range c.Orders() {
		if order.Name() == product.Name() {
			order.AddItem()
			product.RemoveItem()
			return
		}
	}
	newOrder := NewOrder(product.Name(), product.Price(), 1)
	c.orders = append(c.orders, &newOrder)
	c.products = append(c.products, product)
	product.RemoveItem()
}

func (c *Cart) AddOrders(products []*Product) {
	for _, prod := range products {
		c.AddOrder(prod)
	}
}

func (c *Cart) RemoveOrder(prevOrder *Order) {
	for i, order := range c.orders {
		if order.Name() == prevOrder.Name() {
			order.RemoveItem()
			prod := c.products[i]
			prod.AddItem()
			if order.Stock() == 0 {
				c.orders = append(c.orders[:i], c.orders[i+1:]...)
				c.products = append(c.products[:i], c.products[i+1:]...)
			}
			return
		}
	}
	log.Fatalf("Item %s not found in cart\n", prevOrder.Name())
}

func (c *Cart) RemoveOrders(items []*Order) {
	for _, item := range items {
		c.RemoveOrder(item)
	}
}

func (c *Cart) Clear() {
	for i, order := range c.orders {
		prod := c.products[i]
		for order.Stock() > 0 {
			order.RemoveItem()
			prod.AddItem()
		}
	}
	c.orders = []*Order{}
	c.products = []*Product{}
	c.total = decimal.Zero
}

func (c *Cart) ClearPostPay() {
	c.orders = []*Order{}
	c.products = []*Product{}
	c.total = decimal.Zero
}
