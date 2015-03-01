package broker

type Worker struct {
	Name     string     // human readable id
	Identity string     // routing frame
	Ready    bool       // ready to recieve work
	Queue    [][]string // pending messages to process
}
