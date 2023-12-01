package worker

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"sync"

	"github.com/ShivamIITK21/mini-mapreduce/core"
)

type Worker struct {
	master	string
	Port	string
	Wg		sync.WaitGroup
	Map		func(string, string) []core.KeyValue
	Reduce	func(string, []string) string
}

func New(Port string, pluginFileName string) *Worker{
	w := &Worker{Port: Port, master: "", Wg: sync.WaitGroup{}}
	mapf, reducef := core.ReadMapReduceFuncs(pluginFileName)
	w.Map = mapf
	w.Reduce = reducef
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

func (w *Worker) DoTask(task core.Task) error {
	if(task.Type == core.MAP) {
		file, err := os.ReadFile(task.File)
		if err != nil {
			return err
		}
		fileContent := string(file)
		kva := w.Map(task.File, fileContent)

		oname := "mr-" + strconv.Itoa(task.Id) + ".txt"
		os.Remove(oname)
		oFile, err := os.Create(oname)
		if err != nil {
			return err
		}

		for _, kv := range kva {
			fmt.Fprintf(oFile, "%s %s\n", kv.Key, kv.Value)
		}

	} else if(task.Type == core.REDUCE) {

	}

	return nil
}