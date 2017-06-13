package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func UpdateProcess(config Process) {
	log.Printf("%d %d\tstarted", config.Level, UpdateThreadID)
	Tick := time.NewTicker(120 * time.Second)
	for {
		select {
		case comm := <-config.Channel:
			if comm == StopCode {
				Tick.Stop()
				config.Channel <- ConfirmCode
				return
			}
		default:
			select {
			case <-Tick.C:
				go log.Println("Update Check")
			default:
				time.Sleep(100) // Prevent high CPU usage
			}
		}
	}
}

func WebServer(config Process) {
	log.Printf("%d %d\tstarted", config.Level, WebserverThreadID)
	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "BTSOOT SERVER MAIN PAGE")
		}),
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	go func() {
		log.Println(server.ListenAndServe())
	}()
	for {
		select {
		case comm := <-config.Channel:
			if comm == StopCode {
				err := server.Shutdown(ctx)
				if err != nil {
					log.Println("HTTP error stop unsuccessful")
					config.Channel <- ErrorCode
					return
				}
				config.Channel <- ConfirmCode
				return
			}
		default:
			time.Sleep(100)
		}
	}
}

func ScanningProcess(config Process) {
	log.Printf("%d %d\tstarted", config.Level, ScanThreadID)
	scanfilescomm := make(chan int, 2)
	go scanfiles(".", scanfilescomm)
	for {
		select {
		case comm := <-config.Channel:
			if comm == StopCode {
				config.Channel <- ConfirmCode
				close(scanfilescomm)
				return
			}
		default:
			time.Sleep(100)

		}
	}
}
