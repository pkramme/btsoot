package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paulkramme/toml"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	db, err := sql.Open("sqlite3", Config.DBFileLocation)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	ProcessList := CreateMasterProcessList()

	// NOTE: Init standard threads...
	go UpdateProcess(ProcessList[UpdateThreadID])
	go WebServer(ProcessList[WebserverThreadID])
	go ScanningProcess(ProcessList[ScanThreadID])
	signals := make(chan os.Signal, 1)

	log.Println("Startup complete...")

	signal.Notify(signals, syscall.SIGINT)
	<-signals
	fmt.Println("Exiting. Please wait...")

	for k, v := range ProcessList {
		fmt.Printf("Sending stop (%x) to THREADID=%d\n", StopCode, k)
		v.Channel <- StopCode
	}
	for len(ProcessList) != 0 {
		for k, v := range ProcessList {
			select {
			case callback := <-v.Channel:
				if callback == ConfirmCode {
					fmt.Printf("THREADID=%d: Confirmation (%x)\n", k, ConfirmCode)
					delete(ProcessList, k)
				} else if callback == ErrorCode {
					fmt.Println("THREADID=%d: Error (%x)\n", k, ErrorCode)
				}
			default:
				// NOTE: Thread did not respond, wait for the next flyby
			}
		}
	}
}
