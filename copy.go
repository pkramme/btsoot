package main

import (
	"io"
	"os"
)

// Copy copies files and returnes error if it fails.
func Copy(Source string, Destination string) (err error) {
	fdSource, err := os.Open(Source)
	if err != nil {
		return
	}
	defer fdSource.Close()
	fdDestination, err := os.Create(Destination)
	if err != nil {
		return
	}
	defer fdDestination.Close()
	buf := make([]byte, 512)
	_, err = io.CopyBuffer(fdDestination, fdSource, buf)
	if err != nil {
		return
	}
	return
}
