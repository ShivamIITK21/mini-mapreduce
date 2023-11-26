package worker

import "log"

func(w *Worker) AssignMaster(m_port string, status *int) error {
	*status = 1
	w.master = m_port
	log.Printf("%s got assigned a master\n", w.Port)
	return nil
}