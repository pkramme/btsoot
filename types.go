package main

import (
	"os"
)

type File struct {
	Path     string
	Finfo    os.FileInfo
	Checksum string
}

type Configuration struct {
	LogFileLocation  string
	DBFileLocation   string
	MaxWorkerThreads int
	Source           string
}
