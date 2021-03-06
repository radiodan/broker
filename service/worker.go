package service

import (
	"fmt"
	"time"
)

type Worker struct {
	Name     string     // Human readable id
	Identity string     // Routing frame
	Ready    bool       // Ready to recieve work
	Queue    []*Request // Pending Requests to process
	Services Services   // Array of registered services
	Expiry   time.Time  // Expires at unless heartbeat
}

func NewWorker(identity string, services Services) *Worker {
	name := fmt.Sprintf("%q", identity)

	worker := &Worker{
		Identity: identity,
		Name:     name,
		Ready:    true,
		Services: services,
	}

	worker.Refresh()

	return worker
}

func (w *Worker) AppendToQueue(req *Request) {
	w.Queue = append([]*Request{req}, w.Queue...)
}

func (w *Worker) Refresh() {
	w.Expiry = time.Now().Add(HEARTBEAT_EXPIRY)
}

func (w *Worker) NextMsg() (msg *Request, exists bool) {
	queueLength := len(w.Queue)
	if queueLength == 0 {
		exists = false
		return
	}

	exists = true
	msg = w.Queue[0]

	if queueLength == 1 {
		w.Queue = []*Request{}
	} else {
		w.Queue = w.Queue[1:]
	}

	return
}
