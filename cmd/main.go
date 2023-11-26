package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/ShivamIITK21/mini-mapreduce/master"
	"github.com/ShivamIITK21/mini-mapreduce/worker"
)

func main(){
	master := master.New(":3000")
	worker := worker.New(":3001")
	rpc.Register(worker)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", worker.Port)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)

	master.CallWorker(":3001")
	fmt.Println(master.Workers[0].Port)

}