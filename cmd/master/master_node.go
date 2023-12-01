package main

import (
	"flag"
	"fmt"
	"strings"
	"net/rpc"
	"net"
	"net/http"
	"log"

	"github.com/ShivamIITK21/mini-mapreduce/master"
)

func getArgs() (string, []string, []string) {
	var master_port string
	var worker_ports_str string
	var input_files_str	string

	flag.StringVar(&master_port, "m", "", "Port of master node")
	flag.StringVar(&worker_ports_str, "w", "", "Comma seperated list of worker ports")
	flag.StringVar(&input_files_str, "f", "", "Comma seperated list of input files")
	flag.Parse()
	worker_ports := strings.Split(worker_ports_str, ",")
	input_files := strings.Split(input_files_str, ",")

	return master_port, worker_ports, input_files
}

func main() {
	
	master_port, worker_ports, input_files := getArgs()
	m := master.New(master_port)

	rpc.Register(m)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", m.Port)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)

	
	m.CallAllWorkers(worker_ports)
	go m.PingAllWorkers()
	m.StoreMapTasks(input_files)
	fmt.Println(m.Tasks)

	for{}
}