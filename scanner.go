package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

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

func Worker(in chan File, out chan File) {
	for {
		FileToProcess := <-in
		if FileToProcess.Finfo.IsDir() {
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

func ScanningProcess(procconfig Process, config Configuration) {
	log.Println("SCANNERPROC: Startup complete")
	procconfig.Subprocesses = make(map[int]Process)
	//go scanfiles(".", 4, scanfilescomm)
	for {
		select {
		case comm := <-procconfig.Channel:
			if comm == StopCode {
				log.Println("SCANNERPROC: Shutdown")
				procconfig.Channel <- ConfirmCode
				return
			}
		default:
			time.Sleep(100)
		}
	}
}
