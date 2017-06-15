package main

import (
	"flag"
	"fmt"
	"github.com/paulkramme/toml"
	"log"
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

	f, err := os.OpenFile(Config.LogFileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	log.SetOutput(f)

	// NOTE: Create a spacer in the log
		log.Println("\n\n\n\n\n")

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
	log.Println("Exiting...")
	KillAll(ProcessList)
	time.Sleep(1 * time.Second) // wait for scanners
}
