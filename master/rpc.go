package master

import "github.com/ShivamIITK21/mini-mapreduce/core"

func (m *Master) RespondToTaskRequest(worker_port string, task_ptr *core.Task) error {
	task_ptr.File = ""
	m.mu.Lock()
	for idx, task := range m.Tasks {
		if(task.Status == core.UNASSIGNED) {
			if open_task, ok := m.Tasks[idx]; ok {
				open_task.Status = core.ASSIGNED
				m.Tasks[idx] = open_task
				task_ptr.File = open_task.File
				task_ptr.Type = open_task.Type
				task_ptr.Id = open_task.Id
				
				if worker, ok := m.Workers[worker_port]; ok {
					worker.Doing = idx
					m.Workers[worker_port] = worker
				}
			}
			break
		}
	}
	m.mu.Unlock()
	return nil
}