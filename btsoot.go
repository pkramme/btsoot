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
	"flag"
	"fmt"
	"time"
	"github.com/paulkramme/toml"
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
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
