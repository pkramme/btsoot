package main

import (
	"os"
	"time"
)

type File struct {
	Path     string
	Finfo    os.FileInfo
	Checksum string
}

type Block struct {
	Blockname   string
	Source      string
	Destination string
	Interval    Duration
	Files       []File
}

type Duration struct {
	time.Duration
}

// This function is needed to convert intervals (4m3s) to understandable formats
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type WebResponse struct {
	Type    string
	Message string
}
