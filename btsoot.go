/*BTSOOT
BTSOOT is crossplatform incremental backup tool written in golang.
It is able to identify files based on lastmod date or by a blake2b
checksum. It then copies changed files to a remote destination.

BTSOOT is licensed under BSD 3 Clause and is created by Paul Kramme.
*/

package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	cli "gopkg.in/urfave/cli.v1"
)

const (
	// Version string used for all version checking and CLI assignment
	Version = "0.7.0"
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
				var newfiles []File
				if Config.Scantype.Blake2bBased {
					newfiles = ScanFilesBlake2b(Config.Source, Config.MaxWorkerThreads)
				} else if Config.Scantype.TimestampBased {
					newfiles = ScanFilesTimestamp(Config.Source)
				} else {
					panic("Unsupported scantype")
				}

				for _, v := range newfiles {
					// Create dirs first
					if v.Directory == true {
						os.MkdirAll(filepath.Join(Config.Destination, v.Path), 0777)
					}
				}

				for _, v := range newfiles {
					if v.Directory == false {
						if Config.Copy.UseExternalCopy {
							cmd := exec.Command(Config.Copy.ExternalCopyPath, filepath.Join(Config.Source, v.Path), filepath.Join(Config.Destination, v.Path))
							var out bytes.Buffer
							cmd.Stdout = &out
							err := cmd.Run()
							if err != nil {
								log.Println(err)
							}
							log.Print(out.String())
						} else {
							err := CopyFile(filepath.Join(Config.Source, v.Path), filepath.Join(Config.Destination, v.Path))
							if err != nil {
								log.Println(err)
								panic(err)
							}
						}
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
				cli.BoolFlag{
					Name:  "override, o",
					Usage: "Overrides the saveguard. It has to be enabled. USE WITH CAUTION!",
				},
				cli.BoolFlag{
					Name:  "dry-run, d",
					Usage: "Scans the block, but doesn't change or copies anything",
				},
			},
			Action: func(c *cli.Context) error {
				Config, err := LoadConfig(c.String("config"))
				if err != nil {
					panic(err)
				}

				if Config.LogFileLocation != "" {
					var f *os.File
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
					panic("Datafile not found. Please initialize the file.")

				}
				if Data.Version == "0.7.0" {
					log.Println("Block Version is 0.7.0")
				}
				log.Println("Scan initialized")

				if Config.Scantype.Blake2bBased {
					Data.Scans[time.Now()] = ScanFilesBlake2b(Config.Source, Config.MaxWorkerThreads)
				} else if Config.Scantype.TimestampBased {
					Data.Scans[time.Now()] = ScanFilesTimestamp(Config.Source)
				} else {
					panic("Unsupported scantype")
				}

				log.Println("Scan finished.")
				sortingslice := make(timeSlice, 0, len(Data.Scans))
				for k := range Data.Scans {
					sortingslice = append(sortingslice, k)
				}

				sort.Sort(sortingslice)
				log.Println(sortingslice[len(sortingslice)-1])
				log.Println(sortingslice[len(sortingslice)-2])

				newandchanged, deleted := Compare(Data.Scans[sortingslice[len(sortingslice)-1]], Data.Scans[sortingslice[len(sortingslice)-2]])

				if Config.SaveguardEnable {
					scanlen := len(Data.Scans[sortingslice[len(sortingslice)-1]])
					deletedlen := len(deleted)
					percentage := (deletedlen / scanlen) * 100
					if percentage >= Config.SaveguardMaxPercentage {
						if c.Bool("override") != true {
							log.Println("The change percentage exceeds the maximum saveguard percentage. Aborting.")
							os.Exit(1)
						}
						log.Println("The change percentage exceeds the maximum saveguard percentage, but the override flag is set.")
					}

				}

				log.Println("New or changed:", len(newandchanged))
				log.Println("Deleted or changed:", len(deleted))
				if c.Bool("dry-run") {
					log.Println("dry-run flag is set, quitting.")
					return nil
				}

				for _, v := range deleted {
					err := os.RemoveAll(filepath.Join(Config.Destination, v.Path))
					if err != nil {
						log.Println(err)
						panic(err)
					}
				}

				// Create dirs first
				for _, v := range newandchanged {
					if v.Directory {
						err := os.MkdirAll(filepath.Join(Config.Destination, v.Path), 0777)
						if err != nil {
							log.Println(err)
							panic(err)
						}
					}
				}
				// Now copy all files
				for _, v := range newandchanged {
					if !v.Directory {
						if Config.Copy.UseExternalCopy {
							exec.Command(Config.Copy.ExternalCopyPath, filepath.Join(Config.Source, v.Path), filepath.Join(Config.Destination, v.Path))
						} else {
							err := CopyFile(filepath.Join(Config.Source, v.Path), filepath.Join(Config.Destination, v.Path))
							if err != nil {
								log.Println(err)
								panic(err)
							}
						}
					}
				}

				if len(Data.Scans) > 3 {
					delete(Data.Scans, sortingslice[0])
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
