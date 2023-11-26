package master

import (
	"log"
	"net/rpc"
	"sync"
)

type WorkerInfo struct{
	Port	string
	Client	*rpc.Client
}

type Master struct {
	port    string
	Workers []*WorkerInfo
	mu		sync.RWMutex
}

func New(port string) *Master{
	m := &Master{port: port}
	return m
}

func (m *Master) CallWorker(port string) {
	client, err := rpc.DialHTTP("tcp", port)
	if err != nil {
		log.Printf("Could not Dial to %s\n", port)
		return
	}
	var status int
	err = client.Call("Worker.AssignMaster", m.port, &status)
	if err != nil {
		log.Printf("Error in Calling %s\n", port)
		return
	}

	if status == 1{
		m.mu.Lock()
		m.Workers = append(m.Workers, &WorkerInfo{Port: port, Client: client})
		log.Printf("Got a repy from %s\n", port)
		m.mu.Unlock()
		return
	} else {
		log.Printf("Worker at %s is not online\n", port)
		return
	}
}
