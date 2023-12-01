package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/ShivamIITK21/mini-mapreduce/worker"
)

func main() {
	port := os.Args[1]
	pluginName := os.Args[2]

	worker := worker.New(port, pluginName)
	rpc.Register(worker)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", worker.Port)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)

	worker.Wg.Wait()

	for {
		task, err := worker.AskForTask()
		if err != nil {
			log.Printf("Error in fetching task\n")
		}
		if(task.File == ""){
			log.Printf("No Open task was available\n")
		} else {
			log.Printf("Recived task on %s\n", task.File)
			err := worker.DoTask(task)
			log.Print(err)
		}
		time.Sleep(1*time.Second)
	}

}