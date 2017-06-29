package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func WebServer(procconfig Process, config Configuration, db *sql.DB) {
	log.Println("WEBSERVERPROC: Startup complete")
	router := httprouter.New()
	server := http.Server{
		Addr:    config.Listen,
		Handler: router,
	}
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		RootHandler(w, r, ps, db)
	})
	router.POST("/block/:name", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		NewBlockHandler(w, r, ps, db)
	})

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

func RootHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *sql.DB) {
	fmt.Fprintf(w, "Many information about the blocks, such as possible blocked blocks, may also be used as a healthcheck")

}

func NewBlockHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *sql.DB) {
	fmt.Println("Creating a new block:", ps.ByName("name"))
	// if block exists, abort
	// if not, create
}
