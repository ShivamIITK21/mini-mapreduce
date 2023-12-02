## Mini-MapReduce

This is a small implementation of [Google's MapReduce(2004) Paper](https://static.googleusercontent.com/media/research.google.com/en//archive/mapreduce-osdi04.pdf).

### Architecture

There is one master and muliple worker nodes, which are supposed to perform map and reduce tasks. The master and worker nodes communicate via RPCs. The master node is responsible for managing the workers and assining tasks to the workers, combination and cleanup of the intermediate files. The master node uses several threads for this. The worker nodes on the other hand are mostly single threaded, except a few threads which respond to pings and RPCs. This architecture requires a shared filesystem to work.

### Build and Run

Write you map and reduce functions in a single Go file and compile it into a shared lib. For example, 

```
go build -buildmode=plugin ./examples/wordcount/wc.go
```

Start any number of worker nodes by speciying the port and the Go plugin.

```
go run cmd/worker/worker_node.go :3002 indexer.so
```

Finally, start the master node by specifying the port, all the worker ports, the input files and the number of reduce tasks.

```
go run -race  cmd/master/master_node.go -m :3000 -w :3001,:3002,:3003 -f data/1.txt,data/2.txt,data/3.txt -r 5
```

PS - Some Fault tolerance features and error handling are yet to be added