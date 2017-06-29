package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"context"
	"time"
)

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

func RootHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Many information about the blocks, such as possible blocked blocks, may also be used as a healthcheck")
}

func NewBlockHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("Creating a new block:", ps.ByName("name"))
	// if block exists, abort
	// if not, create
}
