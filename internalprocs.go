package main

import (
	"fmt"
	"log"
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
	Channel      chan interface{}
	Subprocesses map[int]Process
}

func CreateMasterProcessList() map[int]Process {
	pmap := make(map[int]Process)

	pmap[0] = Process{
		Channel: make(chan interface{}),
	}
	pmap[1] = Process{
		Channel: make(chan interface{}),
	}
	pmap[2] = Process{
		Channel: make(chan interface{}),
	}
	return pmap
}

func (p Process) Kill() {
	p.Channel <- StopCode
	// NOTE: Wait 10 seconds, which is 100 loops with 100 milliseconds delay
	// to increase responsiveness
	for i := 100; i > 0; i-- {
		select {
		case callback := <-p.Channel:
			if callback == ErrorCode {
				log.Println("Error %x. Could not kill thread", ErrorCode)
			}
		default:
			time.Sleep(100)
			// NOTE: Wait for the next loop
		}
	}
}

// This function only kills all level 0 threads. Subthreads need to be handled
// by the threads itself.
func KillAll(m map[int]Process) {
	go func() {
		fmt.Println("Waiting 10 seconds...")
		time.Sleep(10 * time.Second)
		fmt.Println("One or more threads do not answer. You may have to kill the program or wait.")
	}()
	var wg sync.WaitGroup
	for _, v := range m {
		wg.Add(1)
		go func(p Process, wg *sync.WaitGroup) {
			p.Channel <- StopCode

			for {
				select {
				case callback := <-p.Channel:
					if callback == ErrorCode {
						fmt.Println("Error (%x)\nOne thread did not answer.", ErrorCode)
					}
					wg.Done()
				default:
					time.Sleep(100)
					// NOTE: Wait for the next loop
				}
			}
		}(v, &wg)
	}
	wg.Wait()
	fmt.Println("Everything is shutted down.")
}
