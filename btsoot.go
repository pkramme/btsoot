package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	cli "gopkg.in/urfave/cli.v1"
)

const (
	// Version string used for all version checking and CLI assignment
	Version = "0.7.0"

	// StopCode signals a thread to stop.
	StopCode = 1000
	// ConfirmCode confirms the execution of a signal.
	ConfirmCode = 1001
	// ErrorCode denies the save execution of a signal due to an error.
	ErrorCode = 1002
)

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Copyright = "Copyright (c) 2017 Paul Kramme All Rights Reserved. Distributed under BSD 3-Clause License."
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Paul Kramme",
			Email: "pjkramme@gmail.com",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Initialize a new block",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Specifies the location of the configfile",
				},
			},
			Action: func(c *cli.Context) error {
				Config, err := LoadConfig(c.String("config"))
				if err != nil {
					panic(err)
				}
				df, err := os.Create(Config.DBFileLocation)
				df.Close()

				var f *os.File
				if Config.LogFileLocation != "" {
					f, err = os.OpenFile(Config.LogFileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						log.Fatalln(err)
					}
					defer f.Close()
					log.SetOutput(f)
				}

				Data := new(Block)

				Data.Scans = make(map[time.Time][]File)
				Data.Scans[time.Now()] = ScanFiles(Config.Source, Config.MaxWorkerThreads)
				Data.Version = Version
				err = Save(Config.DBFileLocation, Data)
				if err != nil {
					log.Println(err)
					panic(err)
				}
				return err
			},
		},

		{
			Name:  "backup",
			Usage: "Backups the dataset",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Specifies the location of the configfile",
				},
			},
			Action: func(c *cli.Context) error {
				Config, err := LoadConfig(c.String("config"))
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

				Data := new(Block)
				err = Load(Config.DBFileLocation, Data)
				if err != nil {
					log.Println(err)
					fmt.Println("Datafile not found. Please initialize the file")
				}
				if Data.Version == "0.7.0" {
					fmt.Println("Block Version is 0.7.0")
				}
				Data.Scans[time.Now()] = ScanFiles(Config.Source, Config.MaxWorkerThreads)
				sortingslice := make(timeSlice, 0, len(Data.Scans))
				for k := range Data.Scans {
					sortingslice = append(sortingslice, k)
				}

				sort.Sort(sortingslice)
				fmt.Println(sortingslice[len(sortingslice)-1])
				fmt.Println(sortingslice[len(sortingslice)-2])

				newandchanged, deleted := Compare(Data.Scans[sortingslice[len(sortingslice)-1]], Data.Scans[sortingslice[len(sortingslice)-2]])

				for i, v := range newandchanged {
					fmt.Println("NEW:", i, v.Path, v.Checksum)
				}

				for i, v := range deleted {
					fmt.Println("DEL:", i, v.Path, v.Checksum)
				}
				err = Save(Config.DBFileLocation, Data)
				if err != nil {
					log.Println(err)
					panic(err)
				}
				return err
			},
		},
	}
	app.Run(os.Args)
}

type timeSlice []time.Time

func (ts timeSlice) Len() int {
	return len(ts)
}

func (ts timeSlice) Less(i, j int) bool {
	return ts[i].Before(ts[j])
}

func (ts timeSlice) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}
