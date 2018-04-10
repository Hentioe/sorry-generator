package main

import (
	"os"
	"path/filepath"
	"log"
	"fmt"
)

func init() {
	if parentPath, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		log.Fatal(err)
	} else {
		distDir := parentPath + "/dist"
		if _, err := os.Stat(distDir); os.IsNotExist(err) {
			fmt.Println(distDir)
			os.Mkdir(distDir, os.FileMode(0774))
		}
	}
}

func main() {
}
