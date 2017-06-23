package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"github.com/julienschmidt/httprouter"
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
	router := httprouter.New()
	server := http.Server{
		Addr: ":8080",
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
	config.Subprocesses = make(map[int]Process)
	//go scanfiles(".", 4, scanfilescomm)
	for {
		select {
		case comm := <-config.Channel:
			if comm == StopCode {
				log.Println("SCANNERPROC: Shutdown")
				config.Channel <- ConfirmCode
				return
			}
		default:
			time.Sleep(100)
		}
	}
}
