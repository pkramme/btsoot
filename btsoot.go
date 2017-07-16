package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
				df, err := os.Create(Config.DataFileLocation)
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
				newfiles := ScanFiles(Config.Source, Config.MaxWorkerThreads)

				for _, v := range newfiles {
					// Create dirs first
					if v.Directory == true {
						os.MkdirAll(filepath.Join(Config.Destination, v.Path), 0777)
					}
				}

				for _, v := range newfiles {
					// fmt.Println("NEW:", i, v.Path, v.Checksum)
					if v.Directory == false {
						err := Copy(filepath.Join(Config.Source, v.Path), filepath.Join(Config.Destination, v.Path))
						if err != nil {
							log.Println(err)
						}
					}
				}

				Data.Scans[time.Now()] = newfiles
				Data.Version = Version
				err = Save(Config.DataFileLocation, Data)
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
				err = Load(Config.DataFileLocation, Data)
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

				// for _, v := range Data.Scans[sortingslice[len(sortingslice)-1]] {
				// 	fmt.Println("NEWSCN:", v)
				// }

				newandchanged, deleted := Compare(Data.Scans[sortingslice[len(sortingslice)-1]], Data.Scans[sortingslice[len(sortingslice)-2]])

				// for _, v := range deleted {
				// 	fmt.Println("DEL:", v)
				// }
				// for _, v := range newandchanged {
				// 	fmt.Println("NEW:", v)
				// }

				for _, v := range deleted {
					err := os.RemoveAll(filepath.Join(Config.Destination, v.Path))
					if err != nil {
						log.Println(err)
						fmt.Println(err)
					}
				}

				// Create dirs first
				for _, v := range newandchanged {
					if v.Directory == true {
						err := os.MkdirAll(filepath.Join(Config.Destination, v.Path), 0777)
						if err != nil {
							log.Println(err)
							panic(err)
						}
					}
				}
				// Now copy all files
				for _, v := range newandchanged {
					if v.Directory == false {
						err := Copy(filepath.Join(Config.Source, v.Path), filepath.Join(Config.Destination, v.Path))
						if err != nil {
							log.Println(err)
							panic(err)
						}
					}
				}
				err = Save(Config.DataFileLocation, Data)
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
