package main

import "os"
import "time"

type Duration struct {
	time.Duration
}

type File struct {
	Block
	Path     string
	Finfo    os.FileInfo
	Checksum string
}

type Block struct {
	Blockname string
	Source      string
	Destination string
	Interval    Duration
}

type Configuration struct {
	LogFileLocation string
}
