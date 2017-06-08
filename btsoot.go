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

	// NOTE: Wait for SIGINT
	signal.Notify(signals, syscall.SIGINT)
	<-signals
	fmt.Println("Exiting. Please wait...")

	// NOTE: Copy ProcessList
	ContactList := make(map[int]Process)
	for k, v := range ProcessList {
		ContactList[k] = v
	}

	// NOTE: Sending StopCode to registered threads inside ProcessList
	//ContactList := ProcessList
	for len(ContactList) != 0 {
		fmt.Println("Masterloop")
		for ContactListKey, _ := range ContactList {
			fmt.Println("Shortloop")
			v := ProcessList[ContactListKey]
			select {
			case v.Channel <- StopCode:
				fmt.Printf("Sending stop (%x) to THREADID=%d\n", StopCode, ContactListKey)
				delete(ContactList, ContactListKey)
			default:
				// NOTE: Wait for the next flyby to contact thread
			}
		}
	}

	fmt.Println(ProcessList)
	// NOTE: Receiving reponse code and deleting registration from ProcessList
	for len(ProcessList) != 0 {
		//fmt.Println(len(ProcessList), "threads are registered.")
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
	fmt.Print(ProcessList)
}
