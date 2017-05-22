package main

import "os"
import "time"

type Duration struct {
	time.Duration
}

type File struct {
	Path     string
	Finfo    os.FileInfo
	Checksum string
}

type Block struct {
	Source      string
	Destination string
	Interval    Duration
}

type Configuration struct {
	LogFileLocation string
	Blocks          map[string]Block
}
