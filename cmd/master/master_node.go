package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/ShivamIITK21/mini-mapreduce/master"
)

func main() {
	port := os.Args[1]
	m := master.New(port)

	var wg sync.WaitGroup
	wg.Add(1)
	m.CallAllWorkers(os.Args[2:])
	for _, w := range m.Workers {
		fmt.Printf("%s ", w.Port)
	}
	go m.PingAllWorkers()
	wg.Wait()
}