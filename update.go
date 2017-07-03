package main

import (
	"log"
	"time"
)

func UpdateProcess(procconfig Process, config Configuration) {
	log.Println("UPDATEPROC: Startup complete")
	Tick := time.NewTicker(120 * time.Second)
	for {
		select {
		case comm := <-procconfig.Channel:
			if comm == StopCode {
				Tick.Stop()
				log.Println("UPDATEPROC: Shutdown")
				procconfig.Channel <- ConfirmCode
				return
			}
		default:
			select {
			case <-Tick.C:
				go log.Println("UPDATEPROC: Update Check")
			default:
				time.Sleep(100) // Prevent high CPU usage
			}
		}
	}
}
