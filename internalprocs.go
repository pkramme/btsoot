package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	UpdateThreadID    = 0
	ScanThreadID      = 1
	WebserverThreadID = 2

	StopCode    = 1000
	ConfirmCode = 1001
	ErrorCode   = 1002
)

type Process struct {
	Channel      chan int
	Level        int
	Subprocesses map[int]Process
}

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

func (p Process) Kill(wg *sync.WaitGroup) {
	p.Channel <- StopCode

	for {
		select {
		case callback := <-p.Channel:
			if callback == ErrorCode {
				fmt.Println("Error (%x)\nYou may have to kill btsoot. Shutdown is now unsafe.", ErrorCode)
			}
			wg.Done()
		default:
			time.Sleep(100)
			// NOTE: Wait for the next loop
		}
	}
}

func KillAll(m map[int]Process) {
	go func() {
		fmt.Println("Waiting 10 seconds...")
		time.Sleep(10 * time.Second)
		fmt.Println("One or more threads do not answer. You may have to kill the program or wait.")
	}()
	var wg sync.WaitGroup
	for _, v := range m {
		wg.Add(1)
		go v.Kill(&wg)
	}
	wg.Wait()
	fmt.Println("Everything is shutted down.")
}
