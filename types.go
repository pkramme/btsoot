package main

type File struct {
	Path     string
	Name string
	Checksum string
	Directory bool
	Size int64
}

type Configuration struct {
	LogFileLocation  string
	DBFileLocation   string
	MaxWorkerThreads int
	Source           string
}

type Block struct {
	Scans map[string][]File
}
