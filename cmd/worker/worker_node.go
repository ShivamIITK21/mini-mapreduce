package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"

	"github.com/ShivamIITK21/mini-mapreduce/worker"
)

func main() {
	ports := os.Args[1:]

	var wg sync.WaitGroup
	wg.Add(1)
	for _, port := range ports {
		worker := worker.New(port)
		rpc.Register(worker)
		rpc.HandleHTTP()
		l, err := net.Listen("tcp", worker.Port)
		if err != nil {
			log.Fatal("listen error:", err)
		}
		go http.Serve(l, nil)
	}
	wg.Wait()
}