package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func sha256sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}
	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func WorkFiller(in chan File, out chan File) {
	
}

func Worker(in chan File, out chan File) {

}

func scanfiles(location string, WorkFiller chan File) (err error) {
	var walkcallback = func(path string, fileinfo os.FileInfo, inputerror error) (err error) {
		var f File
		f.Path = path
		f.Finfo = fileinfo
		WorkFiller <- f
		return
	}
	err = filepath.Walk(location, walkcallback)
	return
}
