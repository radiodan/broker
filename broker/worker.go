package broker

type Worker struct {
	Name     string     // human readable id
	Identity string     // routing frame
	Ready    bool       // ready to recieve work
	Queue    [][]string // pending messages to process
}

func (w *Worker) NextMsg() (msg []string, exists bool) {
	queueLength := len(w.Queue)
	if queueLength == 0 {
		exists = false
		return
	}

	exists = true
	msg = w.Queue[0]

	if queueLength == 1 {
		w.Queue = [][]string{}
	} else {
		w.Queue = w.Queue[1:]
	}

	return
}
