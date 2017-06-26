package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
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

func WebServer(procconfig Process, config Configuration) {
	log.Println("WEBSERVERPROC: Startup complete")
	router := httprouter.New()
	server := http.Server{
		Addr:    config.Listen,
		Handler: router,
	}
	router.GET("/", RootHandler)
	router.POST("/block/:name", NewBlockHandler)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	go func() {
		log.Println(server.ListenAndServe())
	}()
	for {
		select {
		case comm := <-procconfig.Channel:
			if comm == StopCode {
				err := server.Shutdown(ctx)
				if err != nil {
					log.Println("HTTP error stop unsuccessful")
					procconfig.Channel <- ErrorCode
					return
				}
				log.Println("WEBSERVERPROC: Shutdown")
				procconfig.Channel <- ConfirmCode
				return
			}
		default:
			time.Sleep(100)
		}
	}
}

func ScanningProcess(procconfig Process, config Configuration) {
	log.Println("SCANNERPROC: Startup complete")
	procconfig.Subprocesses = make(map[int]Process)
	//go scanfiles(".", 4, scanfilescomm)
	for {
		select {
		case comm := <-procconfig.Channel:
			if comm == StopCode {
				log.Println("SCANNERPROC: Shutdown")
				procconfig.Channel <- ConfirmCode
				return
			}
		default:
			time.Sleep(100)
		}
	}
}
