package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paulkramme/toml"
	"log"
	"os"
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

	var Config Configuration

	_, err := toml.DecodeFile(*ConfigLocation, &Config)
	if err != nil {
		fmt.Println("Couldn't find or open config file.")
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

	// NOTE: Create a spacer in the log
	log.Println("\n\n")

	db, err := sql.Open("sqlite3", Config.DBFileLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	err = DatabaseSetup(db)
	if err != nil {
		panic(err)
	}

	_ = ScanFiles(Config.Source, Config.MaxWorkerThreads)
}

func DatabaseSetup(db *sql.DB) (err error) {
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS blocks(
			name text PRIMARY KEY NOT NULL,
			interval VARCHAR,
			path text NOT NULL,
			destination text NOT NULL,
			absolutepath BOOL)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files(
			block text NOT NULL,
			checksum TEXT NOT NULL,
			filename text,
			absolutepath text,
			FOREIGN KEY (block) REFERENCES blocks(name))`)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS users(
			username text,
			password hash)`)
	if err != nil {
		return err
	}
	// _, err = db.Exec("CREATE TABLE IF NOT EXISTS version")
	// if err != nil {
	// 	return err
	// }
	return
}
