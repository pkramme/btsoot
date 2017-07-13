package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
	"gopkg.in/urfave/cli.v1"
)

const (
	StopCode    = 1000
	ConfirmCode = 1001
	ErrorCode   = 1002
)

func main() {
	app := cli.NewApp()
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
				Config := new(Configuration)
				err := ini.MapTo(Config, c.String("config"))
				if err != nil {
					panic(err)
				}
				f, err := os.Create(Config.DBFileLocation)
				f.Close()
				if err == nil {
					Data := new(Block)
					Data.Version = "0.7.0"
					Save(Config.DBFileLocation, Data)
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
				Config := new(Configuration)

				err := ini.MapTo(Config, c.String("config"))
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
				Data.Scans[time.Now().Format(time.RFC3339)] = ScanFiles(Config.Source, Config.MaxWorkerThreads)
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
