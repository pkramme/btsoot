package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

// Save saves a datastructure to a specified path with gob encoding.
func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		err := encoder.Encode(object)
		if err != nil {
			fmt.Println(err)
		}
	}
	file.Close()
	return err
}

// Load loads gob data from a file into a datastructure.
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
