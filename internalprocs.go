package main

const (
	UpdateThreadID    = 0
	ScanThreadID      = 1
	WebserverThreadID = 2

	StopCode    = 1000
	ConfirmCode = 1001
	ErrorCode   = 1002
)

func CreateMasterProcessList() map[int]Process {
	pmap := make(map[int]Process)

	pmap[0] = Process{
		Channel: make(chan int),
		Level:   0,
	}
	pmap[1] = Process{
		Channel:      make(chan int),
		Level:        0,
		Subprocesses: make(map[int]Process),
	}
	pmap[2] = Process{
		Channel: make(chan int),
		Level:   0,
	}
	return pmap
}

func ListProcessList() {
	return
}

func (p Process) AddProcessToList() {
	return
}

func (p Process) Kill() {
	return
}

type Process struct {
	Channel      chan int
	Level        int
	Subprocesses map[int]Process
}
