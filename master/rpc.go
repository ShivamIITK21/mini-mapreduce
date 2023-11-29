package master

import "github.com/ShivamIITK21/mini-mapreduce/core"

func (m *Master) RespondToTaskRequest(_ int, file_name *string) error {
	*file_name = ""
	for idx, task := range m.Tasks {
		if(task.Status == core.UNASSIGNED) {
			m.mu.Lock()
			if open_task, ok := m.Tasks[idx]; ok {
				open_task.Status = core.ASSIGNED
				m.Tasks[idx] = open_task
			}
			*file_name = m.Tasks[idx].File
			m.mu.Unlock()
		}
	}
	return nil
}