package worker

import "log"

func(w *Worker) AssignMaster(m_port string, status *int) error {
	*status = 1
	w.master = m_port
	log.Printf("%s got assigned a master\n", w.Port)
	return nil
}

func(w *Worker) Ping(_ int, status *int) error {
	*status = 1
	log.Printf("%s got pinged\n", w.Port)
	return nil
}
