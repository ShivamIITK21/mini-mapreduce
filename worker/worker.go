package worker

import (
	"log"
	"net/rpc"
	"sync"

	"github.com/ShivamIITK21/mini-mapreduce/core"
)

type Worker struct {
	master	string
	Port	string
	Wg		sync.WaitGroup
}

func New(Port string) *Worker{
	w := &Worker{Port: Port, master: "", Wg: sync.WaitGroup{}}
	w.Wg.Add(1)
	return w
}

func (w *Worker) AskForTask() (core.Task, error) {
	var recievedTask core.Task

	client, err := rpc.DialHTTP("tcp", w.master)
	if err != nil {
		log.Printf("Could not Dial to %s\n", w.master)
		return recievedTask, err
	}

	err = client.Call("Master.RespondToTaskRequest", w.Port, &recievedTask)
	if err != nil {
		log.Printf("Error in Calling %s\n", w.master)
		return recievedTask, err
	}

	return recievedTask, nil
}