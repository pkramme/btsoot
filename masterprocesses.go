package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func UpdateProcess(config Process) {
	log.Println("UPDATEPROC: Startup complete")
	Tick := time.NewTicker(120 * time.Second)
	for {
		select {
		case comm := <-config.Channel:
			if comm == StopCode {
				Tick.Stop()
				log.Println("UPDATEPROC: Shutdown")
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
	log.Println("WEBSERVERPROC: Startup complete")
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
				log.Println("WEBSERVERPROC: Shutdown")
				config.Channel <- ConfirmCode
				return
			}
		default:
			time.Sleep(100)
		}
	}
}

func ScanningProcess(config Process) {
	log.Println("SCANNERPROC: Startup complete")
	scanfilescomm := make(chan int, 2)
	go scanfiles(".", 4, scanfilescomm)
	for {
		select {
		case comm := <-config.Channel:
			if comm == StopCode {
				close(scanfilescomm)
				log.Println("SCANNERPROC: Shutdown")
				config.Channel <- ConfirmCode
				return
			}
		default:
			time.Sleep(100)
		}
	}
}
