package master

import (
	"log"
	"net/rpc"
	"sync"
	"time"
)

type WorkerInfo struct{
	Port	string
	Client	*rpc.Client
}

type Master struct {
	port    		string
	Workers 		[]*WorkerInfo
	WorkerClient 	map[string]*rpc.Client
	mu				sync.RWMutex
}

func New(port string) *Master{
	m := &Master{port: port, WorkerClient: make(map[string]*rpc.Client)}
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
		m.WorkerClient[port] = client
		log.Printf("Got a reply from %s\n", port)
		m.mu.Unlock()
		return
	} else {
		log.Printf("Worker at %s is not online\n", port)
		return
	}
}

func (m *Master)CallAllWorkers(ports []string){
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		p := port
		go func ()  {
			defer wg.Done()
			m.CallWorker(p)
		}()
	}

	wg.Wait()
}

func (m *Master)pingWorkerPeriodically(port string){
	client := m.WorkerClient[port]
	for {
		var status int
		err := client.Call("Worker.Ping", 0, &status)
		if err != nil || status != 1{
			log.Printf("Error in pinging %s\n", port)
			m.mu.Lock()
			delete(m.WorkerClient, port)
			m.mu.Unlock()
			return
		} else {
			log.Printf("Pinged %s successfully\n", port)
			time.Sleep(5*time.Second)
		}
	}
}

func (m* Master)PingAllWorkers(){
	for port := range m.WorkerClient {
		p := port
		go m.pingWorkerPeriodically(p)
	}
}