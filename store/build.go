package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	e "github.com/chijiokejoseph/gocart-mock/err"
	"github.com/shopspring/decimal"
)

type Data struct {
	count int
	price decimal.Decimal
}

func BuildStore(dataPath string) *Store {
	srcNos := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	inventory := make(map[string]*Data)
	total := 0
	for _, srcNo := range srcNos {
		src := filepath.Join(dataPath, fmt.Sprintf("warehouse%d.txt", srcNo))
		content, err := os.ReadFile(src)
		if err != nil {
			log.Fatalln(err)
		}
		contents := strings.ReplaceAll(string(content), "\r\n", "\n")
		contents = strings.ReplaceAll(contents, "\n", "")
		contents = strings.TrimSuffix(contents, ";")
		itemContents := strings.Split(contents, ";")
		for _, itemContent := range itemContents {
			itemContent = strings.TrimSpace(itemContent)
			details := strings.Split(itemContent, ":")
			name := details[0]
			price := details[1]
			name = strings.TrimSpace(name)
			price = strings.TrimSpace(price)
			if entry, exists := inventory[name]; exists {
				entry.count++
				continue
			}
			inventory[name] = &Data{count: 1, price: e.Must(decimal.NewFromString(price))}
			total++
		}
	}
	shelf := make([]Product, total)
	i := 0
	for name, data := range inventory {
		shelf[i] = NewProduct(name, data.price, data.count)
		i++
	}
	return NewStore(shelf)
}
