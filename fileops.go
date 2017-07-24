package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Save saves a datastructure to a specified path with json encoding.
func Save(path string, object interface{}) error {
	encodedobject, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	err = ioutil.WriteFile(path, encodedobject, 0644)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return err
}

// Load loads json data from a file into a datastructure.
func Load(path string, object interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	json.Unmarshal(file, object)
	return err
}
