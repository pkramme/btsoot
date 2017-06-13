package main

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"fmt"
	"log"
	"errors"
)

func sha512sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := sha512.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}
	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func WorkFiller(in chan File, out chan File) {
	Win := make(chan File)
	Wout := make(chan File)
	for i := 4; i != 0; i-- {
		fmt.Println(i)
		go Worker(Win, Wout, i)
	}
	for {
		FileToProcess, ok := <-in
		if ok == false {
			close(Win)
			return
		}
		Win <- FileToProcess
	}
}

func Worker(in chan File, out chan File, i int) {
	for {
		FileToProcess, ok := <-in
		if ok == false {
			return
		}
		hash, err := sha512sum(FileToProcess.Path)
		if err != nil {
			log.Println(err)
			continue
		}
		FileToProcess.Checksum = hash
		//out <- FileToProcess
		fmt.Println(i, FileToProcess.Path)
	}
}

func scanfiles(location string, comm chan int) (err error) {
	WFin := make(chan File)
	WFout := make(chan File)
	go WorkFiller(WFin, WFout)
	var walkcallback = func(path string, fileinfo os.FileInfo, inputerror error) (err error) {
		select {
		case _, ok := <-comm:
			if ok == false {
				fmt.Println("Shutting down scanner")
				return errors.New("Shutting down filewalker. THIS IS NOT AN ERROR!")
			}
		default:
		}
		var f File
		f.Path = path
		f.Finfo = fileinfo
		WFin <- f
		return
	}
	err = filepath.Walk(location, walkcallback)
	close(WFin)
	return
}
