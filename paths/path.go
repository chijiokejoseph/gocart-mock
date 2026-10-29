package paths

import (
	"log"
	"os"
	"path/filepath"

	e "github.com/chijiokejoseph/gocart-mock/err"
)

func Path() string {
	return e.Must(os.Getwd())
}

func DataPath() string {
	dir := Path()
	dataPath := filepath.Join(dir, "data")
	// check if a directory (not a file) specified as a string exists
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		// create the directory
		err := os.MkdirAll(dataPath, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return dataPath
}
