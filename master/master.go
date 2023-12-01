package master

import (
	"log"
	"net/rpc"
	"sync"
	"time"

	"github.com/ShivamIITK21/mini-mapreduce/core"
)

type WorkerInfo struct{
	Client	*rpc.Client
	Status	int
	Doing	int
}

type Master struct {
	Port    		string
	Workers 		map[string]WorkerInfo
	mu				sync.RWMutex
	Tasks			map[int]core.Task
}

func New(port string) *Master{
	m := &Master{Port: port, Workers: make(map[string]WorkerInfo), Tasks: make(map[int]core.Task)}
	return m
}

func (m *Master) CallWorker(port string) {
	client, err := rpc.DialHTTP("tcp", port)
	if err != nil {
		log.Printf("Could not Dial to %s\n", port)
		return
	}
	var status int
	err = client.Call("Worker.AssignMaster", m.Port, &status)
	if err != nil {
		log.Printf("Error in Calling %s\n", port)
		return
	}

	if status == 1{
		m.mu.Lock()
		if worker, ok := m.Workers[port]; ok {
			worker.Status = core.IDLE
			worker.Client = client
			worker.Doing = -1
			m.Workers[port] = worker
		} else {
			m.Workers[port] = WorkerInfo{Client: client, Status: core.IDLE, Doing: -1}
		}
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
	for {
		m.mu.RLock()
		client := m.Workers[port].Client
		m.mu.RUnlock()
		var status int
		err := client.Call("Worker.Ping", 0, &status)
		if err != nil || status != 1{
			log.Printf("Error in pinging %s\n", port)
			m.mu.Lock()
			if offline_worker, ok := m.Workers[port]; ok {
				offline_worker.Status = core.OFFLINE
				m.Workers[port] = offline_worker
			}
			m.mu.Unlock()
			m.CallWorker(port)
		} else {
			log.Printf("Pinged %s successfully\n", port)
		}
		time.Sleep(5*time.Second)
	}
}

func (m* Master)PingAllWorkers(){
	var ports []string
	m.mu.RLock()
	for port := range m.Workers {
		ports = append(ports, port)
	}
	m.mu.RUnlock()
	for _, port := range ports {
		p := port
		go m.pingWorkerPeriodically(p)
	}
}

func (m* Master)StoreMapTasks(files []string){
	for idx, file := range files {
		m.Tasks[idx] = core.Task{File: file, Status: core.UNASSIGNED, Type: core.MAP, Id: idx}
	}
}