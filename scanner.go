package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

func Worker(in chan File, out chan File) {
	for {
		FileToProcess, ok := <-in
		if ok == false {
							fmt.Println("Shutting down a thread :)")
			return
		}
		hash, err := sha512sum(FileToProcess.Path)
		if err != nil {
			log.Println(err)
			continue
		}
		FileToProcess.Checksum = hash
		out <- FileToProcess
	}
}

func scanfiles(location string, MaxWorkerThreads int, comm chan int) (files []File) {
	WFin := make(chan File)
	WFout := make(chan File)
	for i := MaxWorkerThreads; i > 0; i-- {
		go Worker(WFin, WFout)
	}
	var walkcallback = func(path string, fileinfo os.FileInfo, inputerror error) (err error) {
		var f File
		f.Path = path
		f.Finfo = fileinfo
		WFin <- f
		return
	}
	go func() {
		err := filepath.Walk(location, walkcallback)
		close(WFin)
		if err != nil {
			panic(err)
		}
	}()
	for {
		select {
		case file := <-WFout:
			fmt.Println(file)
			files = append(files, file)
		default:
			break
		}
		select {
		case msg := <-comm:
			if msg == StopCode {
				close(WFin)
				files = nil
				break
			}
		default:
		}
	}
	fmt.Println(files)
	close(WFin)
	return
}
