package main

import (
	"flag"
	"fmt"
	"github.com/paulkramme/ini"
	"log"
	"os"
	"time"
)

const (
	StopCode    = 1000
	ConfirmCode = 1001
	ErrorCode   = 1002
)

func main() {
	fmt.Println("BTSOOT - Copyright (c) 2017 Paul Kramme All Rights Reserved.")

	ConfigLocation := flag.String("c", "", "Specifies configfile location")
	flag.Parse()

	Config := new(Configuration)

	err := ini.MapTo(Config, *ConfigLocation)
	if err != nil {
		panic(err)
	}

	var f *os.File
	if Config.LogFileLocation != "" {
		f, err = os.OpenFile(Config.LogFileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	log.Println("BTSOOT started")

	Data := new(Block)
	err = Load(Config.DBFileLocation, Data)
	if err != nil {
		fmt.Println("Datafile not found. Please initialize the file")
	}

	Scan.Scans[time.Now().Format(time.RFC3339)] = ScanFiles(Config.Source, Config.MaxWorkerThreads)
	err = Save(Config.DBFileLocation, NewScan)
	if err != nil {
		panic(err)
	}
}
