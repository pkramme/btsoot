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

// Block is a file superstructure which also contains a version string,
// which can be used for backwards compatibility.
type Block struct {
	Version string
	Scans   map[time.Time][]File
}
