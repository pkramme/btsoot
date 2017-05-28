/*
TODO: create startup
- Update check async
- Create and load database
- load jobs from database
- check if blocks exist
- start webserver async
- listen for SIGINT
- start all jobs in their own timer thread -> after x amount of time, start scan,
*/
package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paulkramme/toml"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("BTSOOT - Copyright (C) 2016-2017 Paul Kramme")

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
	NumberOfJobs := 0 // This should hold the number of all jobs loaded from the database
	signals := make(chan os.Signal, 1)
	killall := make(chan bool, 1)
	done := make(chan bool, 1)

	//MasterWorkerGroup.Add(1)
	for j := 0; j < NumberOfJobs; j++ {
		go Job(killall, done)
	}

	signal.Notify(signals, syscall.SIGINT)
	<-signals
	fmt.Println("Exiting. This may take a while.")
	killall <- true
	for j := 0; j < NumberOfJobs; j++ {
		<-done
	}
}


func Job(killall chan bool, done chan bool) {
	killtiny := make(chan bool, 1)
	go Something(killtiny)
	<-killall
	fmt.Println("Setting killtiny to true")
	killtiny <- true
	done <- true
	return
}

func Something(killall chan bool) {
	var cnt int
	for {
		cnt++
		fmt.Println("Hello World", cnt)
		select {
		case <-killall:
			fmt.Println("Killing something")
			return
		default:
		}
	}
}


func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
