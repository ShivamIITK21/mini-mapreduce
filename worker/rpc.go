package worker

import (
	"log"

	"github.com/ShivamIITK21/mini-mapreduce/core"
)

func(w *Worker) AssignMaster(info core.SharedInfo, status *int) error {
	*status = 1
	w.master = info.Port
	w.Nreduce = info.NReduce
	w.Wg.Done()
	log.Printf("%s got assigned a master\n", w.Port)
	return nil
}

func(w *Worker) Ping(_ int, status *int) error {
	*status = 1
	log.Printf("%s got pinged\n", w.Port)
	return nil
}
