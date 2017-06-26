package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Many information about the blocks, such as possible blocked blocks, may also be used as a healthcheck")
}

func NewBlockHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("Creating a new block:", ps.ByName("name"))
	// if block exists, abort
	// if not, create
}
