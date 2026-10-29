package wallets

import (
	"fmt"
	"log"
)

type AirDrop float64

const (
	ONE  AirDrop = 10_000
	TWO  AirDrop = 20_000
	FIVE AirDrop = 50_000
	TEN  AirDrop = 100_000
)

func (a AirDrop) String() string {
	switch a {
	case ONE:
		return "10,000"
	case TWO:
		return "20,000"
	case FIVE:
		return "50,000"
	case TEN:
		return "100,000"
	default:
		log.Fatal(fmt.Errorf("invalid AirDrop value: %d", int(a)))
	}
	return ""
}
