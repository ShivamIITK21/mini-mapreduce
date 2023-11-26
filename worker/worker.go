package worker

type Worker struct {
	master	string
	Port	string
}

func New(Port string) *Worker{
	w := &Worker{Port: Port}
	return w
}

