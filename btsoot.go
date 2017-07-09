package main

import (
	"fmt"
	"github.com/paulkramme/cli"
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
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Specifies the location of the configfile",
		},
		cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Initialize a new block",
			Action: func(c *cli.Context) {
				// Init(c.String("config"))
				fmt.Println("init")
				fmt.Println(c.String("lang"))
			},
		},
		{
			Name:  "license",
			Usage: "Show all licenses associated with this project",
			Action: func(c *cli.Context) {
				License()
			},
		},
		{
			Name:  "backup",
			Usage: "Backups the dataset",
			Action: func(c *cli.Context) {
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

				log.Println("BTSOOT started")

				Data := new(Block)
				err = Load(Config.DBFileLocation, Data)
				if err != nil {
					fmt.Println("Datafile not found. Please initialize the file")
				}

				Data.Scans[time.Now().Format(time.RFC3339)] = ScanFiles(Config.Source, Config.MaxWorkerThreads)
				err = Save(Config.DBFileLocation, Data)
				if err != nil {
					panic(err)
				}
			},
		},
	}
	app.Run(os.Args)
}
