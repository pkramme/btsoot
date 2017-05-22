package main

import (
	"flag"
	"fmt"
	"github.com/paulkramme/toml"
)

func main() {
	fmt.Println("BTSOOT - Copyright (C) 2016-2017 Paul Kramme")

	//Verbose := flag.Bool("verbose", false, "Verbose output for better debugging or just to see whats going on. This can slow BTSOOT down.")
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
