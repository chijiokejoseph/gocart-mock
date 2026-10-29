package wallets

import (
	"fmt"
	"maps"

	"github.com/chijiokejoseph/gocart-mock/conv"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	balance       decimal.Decimal
	creditHistory map[AirDrop]int
	debitHistory  []decimal.Decimal
	cache         *Wallet
}

// NewWallet : Wallet Constructor
func NewWallet() *Wallet {
	return &Wallet{creditHistory: make(map[AirDrop]int)}
}

// Wallet Field Receivers

func (w *Wallet) Balance() decimal.Decimal {
	return w.balance
}

func (w *Wallet) CreditHistory() map[AirDrop]int {
	return w.creditHistory
}

func (w *Wallet) DebitHistory() []decimal.Decimal {
	return w.debitHistory
}

func (w *Wallet) Fund(amount AirDrop) {
	creditKeys := maps.Keys(w.creditHistory)
	for key := range creditKeys {
		if key == amount {
			w.creditHistory[key] += 1
			w.balance = w.balance.Add(conv.Int2Dec(int(amount)))
			return
		}
	}
	maps.Insert(w.creditHistory, func(yield func(AirDrop, int) bool) {
		yield(amount, 1)
	})
	w.balance = w.balance.Add(conv.Int2Dec(int(amount)))
}

var InadequateFundsErr = fmt.Errorf("not enough balance to pay for this amount")

func (w *Wallet) Charge(amount decimal.Decimal) error {
	if w.balance.LessThan(amount) {
		fmt.Printf("\n%s\n", InadequateFundsErr)
		return InadequateFundsErr
	}
	w.balance = w.balance.Sub(amount)
	w.debitHistory = append(w.debitHistory, amount)
	return nil
}

func (w *Wallet) Clear() {
	w.cache = &Wallet{balance: w.balance, creditHistory: w.creditHistory, debitHistory: w.debitHistory, cache: nil}
	w.creditHistory = make(map[AirDrop]int)
	w.debitHistory = make([]decimal.Decimal, 0)
	w.balance = decimal.NewFromInt(0)
}

func (w *Wallet) Unclear() {
	if w.cache == nil {
		return
	}
	oldWallet := w.cache
	*w = *oldWallet
}
