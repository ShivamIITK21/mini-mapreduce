package master

import (
	"io"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
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
	TaskCounter		atomic.Int32
	NMap			int
	NReduce 		int
}

func New(port string, nreduce int) *Master{
	m := &Master{Port: port, Workers: make(map[string]WorkerInfo), Tasks: make(map[int]core.Task), NReduce: nreduce}
	return m
}

func (m *Master) CallWorker(port string) {
	client, err := rpc.DialHTTP("tcp", port)
	if err != nil {
		log.Printf("Could not Dial to %s\n", port)
		return
	}
	var status int
	err = client.Call("Worker.AssignMaster", core.SharedInfo{Port: m.Port, NReduce: m.NReduce, NMap: m.NMap}, &status)
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
	m.NMap = len(files)
	m.TaskCounter.Add(int32(len(files)))
	for idx, file := range files {
		m.Tasks[idx] = core.Task{File: file, Status: core.UNASSIGNED, Type: core.MAP, Id: idx}
	}
}

func (m *Master)CheckCompletion(){
	cnt := 0
	for {
		val := m.TaskCounter.Load()
		log.Print(val)
		if(val == 0) {
			if(cnt == 0){
				log.Printf("All Map tasks done, preparing for Reduce....")
				cnt++
				m.PrepareForReduce()
			} else{
				m.mu.Lock()
				m.CombineOutputs()
				log.Printf("MapReduce Completed, output is in output.txt")
				m.CleanUp()
				os.Exit(0)
				m.mu.Lock()
			}

		}
		time.Sleep(5*time.Second)
	}
}

func (m *Master)PrepareForReduce(){
	m.TaskCounter.Add(int32(m.NReduce))
	m.mu.Lock()
	for idx := 0; idx < m.NReduce; idx++ {	
		m.Tasks[m.NMap + idx] = core.Task{File: "", Status: core.UNASSIGNED, Type: core.REDUCE, Id: idx}	
	}
	m.mu.Unlock()
}


func (m *Master)CombineOutputs(){
	os.Remove("output.txt")
	out, err := os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create output file\n")
	}
	defer out.Close()

	for i := 0; i < m.NReduce; i++ {
		fname := "mr-ouput-" + strconv.Itoa(i) + ".txt"
		file, err := os.Open(fname)
		defer file.Close()
		if err != nil {
			log.Fatalf("Can't Open output file\n")
		}

		n, err := io.Copy(out, file)
		if err != nil {
			log.Fatalf("Unable to copy from reduce output\n")
		}
		log.Printf("Copied %d bytes from %s", n, fname)

	}
}

func (m *Master) CleanUp(){
	files, _ := filepath.Glob("mr*")
	for _, f := range files {
		os.Remove(f)
	}
}
