package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func sha256sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	buf := make([]byte, 64)

	hash := sha256.New()
	_, err = io.CopyBuffer(hash, file, buf)
	if err != nil {
		return
	}
	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func sha512sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	buf := make([]byte, 128)

	hash := sha512.New()
	_, err = io.CopyBuffer(hash, file, buf)
	if err != nil {
		return
	}
	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func worker(in chan File, out chan File, comm chan bool) {
	for {
		FileToProcess, ok := <-in
		if ok != true {
			comm <- true
			return
		}
		if FileToProcess.Directory == true {
			out <- FileToProcess
			continue
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

// ScanFiles takes a folder and a maximum thread number, scans the directory, and returnes File type with SHA512 checksums.
func ScanFiles(location string, MaxWorkerThreads int) (files []File) {
	WFin := make(chan File)
	WFout := make(chan File)

	CheckIfDone := make(chan bool)
	WorkerMap := make(map[int]chan bool)
	for i := MaxWorkerThreads; i > 0; i-- {
		WorkerMap[i] = make(chan bool)
		go worker(WFin, WFout, WorkerMap[i])
	}
	var walkcallback = func(path string, fileinfo os.FileInfo, inputerror error) (err error) {
		if inputerror != nil {
			fmt.Println(inputerror)
			return
		}
		var f File
		f.Path = path
		f.Name = fileinfo.Name()
		f.Size = fileinfo.Size()
		f.Directory = fileinfo.IsDir()
		WFin <- f
		return
	}
	go func() {
		err := filepath.Walk(location, walkcallback)
		close(WFin)
		if err != nil {
			panic(err)
		}
		CheckIfDone <- true
	}()

Resultloop:
	for {
		file := <-WFout
		files = append(files, file)
		select {
		case <-CheckIfDone:
			for len(WorkerMap) != 0 {
				for key, value := range WorkerMap {
					select {
					case <-value:
						delete(WorkerMap, key)
					default:
						file := <-WFout
						files = append(files, file)
					}
				}
			}
			break Resultloop
		default:
		}
	}
	for i, v := range files {
		fmt.Println(i, v.Path, v.Checksum)
	}
	return
}
