package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paulkramme/toml"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("BTSOOT - Copyright (c) 2017 Paul Kramme All Rights Reserved.")

	ConfigLocation := flag.String("config", "./btsoot.conf", "Specifies configfile location")
	flag.Parse()

	var Config Configuration
	_, err := toml.DecodeFile(*ConfigLocation, &Config)
	if err != nil {
		fmt.Println("Couldn't find or open config file.")
		panic(err)
	}
	db, err := sql.Open("sqlite3", "./helloworld.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ProcessList := CreateMasterProcessList()

	// NOTE: Init standard threads...
	go UpdateProcess(ProcessList[UpdateThreadID])
	go WebServer(ProcessList[WebserverThreadID])
	go ScanningProcess(ProcessList[ScanThreadID])
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT)
	<-signals
	fmt.Println("Exiting. Please wait...")

	for k, v := range ProcessList {
		fmt.Printf("Sending stop (%x) to THREADID=%d\n", StopCode, k)
		v.Channel <- StopCode
	}
	for k, v := range ProcessList {
		callback := <-v.Channel
		if callback == ConfirmCode {
			fmt.Printf("THREADID=%d: Confirmation (%x)\n", k, ConfirmCode)
		} else {
			fmt.Println("Problems...")
			// FIXME: Holy, please make this a select
		}
	}
}

func UpdateProcess(config Process) {
	log.Printf("%d %d\tstarted", config.Level, UpdateThreadID)
	for {

		select {
		case comm := <-config.Channel:
			if comm == StopCode {
				config.Channel <- ConfirmCode
				return
			}
		default:
			time.Sleep(10 * time.Second)
		}
	}
}

func WebServer(config Process) {
	log.Printf("%d %d\tstarted", config.Level, WebserverThreadID)
	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello World!")
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
			time.Sleep(1 * time.Second)
		}
	}
}

func ScanningProcess(config Process) {
	log.Printf("%d %d\tstarted", config.Level, ScanThreadID)
	for {
		select {
		case comm := <-config.Channel:
			if comm == StopCode {
				config.Channel <- ConfirmCode
				return
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
