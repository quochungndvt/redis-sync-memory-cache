package rsmemory

import (
	"fmt"
	"log"
	"os"
)

func GetRunDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println(dir)
	return dir, nil
}
