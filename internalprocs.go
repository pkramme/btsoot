package main

import (
	"fmt"
	"sync"
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	p.Channel <- StopCode
	for {
		select {
		case callback := <-p.Channel:
			if callback == ErrorCode {
				fmt.Println("Error (%x)\nYou may have to kill btsoot. Shutdown is now unsafe.", ErrorCode)
			}
			time.Sleep(1 * time.Second)
			wg.Done()
		default:
			select {
			case <-ctx.Done():
				fmt.Println("One thread does not answer. Program has to be killed manually.")
				fmt.Println(ctx.Err())
			default:
				// NOTE: Wait for the next loop
			}
		}

	}
}

func KillAll(m map[int]Process) {
	var wg sync.WaitGroup
	for _, v := range m {
		wg.Add(1)
		go v.Kill(&wg)
	}
	wg.Wait()
	fmt.Println("Everything is shutted down.")
}
