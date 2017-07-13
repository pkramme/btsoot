package main

type File struct {
	Path      string
	Name      string
	Checksum  string
	Directory bool
	Size      int64
}

type Block struct {
	Version string
	Scans   map[string][]File
}
