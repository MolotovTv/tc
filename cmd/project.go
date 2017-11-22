package cmd

import (
	"log"
	"os"
	"path"
)

func projectName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return path.Base(dir)
}
