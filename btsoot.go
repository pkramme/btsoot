package main

import "fmt"
import "flag"
import "encoding/json"
import "io/ioutil"

type location struct {
	Name                  string
	Source                string
	Dest                  string
	More_special_settings bool
}

type config struct {
	Some_generell_settings bool
	Locations              []location
}

func main() {
	fmt.Println("BTSOOT - Copyright (C) 2016-2017 Paul Kramme")

	verbose := flag.Bool("verbose", false, "Verbose output for better debugging or just to see whats going on. This can slow BTSOOT down.")
	add_new := flag.String("add", "", "Add new block")
	add_new_src := flag.String("src", "", "Add new source location, can only be used with -add")
	add_new_dest := flag.String("dest", "", "Add new destination, can only be used with -add")
	rm := flag.String("rm", "", "Remove a block from config")
	flag.Parse()

	var conf config
	configfile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	err = fromjson(string(configfile), &conf)
	if err != nil {
		panic(err)
	}

	if *verbose == true {
		fmt.Println("Verbose printing activated.")
	}

	if *add_new != "" {
		var loc location
		loc.Name = *add_new
		if *add_new_src != "" {
			loc.Source = *add_new_src
		}
		if *add_new_dest != "" {
			loc.Dest = *add_new_dest
		}
		conf.Locations = append(conf.Locations, loc)
	}

	if *rm != "" {
		for n, location_iterr := range conf.Locations {
			if location_iterr.Name == *rm {
				// Removing an slice element without preserving order
				conf.Locations[n] = conf.Locations[len(conf.Locations)-1]
				conf.Locations = conf.Locations[:len(conf.Locations)-1]
			}
		}
	}
	resulting_config, err := json.MarshalIndent(conf, "", "    ")
	err = ioutil.WriteFile("./config.json", resulting_config, 0664)
	if err != nil {
		panic(err)
	}
}
