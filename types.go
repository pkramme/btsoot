package main

import "time"

// File is the struct which represents a file.
type File struct {
	Path      string
	Name      string
	Checksum  string
	Directory bool
	Size      int64
}

// Block represents a file superstructure.
type Block struct {
	Version string
	Scans   map[time.Time][]File
}
