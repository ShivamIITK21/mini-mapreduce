package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	dir, err := os.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	out := ""
	for _, e := range dir {
		out = out + "data/" + e.Name() + ","
	}

	fmt.Println(out)
}