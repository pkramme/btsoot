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

func scanfiles(location string) (m map[string]string, err error) {
	m = make(map[string]string)
	var walkcallback = func(path string, fileinfo os.FileInfo, inputerror error) (err error) {
		checksum, _ := sha256sum(path)
		m[path] = checksum
		return
	}
	err = filepath.Walk(location, walkcallback)
	return
}
